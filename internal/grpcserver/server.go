// Package grpcserver implementa a camada de transporte gRPC.
package grpcserver

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pspd-2026-2-trabalho-2/patient-data-service/gen/patientdata/v1"
	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/domain"
	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/repository"
	"github.com/pspd-2026-2-trabalho-2/patient-data-service/internal/service"
)

// Server implementa pb.PatientDataServiceServer.
type Server struct {
	pb.UnimplementedPatientDataServiceServer
	svc *service.Service
}

// New cria o servidor gRPC.
func New(svc *service.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) ListPatientsByDoctor(req *pb.ListPatientsByDoctorRequest, stream pb.PatientDataService_ListPatientsByDoctorServer) error {
	if req.GetDoctorUsername() == "" {
		return status.Error(codes.InvalidArgument, "doctor_username é obrigatório")
	}
	if err := s.svc.PatientsByDoctor(stream.Context(), req.GetDoctorUsername(), int(req.GetPage()), int(req.GetPageSize()), req.GetSearch(), req.GetGender(), func(p domain.Patient) error {
		return stream.Send(toPBPatient(p))
	}); err != nil {
		return toStatus(err)
	}
	return nil
}

func (s *Server) ListSupervisedPatients(req *pb.ListSupervisedPatientsRequest, stream pb.PatientDataService_ListSupervisedPatientsServer) error {
	if req.GetInternUsername() == "" {
		return status.Error(codes.InvalidArgument, "intern_username é obrigatório")
	}
	if err := s.svc.SupervisedPatients(stream.Context(), req.GetInternUsername(), int(req.GetPage()), int(req.GetPageSize()), req.GetSearch(), req.GetGender(), func(p domain.Patient) error {
		return stream.Send(toPBPatient(p))
	}); err != nil {
		return toStatus(err)
	}
	return nil
}

func (s *Server) GetPatient(ctx context.Context, req *pb.GetPatientRequest) (*pb.Patient, error) {
	if req.GetPatientId() == "" {
		return nil, status.Error(codes.InvalidArgument, "patient_id é obrigatório")
	}
	patient, err := s.svc.GetPatient(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatus(err)
	}
	return toPBPatient(*patient), nil
}

func (s *Server) ListEncounters(ctx context.Context, req *pb.ListEncountersRequest) (*pb.EncounterList, error) {
	if req.GetPatientId() == "" {
		return nil, status.Error(codes.InvalidArgument, "patient_id é obrigatório")
	}
	encounters, err := s.svc.Encounters(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatus(err)
	}
	return &pb.EncounterList{Encounters: toPBEncounters(encounters)}, nil
}

func (s *Server) ListClinicalEvents(ctx context.Context, req *pb.ListClinicalEventsRequest) (*pb.ClinicalEventList, error) {
	if req.GetPatientId() == "" {
		return nil, status.Error(codes.InvalidArgument, "patient_id é obrigatório")
	}
	events, err := s.svc.ClinicalEvents(ctx, req.GetPatientId(), req.GetEventType())
	if err != nil {
		return nil, toStatus(err)
	}
	return &pb.ClinicalEventList{Events: toPBEvents(events)}, nil
}

func (s *Server) GetClinicalSummary(ctx context.Context, req *pb.GetClinicalSummaryRequest) (*pb.ClinicalSummary, error) {
	if req.GetPatientId() == "" {
		return nil, status.Error(codes.InvalidArgument, "patient_id é obrigatório")
	}
	summary, err := s.svc.GetClinicalSummary(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatus(err)
	}
	return toPBSummary(summary), nil
}

func (s *Server) GetClinicalHistory(ctx context.Context, req *pb.GetClinicalHistoryRequest) (*pb.ClinicalEventList, error) {
	if req.GetPatientId() == "" {
		return nil, status.Error(codes.InvalidArgument, "patient_id é obrigatório")
	}
	events, err := s.svc.ClinicalHistory(ctx, req.GetPatientId())
	if err != nil {
		return nil, toStatus(err)
	}
	return &pb.ClinicalEventList{Events: toPBEvents(events)}, nil
}

func (s *Server) ListCohortPatients(req *pb.ListCohortPatientsRequest, stream pb.PatientDataService_ListCohortPatientsServer) error {
	if req.GetConditionCode() == "" {
		return status.Error(codes.InvalidArgument, "condition_code é obrigatório")
	}
	if err := s.svc.CohortPatients(stream.Context(), req.GetConditionCode(), func(p domain.Patient) error {
		return stream.Send(toPBPatient(p))
	}); err != nil {
		return toStatus(err)
	}
	return nil
}

func (s *Server) ListCohortExams(ctx context.Context, req *pb.ListCohortExamsRequest) (*pb.CohortExamsList, error) {
	if req.GetConditionCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "condition_code é obrigatório")
	}
	patientExams, err := s.svc.CohortExams(ctx, req.GetConditionCode(), int(req.GetPage()), int(req.GetPageSize()), req.GetEventType())
	if err != nil {
		return nil, toStatus(err)
	}
	return &pb.CohortExamsList{Patients: toPBPatientExamsList(patientExams)}, nil
}

func (s *Server) GetCohortStatistics(ctx context.Context, req *pb.GetCohortStatisticsRequest) (*pb.CohortStatistics, error) {
	if req.GetConditionCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "condition_code é obrigatório")
	}
	stats, err := s.svc.GetCohortStatistics(ctx, req.GetConditionCode())
	if err != nil {
		return nil, toStatus(err)
	}
	return toPBCohortStatistics(stats), nil
}

func (s *Server) ListProjectsByResearcher(ctx context.Context, req *pb.ListProjectsByResearcherRequest) (*pb.ProjectList, error) {
	if req.GetResearcherUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "researcher_username é obrigatório")
	}
	projects, err := s.svc.ProjectsByResearcher(ctx, req.GetResearcherUsername())
	if err != nil {
		return nil, toStatus(err)
	}
	return &pb.ProjectList{Projects: toPBProjects(projects)}, nil
}

func (s *Server) CheckAssignment(ctx context.Context, req *pb.CheckAssignmentRequest) (*pb.CheckAssignmentResponse, error) {
	if req.GetUsername() == "" || req.GetPatientId() == "" {
		return nil, status.Error(codes.InvalidArgument, "username e patient_id são obrigatórios")
	}
	allowed, assignmentType, err := s.svc.CheckAssignment(ctx, req.GetUsername(), req.GetPatientId(), req.GetRole())
	if err != nil {
		return nil, toStatus(err)
	}
	return &pb.CheckAssignmentResponse{Allowed: allowed, AssignmentType: assignmentType}, nil
}

// toStatus converte erros de domínio em status gRPC apropriados.
func toStatus(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
