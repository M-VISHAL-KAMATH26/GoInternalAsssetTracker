package client

import (
	"context"
	"errors"
	"fmt"

	pb "asset-backend/proto/inventory"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var ErrAssetUnavailable = errors.New("no asset available")

// InventoryClient is the interface request-service's handlers depend
// on — mirrors UserClient's shape so both gRPC clients follow the same
// convention.
type InventoryClient interface {
	CheckAvailability(ctx context.Context, assetType, category string) (bool, error)
	ReserveAsset(ctx context.Context, assetType, category string, employeeID uuid.UUID) (assetID uuid.UUID, serialNumber string, err error)
}

type inventoryClient struct {
	client pb.InventoryServiceClient
	conn   *grpc.ClientConn
}

func NewInventoryClient(address string) (*inventoryClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial inventory-service: %w", err)
	}
	return &inventoryClient{
		client: pb.NewInventoryServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *inventoryClient) Close() error {
	return c.conn.Close()
}

func (c *inventoryClient) CheckAvailability(ctx context.Context, assetType, category string) (bool, error) {
	resp, err := c.client.CheckAvailability(ctx, &pb.CheckAvailabilityRequest{
		AssetType: assetType,
		Category:  category,
	})
	if err != nil {
		return false, fmt.Errorf("check availability grpc call failed: %w", err)
	}
	return resp.Available, nil
}

func (c *inventoryClient) ReserveAsset(ctx context.Context, assetType, category string, employeeID uuid.UUID) (uuid.UUID, string, error) {
	resp, err := c.client.ReserveAsset(ctx, &pb.ReserveAssetRequest{
		AssetType:  assetType,
		Category:   category,
		EmployeeId: employeeID.String(),
	})
	if err != nil {
		if status.Code(err) == codes.FailedPrecondition {
			return uuid.UUID{}, "", ErrAssetUnavailable
		}
		return uuid.UUID{}, "", fmt.Errorf("reserve asset grpc call failed: %w", err)
	}

	assetID, err := uuid.Parse(resp.AssetId)
	if err != nil {
		return uuid.UUID{}, "", fmt.Errorf("invalid asset id in response: %w", err)
	}

	return assetID, resp.SerialNumber, nil
}