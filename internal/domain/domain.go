// Package domain define os modelos de negócio (dados crus), sem protobuf nem SQL.
package domain

import "time"

// Tipos de evento clínico (coluna clinical_events.event_type).
const (
	EventTypeCondition   = "Condition"
	EventTypeObservation = "Observation"
	EventTypeMedication  = "Medication"
)

// Patient espelha a tabela patients.
type Patient struct {
	PatientID string
	FullName  string
	BirthDate time.Time
	Gender    string
	City      string
	State     string
	CPF       string
	CNS       string
}

// Encounter espelha a tabela encounters.
type Encounter struct {
	EncounterID   string
	PatientID     string
	StartDate     time.Time
	EndDate       *time.Time // nulo enquanto o atendimento não encerra
	EncounterType string
	Department    string
}

// ClinicalEvent espelha a tabela clinical_events.
type ClinicalEvent struct {
	EventID     string
	PatientID   string
	EncounterID string
	EventType   string
	Code        string
	Description string
	EventDate   time.Time
	Value       *float64 // preenchido para Observation/Medication
	Unit        string
}

// Project espelha a tabela projects.
type Project struct {
	ProjectID          string
	Title              string
	ResearcherUsername string
	ConditionCode      string
	Status             string
	ValidUntil         *time.Time
}

// ClinicalSummary é o resumo clínico montado a partir de várias tabelas.
type ClinicalSummary struct {
	Patient            Patient
	LastEncounter      *Encounter
	Conditions         []ClinicalEvent
	RecentObservations []ClinicalEvent
	ActiveMedications  []ClinicalEvent
}

// Count é um par chave→contagem usado em distribuições e frequências.
type Count struct {
	Key   string
	Count int64
}

// CohortStatistics reúne as agregações de uma coorte.
type CohortStatistics struct {
	ConditionCode       string
	TotalPatients       int64
	BySex               []Count
	ByAgeRange          []Count
	MeanHbA1c           float64
	MedianHbA1c         float64
	MedicationFrequency []Count
}
