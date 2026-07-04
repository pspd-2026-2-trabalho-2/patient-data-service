package grpcserver

import (
	"testing"
	"time"

	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/domain"
)

func TestToPBPatient(t *testing.T) {
	p := domain.Patient{
		PatientID: "P000001",
		FullName:  "João da Silva",
		BirthDate: time.Date(1970, 5, 10, 0, 0, 0, 0, time.UTC),
		Gender:    "male",
		City:      "Brasilia",
		State:     "DF",
		CPF:       "111",
		CNS:       "700",
	}
	got := toPBPatient(p)
	if got.GetPatientId() != "P000001" {
		t.Errorf("patient_id = %q", got.GetPatientId())
	}
	if got.GetBirthDate() != "1970-05-10" {
		t.Errorf("birth_date = %q, quer 1970-05-10", got.GetBirthDate())
	}
	if got.GetCpf() != "111" || got.GetCns() != "700" {
		t.Errorf("cpf/cns não mapeados: %q/%q", got.GetCpf(), got.GetCns())
	}
}

func TestToPBEventValue(t *testing.T) {
	v := 8.1
	withValue := toPBEvent(domain.ClinicalEvent{EventID: "E1", Value: &v})
	if withValue.Value == nil || *withValue.Value != 8.1 {
		t.Errorf("valor não mapeado: %v", withValue.Value)
	}

	withoutValue := toPBEvent(domain.ClinicalEvent{EventID: "E2"})
	if withoutValue.Value != nil {
		t.Errorf("valor deveria ser nil, veio %v", *withoutValue.Value)
	}
}

func TestFmtDatePtr(t *testing.T) {
	if got := fmtDatePtr(nil); got != "" {
		t.Errorf("nil deveria virar string vazia, veio %q", got)
	}
	tm := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	if got := fmtDatePtr(&tm); got != "2024-08-01" {
		t.Errorf("data = %q, quer 2024-08-01", got)
	}
}

func TestToPBCohortStatistics(t *testing.T) {
	stats := &domain.CohortStatistics{
		ConditionCode: "Diabetes",
		TotalPatients: 13,
		BySex:         []domain.Count{{Key: "female", Count: 7}, {Key: "male", Count: 6}},
		MeanHbA1c:     7.76,
		MedianHbA1c:   7.6,
	}
	got := toPBCohortStatistics(stats)
	if got.GetTotalPatients() != 13 {
		t.Errorf("total = %d", got.GetTotalPatients())
	}
	if len(got.GetBySex()) != 2 {
		t.Fatalf("by_sex len = %d", len(got.GetBySex()))
	}
	if got.GetMedianHba1C() != 7.6 {
		t.Errorf("mediana = %v", got.GetMedianHba1C())
	}
}
