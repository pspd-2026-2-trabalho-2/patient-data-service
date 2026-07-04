package grpcserver

import (
	"time"

	pb "github.com/pspd-2026-2-trabalho-2/patient-data-service/gen/patientdata/v1"
	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/domain"
)

// dateLayout é o formato ISO usado nas datas expostas em texto (as colunas são DATE).
const dateLayout = "2006-01-02"

func fmtDate(t time.Time) string { return t.Format(dateLayout) }

func fmtDatePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(dateLayout)
}

func toPBPatient(p domain.Patient) *pb.Patient {
	return &pb.Patient{
		PatientId: p.PatientID,
		FullName:  p.FullName,
		BirthDate: fmtDate(p.BirthDate),
		Gender:    p.Gender,
		City:      p.City,
		State:     p.State,
		Cpf:       p.CPF,
		Cns:       p.CNS,
	}
}

func toPBPatients(in []domain.Patient) []*pb.Patient {
	out := make([]*pb.Patient, 0, len(in))
	for _, p := range in {
		out = append(out, toPBPatient(p))
	}
	return out
}

func toPBEncounter(e domain.Encounter) *pb.Encounter {
	return &pb.Encounter{
		EncounterId:   e.EncounterID,
		PatientId:     e.PatientID,
		StartDate:     fmtDate(e.StartDate),
		EndDate:       fmtDatePtr(e.EndDate),
		EncounterType: e.EncounterType,
		Department:    e.Department,
	}
}

func toPBEncounters(in []domain.Encounter) []*pb.Encounter {
	out := make([]*pb.Encounter, 0, len(in))
	for _, e := range in {
		out = append(out, toPBEncounter(e))
	}
	return out
}

func toPBEvent(e domain.ClinicalEvent) *pb.ClinicalEvent {
	ev := &pb.ClinicalEvent{
		EventId:     e.EventID,
		PatientId:   e.PatientID,
		EncounterId: e.EncounterID,
		EventType:   e.EventType,
		Code:        e.Code,
		Description: e.Description,
		EventDate:   fmtDate(e.EventDate),
		Unit:        e.Unit,
	}
	if e.Value != nil {
		v := *e.Value
		ev.Value = &v
	}
	return ev
}

func toPBEvents(in []domain.ClinicalEvent) []*pb.ClinicalEvent {
	out := make([]*pb.ClinicalEvent, 0, len(in))
	for _, e := range in {
		out = append(out, toPBEvent(e))
	}
	return out
}

func toPBProject(p domain.Project) *pb.Project {
	return &pb.Project{
		ProjectId:          p.ProjectID,
		Title:              p.Title,
		ResearcherUsername: p.ResearcherUsername,
		ConditionCode:      p.ConditionCode,
		Status:             p.Status,
		ValidUntil:         fmtDatePtr(p.ValidUntil),
	}
}

func toPBProjects(in []domain.Project) []*pb.Project {
	out := make([]*pb.Project, 0, len(in))
	for _, p := range in {
		out = append(out, toPBProject(p))
	}
	return out
}

func toPBCounts(in []domain.Count) []*pb.Count {
	out := make([]*pb.Count, 0, len(in))
	for _, c := range in {
		out = append(out, &pb.Count{Key: c.Key, Count: c.Count})
	}
	return out
}

func toPBSummary(s *domain.ClinicalSummary) *pb.ClinicalSummary {
	out := &pb.ClinicalSummary{
		Patient:            toPBPatient(s.Patient),
		Conditions:         toPBEvents(s.Conditions),
		RecentObservations: toPBEvents(s.RecentObservations),
		ActiveMedications:  toPBEvents(s.ActiveMedications),
	}
	if s.LastEncounter != nil {
		out.LastEncounter = toPBEncounter(*s.LastEncounter)
	}
	return out
}

func toPBCohortStatistics(s *domain.CohortStatistics) *pb.CohortStatistics {
	return &pb.CohortStatistics{
		ConditionCode:       s.ConditionCode,
		TotalPatients:       s.TotalPatients,
		BySex:               toPBCounts(s.BySex),
		ByAgeRange:          toPBCounts(s.ByAgeRange),
		MeanHba1C:           s.MeanHbA1c,
		MedianHba1C:         s.MedianHbA1c,
		MedicationFrequency: toPBCounts(s.MedicationFrequency),
	}
}
