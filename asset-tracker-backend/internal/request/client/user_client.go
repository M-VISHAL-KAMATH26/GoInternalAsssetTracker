package client

import (
	"context"
	"errors"
	"fmt"

	"asset-backend/internal/user/domain"
	pb "asset-backend/proto/user"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var ErrEmployeeNotFound = errors.New("employee not found")

// UserClient is the interface request-service's handlers depend on —
// same "accept interfaces" pattern as the repositories, so it can be
// faked in tests without a real gRPC connection.
type UserClient interface {
	GetEmployee(ctx context.Context, employeeID uuid.UUID) (*domain.Employee, error)
	GetManagerOf(ctx context.Context, employeeID uuid.UUID) (*domain.Employee, error)
}

type userClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserClient dials a connection to user-service's gRPC address.
func NewUserClient(address string) (*userClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial user-service: %w", err)
	}
	return &userClient{
		client: pb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *userClient) Close() error {
	return c.conn.Close()
}

func (c *userClient) GetEmployee(ctx context.Context, employeeID uuid.UUID) (*domain.Employee, error) {
	resp, err := c.client.GetEmployee(ctx, &pb.GetEmployeeRequest{EmployeeId: employeeID.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("get employee grpc call failed: %w", err)
	}
	return fromProtoEmployee(resp)
}

func (c *userClient) GetManagerOf(ctx context.Context, employeeID uuid.UUID) (*domain.Employee, error) {
	resp, err := c.client.GetManagerOf(ctx, &pb.GetManagerOfRequest{EmployeeId: employeeID.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("get manager grpc call failed: %w", err)
	}
	return fromProtoEmployee(resp)
}

func fromProtoEmployee(resp *pb.EmployeeResponse) (*domain.Employee, error) {
	id, err := uuid.Parse(resp.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id in response: %w", err)
	}

	var managerID *uuid.UUID
	if resp.ManagerId != "" {
		mid, err := uuid.Parse(resp.ManagerId)
		if err != nil {
			return nil, fmt.Errorf("invalid manager id in response: %w", err)
		}
		managerID = &mid
	}

	return &domain.Employee{
		ID:         id,
		Name:       resp.Name,
		Email:      resp.Email,
		Role:       domain.EmployeeRole(resp.Role),
		Department: resp.Department,
		ManagerID:  managerID,
	}, nil
}