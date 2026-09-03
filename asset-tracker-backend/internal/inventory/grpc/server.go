package grpc

import (
	"context"
	"errors"

	"asset-backend/internal/inventory/repository"
	pb "asset-backend/proto/inventory"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedInventoryServiceServer
	assetRepo repository.AssetRepository
}

func NewServer(assetRepo repository.AssetRepository) *Server {
	return &Server{assetRepo: assetRepo}
}

func (s *Server) CheckAvailability(ctx context.Context, req *pb.CheckAvailabilityRequest) (*pb.CheckAvailabilityResponse, error) {
	assets, err := s.assetRepo.ListByStatus(ctx, "available")
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check availability")
	}

	count := 0
	for _, a := range assets {
		if a.Type == req.AssetType && a.Category == req.Category {
			count++
		}
	}

	return &pb.CheckAvailabilityResponse{
		Available:      count > 0,
		AvailableCount: int32(count),
	}, nil
}

func (s *Server) ReserveAsset(ctx context.Context, req *pb.ReserveAssetRequest) (*pb.ReserveAssetResponse, error) {
	employeeID, err := uuid.Parse(req.EmployeeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid employee id")
	}

	asset, err := s.assetRepo.ReserveAvailableAsset(ctx, req.AssetType, req.Category, employeeID)
	if err != nil {
		if errors.Is(err, repository.ErrNoAssetAvailable) {
			return nil, status.Error(codes.FailedPrecondition, "no asset available")
		}
		return nil, status.Error(codes.Internal, "failed to reserve asset")
	}

	return &pb.ReserveAssetResponse{
		AssetId:      asset.ID.String(),
		SerialNumber: asset.SerialNumber,
	}, nil
}