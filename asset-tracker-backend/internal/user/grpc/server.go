package grpc

import (
	"context"
	"errors"

	"asset-backend/internal/user/domain"
	"asset-backend/internal/user/repository"
	pb "asset-backend/proto/user"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the generated pb.UserServiceServer interface.
type Server struct {
	pb.UnimplementedUserServiceServer
	repo repository.EmployeeRepository
}

func NewServer(repo repository.EmployeeRepository) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetEmployee(ctx context.Context, req *pb.GetEmployeeRequest) (*pb.EmployeeResponse, error) {
	id, err := uuid.Parse(req.EmployeeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid employee id")
	}

	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return nil, status.Error(codes.NotFound, "employee not found")
		}
		return nil, status.Error(codes.Internal, "failed to get employee")
	}

	return toProtoEmployee(employee), nil
}

func (s *Server) GetManagerOf(ctx context.Context, req *pb.GetManagerOfRequest) (*pb.EmployeeResponse, error) {
	id, err := uuid.Parse(req.EmployeeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid employee id")
	}

	manager, err := s.repo.GetManagerOf(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return nil, status.Error(codes.NotFound, "manager not found")
		}
		return nil, status.Error(codes.Internal, "failed to get manager")
	}

	return toProtoEmployee(manager), nil
}

func toProtoEmployee(e *domain.Employee) *pb.EmployeeResponse {
	managerID := ""
	if e.ManagerID != nil {
		managerID = e.ManagerID.String()
	}
	return &pb.EmployeeResponse{
		Id:         e.ID.String(),
		Name:       e.Name,
		Email:      e.Email,
		Role:       string(e.Role),
		Department: e.Department,
		ManagerId:  managerID,
	}
}