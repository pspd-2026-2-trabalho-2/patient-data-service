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
	err := r.PatientsByDoctor(context.Background(), "med.cardoso", 100, 0, "", "", func(p domain.Patient) error {
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
	err := r.SupervisedPatients(context.Background(), "est.souza", 100, 0, "", "", func(p domain.Patient) error {
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

// TestPatientsByDoctorSearchByName: "Silva" só bate com P000001 (João da Silva).
func TestPatientsByDoctorSearchByName(t *testing.T) {
	r := newRepo(t)
	var ps []domain.Patient
	err := r.PatientsByDoctor(context.Background(), "med.cardoso", 100, 0, "silva", "", func(p domain.Patient) error {
		ps = append(ps, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 1 || ps[0].PatientID != "P000001" {
		t.Errorf("busca por 'silva' deveria trazer só P000001, veio %+v", ps)
	}
}

// TestPatientsByDoctorSearchByCPF: "444.444" só bate com o CPF de P000004.
func TestPatientsByDoctorSearchByCPF(t *testing.T) {
	r := newRepo(t)
	var ps []domain.Patient
	err := r.PatientsByDoctor(context.Background(), "med.cardoso", 100, 0, "444.444", "", func(p domain.Patient) error {
		ps = append(ps, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 1 || ps[0].PatientID != "P000004" {
		t.Errorf("busca por CPF '444.444' deveria trazer só P000004, veio %+v", ps)
	}
}

// TestPatientsByDoctorSearchNoResults: busca sem correspondência devolve bundle vazio, não erro.
func TestPatientsByDoctorSearchNoResults(t *testing.T) {
	r := newRepo(t)
	var ps []domain.Patient
	err := r.PatientsByDoctor(context.Background(), "med.cardoso", 100, 0, "zzz-nao-existe-zzz", "", func(p domain.Patient) error {
		ps = append(ps, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 0 {
		t.Errorf("busca sem correspondência deveria trazer 0 pacientes, veio %d", len(ps))
	}
}

// TestPatientsByDoctorSearchWithPagination: busca por "e" bate com P000003, P000004,
// P000005, P000006, P000008 (Oliveira, Pereira, Beatriz, Pedro, Rafael/Almeida) — 5 no
// total entre os pacientes de med.cardoso. Pagina de 3 em 3.
func TestPatientsByDoctorSearchWithPagination(t *testing.T) {
	r := newRepo(t)

	var page1 []domain.Patient
	err := r.PatientsByDoctor(context.Background(), "med.cardoso", 3, 0, "e", "", func(p domain.Patient) error {
		page1 = append(page1, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	wantPage1 := []string{"P000003", "P000004", "P000005"}
	if !samePatientIDs(page1, wantPage1) {
		t.Errorf("página 1 da busca 'e' = %v, quer %v", ids(page1), wantPage1)
	}

	var page2 []domain.Patient
	err = r.PatientsByDoctor(context.Background(), "med.cardoso", 3, 3, "e", "", func(p domain.Patient) error {
		page2 = append(page2, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	wantPage2 := []string{"P000006", "P000008"}
	if !samePatientIDs(page2, wantPage2) {
		t.Errorf("página 2 da busca 'e' = %v, quer %v", ids(page2), wantPage2)
	}
}

// TestPatientsByDoctorGenderFilter: gender=male entre os pacientes de med.cardoso
// traz P000001, P000004, P000006, P000008.
func TestPatientsByDoctorGenderFilter(t *testing.T) {
	r := newRepo(t)
	var ps []domain.Patient
	err := r.PatientsByDoctor(context.Background(), "med.cardoso", 100, 0, "", "male", func(p domain.Patient) error {
		ps = append(ps, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"P000001", "P000004", "P000006", "P000008"}
	if !samePatientIDs(ps, want) {
		t.Errorf("filtro gender=male = %v, quer %v", ids(ps), want)
	}
}

// TestSupervisedPatientsSearchByName: est.souza supervisiona P000001..P000003;
// "Souza" só bate com P000002 (Maria Souza).
func TestSupervisedPatientsSearchByName(t *testing.T) {
	r := newRepo(t)
	var ps []domain.Patient
	err := r.SupervisedPatients(context.Background(), "est.souza", 100, 0, "souza", "", func(p domain.Patient) error {
		ps = append(ps, p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 1 || ps[0].PatientID != "P000002" {
		t.Errorf("busca por 'souza' deveria trazer só P000002, veio %+v", ps)
	}
}

func ids(ps []domain.Patient) []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = p.PatientID
	}
	return out
}

func samePatientIDs(ps []domain.Patient, want []string) bool {
	got := ids(ps)
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
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
