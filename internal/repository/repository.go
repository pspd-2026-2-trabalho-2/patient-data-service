// Package repository executa as consultas SQL sobre o pseudo-prontuário eletrônico.
package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/domain"
	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/observability"
)

var ErrNotFound = errors.New("registro não encontrado")

type Repository struct {
	pool    *pgxpool.Pool
	metrics *observability.Metrics
}

func New(pool *pgxpool.Pool, m *observability.Metrics) *Repository {
	return &Repository{pool: pool, metrics: m}
}

// COALESCE evita erro ao ler colunas TEXT nulas.
const (
	patientCols = `p.patient_id, p.full_name, p.birth_date, p.gender,
		COALESCE(p.city,''), COALESCE(p.state,''), COALESCE(p.cpf,''), COALESCE(p.cns,'')`

	encounterCols = `encounter_id, patient_id, start_date, end_date,
		COALESCE(encounter_type,''), COALESCE(department,'')`

	// value é varchar no banco: converte para número só quando for numérico (senão NULL).
	eventCols = `event_id, patient_id, COALESCE(encounter_id,''), event_type, code,
		COALESCE(description,''), event_date,
		CASE WHEN value ~ '^-?[0-9]+(\.[0-9]+)?$' THEN value::double precision END,
		COALESCE(unit,'')`

	projectCols = `project_id, title, researcher_username, target_condition_code, status, valid_until`
)

func (r *Repository) PatientsByDoctor(ctx context.Context, doctor string, yield func(domain.Patient) error) error {
	return r.streamPatients(ctx, "PatientsByDoctor",
		`SELECT `+patientCols+` FROM patients p
		 JOIN user_patient_assignments a ON a.patient_id = p.patient_id
		 WHERE a.username = $1 AND UPPER(a.assignment_type) = 'ATTENDING' AND a.active
		 ORDER BY p.patient_id`, yield, doctor)
}

func (r *Repository) SupervisedPatients(ctx context.Context, intern string, yield func(domain.Patient) error) error {
	return r.streamPatients(ctx, "SupervisedPatients",
		`SELECT `+patientCols+` FROM patients p
		 JOIN user_patient_assignments a ON a.patient_id = p.patient_id
		 WHERE a.username = $1 AND UPPER(a.assignment_type) = 'TRAINEE' AND a.active
		 ORDER BY p.patient_id`, yield, intern)
}

func (r *Repository) CohortPatients(ctx context.Context, conditionCode string, yield func(domain.Patient) error) error {
	return r.streamPatients(ctx, "CohortPatients",
		`SELECT DISTINCT `+patientCols+` FROM patients p
		 JOIN clinical_events e ON e.patient_id = p.patient_id
		 WHERE UPPER(e.event_type) = 'CONDITION' AND UPPER(e.code) = UPPER($1)
		 ORDER BY p.patient_id`, yield, conditionCode)
}

func (r *Repository) GetPatient(ctx context.Context, patientID string) (p *domain.Patient, err error) {
	start := time.Now()
	found := 0
	defer func() { r.metrics.RecordQuery("GetPatient", start, found, err) }()

	row := r.pool.QueryRow(ctx, `SELECT `+patientCols+` FROM patients p WHERE p.patient_id = $1`, patientID)
	var out domain.Patient
	err = row.Scan(&out.PatientID, &out.FullName, &out.BirthDate, &out.Gender,
		&out.City, &out.State, &out.CPF, &out.CNS)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = ErrNotFound
		}
		return nil, err
	}
	found = 1
	return &out, nil
}

// streamPatients executa a consulta e chama yield para cada paciente, conforme os
// registros são lidos do cursor do banco — sem bufferizar o resultado inteiro. É o
// que permite atender coortes de dezenas de milhares de pacientes por server streaming.
func (r *Repository) streamPatients(ctx context.Context, name, sql string, yield func(domain.Patient) error, args ...any) (err error) {
	start := time.Now()
	sent := 0
	defer func() { r.metrics.RecordQuery(name, start, sent, err) }()

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var p domain.Patient
		if err = rows.Scan(&p.PatientID, &p.FullName, &p.BirthDate, &p.Gender,
			&p.City, &p.State, &p.CPF, &p.CNS); err != nil {
			return err
		}
		if err = yield(p); err != nil {
			return err
		}
		sent++
	}
	err = rows.Err()
	return err
}

func (r *Repository) Encounters(ctx context.Context, patientID string) (out []domain.Encounter, err error) {
	start := time.Now()
	defer func() { r.metrics.RecordQuery("Encounters", start, len(out), err) }()

	rows, err := r.pool.Query(ctx,
		`SELECT `+encounterCols+` FROM encounters WHERE patient_id = $1 ORDER BY start_date DESC`, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var e domain.Encounter
		if err = rows.Scan(&e.EncounterID, &e.PatientID, &e.StartDate, &e.EndDate,
			&e.EncounterType, &e.Department); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	err = rows.Err()
	return out, err
}

// ClinicalEvents ordena do mais recente para o mais antigo; eventType vazio traz todos.
func (r *Repository) ClinicalEvents(ctx context.Context, patientID, eventType string) ([]domain.ClinicalEvent, error) {
	if eventType == "" {
		return r.queryEvents(ctx, "ClinicalEvents",
			`SELECT `+eventCols+` FROM clinical_events WHERE patient_id = $1 ORDER BY event_date DESC`,
			patientID)
	}
	return r.queryEvents(ctx, "ClinicalEvents",
		`SELECT `+eventCols+` FROM clinical_events WHERE patient_id = $1 AND UPPER(event_type) = UPPER($2) ORDER BY event_date DESC`,
		patientID, eventType)
}

// ClinicalHistory ordena em ordem temporal crescente.
func (r *Repository) ClinicalHistory(ctx context.Context, patientID string) ([]domain.ClinicalEvent, error) {
	return r.queryEvents(ctx, "ClinicalHistory",
		`SELECT `+eventCols+` FROM clinical_events WHERE patient_id = $1 ORDER BY event_date ASC`, patientID)
}

func (r *Repository) queryEvents(ctx context.Context, name, sql string, args ...any) (out []domain.ClinicalEvent, err error) {
	start := time.Now()
	defer func() { r.metrics.RecordQuery(name, start, len(out), err) }()

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var e domain.ClinicalEvent
		if err = rows.Scan(&e.EventID, &e.PatientID, &e.EncounterID, &e.EventType, &e.Code,
			&e.Description, &e.EventDate, &e.Value, &e.Unit); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	err = rows.Err()
	return out, err
}

func (r *Repository) ProjectsByResearcher(ctx context.Context, researcher string) (out []domain.Project, err error) {
	start := time.Now()
	defer func() { r.metrics.RecordQuery("ProjectsByResearcher", start, len(out), err) }()

	rows, err := r.pool.Query(ctx,
		`SELECT `+projectCols+` FROM projects WHERE researcher_username = $1 ORDER BY project_id`, researcher)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p domain.Project
		if err = rows.Scan(&p.ProjectID, &p.Title, &p.ResearcherUsername, &p.ConditionCode,
			&p.Status, &p.ValidUntil); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	err = rows.Err()
	return out, err
}

// CheckAssignment verifica vínculo ativo; role vazio ignora o tipo. Retorna (permitido, tipo).
func (r *Repository) CheckAssignment(ctx context.Context, username, patientID, role string) (allowed bool, assignmentType string, err error) {
	start := time.Now()
	found := 0
	defer func() { r.metrics.RecordQuery("CheckAssignment", start, found, err) }()

	sql := `SELECT assignment_type FROM user_patient_assignments
			WHERE username = $1 AND patient_id = $2 AND active`
	args := []any{username, patientID}
	if at := mapAssignmentType(role); at != "" {
		sql += ` AND UPPER(assignment_type) = $3`
		args = append(args, at)
	}
	sql += ` LIMIT 1`

	err = r.pool.QueryRow(ctx, sql, args...).Scan(&assignmentType)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}
	found = 1
	return true, assignmentType, nil
}

// mapAssignmentType traduz o papel (role) para o valor da coluna assignment_type do banco.
func mapAssignmentType(role string) string {
	switch strings.ToUpper(role) {
	case "MEDICO", "ATTENDING":
		return "ATTENDING"
	case "ESTAGIARIO", "TRAINEE":
		return "TRAINEE"
	default:
		return "" // sem filtro por tipo
	}
}

// cohortCTE seleciona os pacientes da coorte (com a condição informada em $1).
const cohortCTE = `WITH cohort AS (
	SELECT DISTINCT patient_id FROM clinical_events WHERE UPPER(event_type) = 'CONDITION' AND UPPER(code) = UPPER($1)
)`

func (r *Repository) CohortTotal(ctx context.Context, conditionCode string) (total int64, err error) {
	start := time.Now()
	defer func() { r.metrics.RecordQuery("CohortTotal", start, 1, err) }()

	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(DISTINCT patient_id) FROM clinical_events
		 WHERE UPPER(event_type) = 'CONDITION' AND UPPER(code) = UPPER($1)`, conditionCode).Scan(&total)
	return total, err
}

func (r *Repository) CohortBySex(ctx context.Context, conditionCode string) ([]domain.Count, error) {
	return r.queryCounts(ctx, "CohortBySex",
		cohortCTE+`
		SELECT p.gender, COUNT(*) FROM patients p
		JOIN cohort c ON c.patient_id = p.patient_id
		GROUP BY p.gender ORDER BY p.gender`, conditionCode)
}

func (r *Repository) CohortByAgeRange(ctx context.Context, conditionCode string) ([]domain.Count, error) {
	return r.queryCounts(ctx, "CohortByAgeRange",
		cohortCTE+`,
		ages AS (
			SELECT date_part('year', age(p.birth_date))::int AS age
			FROM patients p JOIN cohort c ON c.patient_id = p.patient_id
		)
		SELECT CASE
			WHEN age BETWEEN 0 AND 17  THEN '0-17'
			WHEN age BETWEEN 18 AND 39 THEN '18-39'
			WHEN age BETWEEN 40 AND 59 THEN '40-59'
			ELSE '60+'
		END AS age_range, COUNT(*)
		FROM ages GROUP BY age_range ORDER BY age_range`, conditionCode)
}

func (r *Repository) CohortMedicationFrequency(ctx context.Context, conditionCode string) ([]domain.Count, error) {
	return r.queryCounts(ctx, "CohortMedicationFrequency",
		cohortCTE+`
		SELECT e.code, COUNT(*) FROM clinical_events e
		JOIN cohort c ON c.patient_id = e.patient_id
		WHERE UPPER(e.event_type) = 'MEDICATION'
		GROUP BY e.code ORDER BY COUNT(*) DESC, e.code`, conditionCode)
}

// CohortByDepartment conta pacientes distintos por departamento de atendimento.
func (r *Repository) CohortByDepartment(ctx context.Context, conditionCode string) ([]domain.Count, error) {
	return r.queryCounts(ctx, "CohortByDepartment",
		cohortCTE+`
		SELECT e.department, COUNT(DISTINCT e.patient_id) FROM encounters e
		JOIN cohort c ON c.patient_id = e.patient_id
		WHERE e.department IS NOT NULL AND e.department <> ''
		GROUP BY e.department ORDER BY COUNT(DISTINCT e.patient_id) DESC, e.department`, conditionCode)
}

// CohortHbA1c retorna média e mediana (0 se não houver exames).
func (r *Repository) CohortHbA1c(ctx context.Context, conditionCode string) (mean, median float64, err error) {
	start := time.Now()
	defer func() { r.metrics.RecordQuery("CohortHbA1c", start, 1, err) }()

	err = r.pool.QueryRow(ctx,
		cohortCTE+`
		SELECT COALESCE(AVG(e.value::double precision), 0),
		       COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY e.value::double precision), 0)
		FROM clinical_events e JOIN cohort c ON c.patient_id = e.patient_id
		WHERE UPPER(e.event_type) = 'OBSERVATION' AND UPPER(e.code) = 'HBA1C'
		  AND e.value ~ '^-?[0-9]+(\.[0-9]+)?$'`, conditionCode).Scan(&mean, &median)
	return mean, median, err
}

func (r *Repository) queryCounts(ctx context.Context, name, sql string, args ...any) (out []domain.Count, err error) {
	start := time.Now()
	defer func() { r.metrics.RecordQuery(name, start, len(out), err) }()

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c domain.Count
		if err = rows.Scan(&c.Key, &c.Count); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	err = rows.Err()
	return out, err
}
