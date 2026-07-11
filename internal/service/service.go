// Package service aplica as regras de negócio sobre o repositório.
package service

import (
	"context"

	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/domain"
	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/repository"
)

const maxRecentObservations = 5

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) PatientsByDoctor(ctx context.Context, doctor string, yield func(domain.Patient) error) error {
	return s.repo.PatientsByDoctor(ctx, doctor, yield)
}

func (s *Service) SupervisedPatients(ctx context.Context, intern string, yield func(domain.Patient) error) error {
	return s.repo.SupervisedPatients(ctx, intern, yield)
}

func (s *Service) GetPatient(ctx context.Context, patientID string) (*domain.Patient, error) {
	return s.repo.GetPatient(ctx, patientID)
}

func (s *Service) Encounters(ctx context.Context, patientID string) ([]domain.Encounter, error) {
	return s.repo.Encounters(ctx, patientID)
}

func (s *Service) ClinicalEvents(ctx context.Context, patientID, eventType string) ([]domain.ClinicalEvent, error) {
	return s.repo.ClinicalEvents(ctx, patientID, eventType)
}

func (s *Service) ClinicalHistory(ctx context.Context, patientID string) ([]domain.ClinicalEvent, error) {
	return s.repo.ClinicalHistory(ctx, patientID)
}

func (s *Service) CohortPatients(ctx context.Context, conditionCode string, yield func(domain.Patient) error) error {
	return s.repo.CohortPatients(ctx, conditionCode, yield)
}

func (s *Service) ProjectsByResearcher(ctx context.Context, researcher string) ([]domain.Project, error) {
	return s.repo.ProjectsByResearcher(ctx, researcher)
}

func (s *Service) CheckAssignment(ctx context.Context, username, patientID, role string) (bool, string, error) {
	return s.repo.CheckAssignment(ctx, username, patientID, role)
}

// GetClinicalSummary monta paciente + último atendimento + condições + exames recentes + medicações.
func (s *Service) GetClinicalSummary(ctx context.Context, patientID string) (*domain.ClinicalSummary, error) {
	patient, err := s.repo.GetPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}
	encounters, err := s.repo.Encounters(ctx, patientID)
	if err != nil {
		return nil, err
	}
	conditions, err := s.repo.ClinicalEvents(ctx, patientID, domain.EventTypeCondition)
	if err != nil {
		return nil, err
	}
	observations, err := s.repo.ClinicalEvents(ctx, patientID, domain.EventTypeObservation)
	if err != nil {
		return nil, err
	}
	medications, err := s.repo.ClinicalEvents(ctx, patientID, domain.EventTypeMedication)
	if err != nil {
		return nil, err
	}

	summary := &domain.ClinicalSummary{
		Patient:            *patient,
		Conditions:         conditions,
		RecentObservations: firstN(observations, maxRecentObservations),
		ActiveMedications:  medications,
	}
	if len(encounters) > 0 {
		summary.LastEncounter = &encounters[0] // Encounters já vem ordenado por data desc
	}
	return summary, nil
}

// GetCohortStatistics reúne todas as agregações da coorte numa única resposta.
func (s *Service) GetCohortStatistics(ctx context.Context, conditionCode string) (*domain.CohortStatistics, error) {
	total, err := s.repo.CohortTotal(ctx, conditionCode)
	if err != nil {
		return nil, err
	}
	bySex, err := s.repo.CohortBySex(ctx, conditionCode)
	if err != nil {
		return nil, err
	}
	byAge, err := s.repo.CohortByAgeRange(ctx, conditionCode)
	if err != nil {
		return nil, err
	}
	mean, median, err := s.repo.CohortHbA1c(ctx, conditionCode)
	if err != nil {
		return nil, err
	}
	meds, err := s.repo.CohortMedicationFrequency(ctx, conditionCode)
	if err != nil {
		return nil, err
	}
	byDept, err := s.repo.CohortByDepartment(ctx, conditionCode)
	if err != nil {
		return nil, err
	}

	return &domain.CohortStatistics{
		ConditionCode:       conditionCode,
		TotalPatients:       total,
		BySex:               bySex,
		ByAgeRange:          byAge,
		MeanHbA1c:           mean,
		MedianHbA1c:         median,
		MedicationFrequency: meds,
		ByDepartment:        byDept,
	}, nil
}

func firstN[T any](s []T, n int) []T {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
