//go:build integration

// Testes de integração: exigem o Postgres com schema+seed (docker compose up -d db).
package repository_test

import (
	"context"
	"math"
	"os"
	"testing"

	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/db"
	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/domain"
	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/observability"
	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/repository"
)

// NewMetrics registra no registrador default; criar uma vez só (senão pânico de duplicidade).
var testMetrics = observability.NewMetrics()

func newRepo(t *testing.T) *repository.Repository {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://pspd:pspd@localhost:5433/hospital?sslmode=disable"
	}
	pool, err := db.NewPool(context.Background(), dsn)
	if err != nil {
		t.Fatalf("conexão com o banco: %v", err)
	}
	t.Cleanup(pool.Close)
	return repository.New(pool, testMetrics)
}

func TestPatientsByDoctor(t *testing.T) {
	r := newRepo(t)
	var ps []domain.Patient
	err := r.PatientsByDoctor(context.Background(), "med.cardoso", func(p domain.Patient) error {
		ps = append(ps, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 8 {
		t.Errorf("med.cardoso deveria ter 8 pacientes, veio %d", len(ps))
	}
}

func TestSupervisedPatients(t *testing.T) {
	r := newRepo(t)
	var ps []domain.Patient
	err := r.SupervisedPatients(context.Background(), "est.souza", func(p domain.Patient) error {
		ps = append(ps, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 3 {
		t.Errorf("est.souza deveria supervisionar 3 pacientes, veio %d", len(ps))
	}
}

func TestCohortPatientsDiabetes(t *testing.T) {
	r := newRepo(t)
	var ps []domain.Patient
	err := r.CohortPatients(context.Background(), "DIABETES", func(p domain.Patient) error {
		ps = append(ps, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 13 {
		t.Errorf("coorte Diabetes deveria ter 13 pacientes, veio %d", len(ps))
	}
}

func TestCohortTotalAndDistributions(t *testing.T) {
	r := newRepo(t)
	ctx := context.Background()

	total, err := r.CohortTotal(ctx, "DIABETES")
	if err != nil {
		t.Fatal(err)
	}
	if total != 13 {
		t.Errorf("total da coorte = %d, quer 13", total)
	}

	bySexCounts, err := r.CohortBySex(ctx, "DIABETES")
	if bySex := sumCounts(t, bySexCounts, err); bySex != 13 {
		t.Errorf("soma por sexo = %d, quer 13", bySex)
	}

	byAgeCounts, err := r.CohortByAgeRange(ctx, "DIABETES")
	if byAge := sumCounts(t, byAgeCounts, err); byAge != 13 {
		t.Errorf("soma por faixa etária = %d, quer 13", byAge)
	}
}

func TestCohortHbA1c(t *testing.T) {
	r := newRepo(t)
	mean, median, err := r.CohortHbA1c(context.Background(), "DIABETES")
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(median-7.6) > 0.001 {
		t.Errorf("mediana HbA1c = %v, quer 7.6", median)
	}
	if math.Abs(mean-7.7615) > 0.01 {
		t.Errorf("média HbA1c = %v, quer ~7.76", mean)
	}
}

func TestMedicationFrequency(t *testing.T) {
	r := newRepo(t)
	freq, err := r.CohortMedicationFrequency(context.Background(), "DIABETES")
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]int64{}
	for _, c := range freq {
		got[c.Key] = c.Count
	}
	if got["METFORMIN"] != 10 {
		t.Errorf("Metformina = %d, quer 10", got["METFORMIN"])
	}
	if got["INSULIN"] != 4 {
		t.Errorf("Insulina = %d, quer 4", got["INSULIN"])
	}
	if got["LOSARTAN"] != 2 {
		t.Errorf("Losartana = %d, quer 2 (P000013/14 não são da coorte)", got["LOSARTAN"])
	}
}

func TestCohortByDepartment(t *testing.T) {
	r := newRepo(t)
	counts, err := r.CohortByDepartment(context.Background(), "DIABETES")
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]int64{}
	for _, c := range counts {
		got[c.Key] = c.Count
	}
	if got["ENDOCRINOLOGY"] != 13 {
		t.Errorf("Endocrinologia = %d, quer 13", got["ENDOCRINOLOGY"])
	}
	if got["CARDIOLOGY"] != 2 {
		t.Errorf("Cardiologia = %d, quer 2 (P000001 e P000007)", got["CARDIOLOGY"])
	}
}

func TestCheckAssignment(t *testing.T) {
	r := newRepo(t)
	ctx := context.Background()

	allowed, atype, err := r.CheckAssignment(ctx, "med.cardoso", "P000001", "medico")
	if err != nil {
		t.Fatal(err)
	}
	if !allowed || atype != "ATTENDING" {
		t.Errorf("vínculo esperado allowed=true type=ATTENDING, veio %v/%q", allowed, atype)
	}

	allowed, _, err = r.CheckAssignment(ctx, "med.cardoso", "P000015", "medico")
	if err != nil {
		t.Fatal(err)
	}
	if allowed {
		t.Error("med.cardoso não deveria ter vínculo com P000015")
	}
}

func TestGetPatientNotFound(t *testing.T) {
	r := newRepo(t)
	_, err := r.GetPatient(context.Background(), "P999999")
	if err != repository.ErrNotFound {
		t.Errorf("esperado ErrNotFound, veio %v", err)
	}
}

func sumCounts(t *testing.T, counts []domain.Count, err error) int64 {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	var total int64
	for _, c := range counts {
		total += c.Count
	}
	return total
}
