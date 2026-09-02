package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"asset-backend/internal/inventory/domain"
	"asset-backend/internal/inventory/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// fakeAssetRepo is a hand-rolled test double for AssetRepository.
// Each method delegates to an optional func field so tests override only
// the behavior they need.
type fakeAssetRepo struct {
	CreateFunc            func(ctx context.Context, asset *domain.Asset) error
	GetByIDFunc           func(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
	GetBySerialNumberFunc func(ctx context.Context, serial string) (*domain.Asset, error)
	ListByStatusFunc      func(ctx context.Context, status domain.AssetStatus) ([]domain.Asset, error)
	ListFunc              func(ctx context.Context) ([]domain.Asset, error)
	UpdateStatusFunc      func(ctx context.Context, id uuid.UUID, status domain.AssetStatus) error
	UpdateFunc            func(ctx context.Context, asset *domain.Asset) error
}

func (f *fakeAssetRepo) Create(ctx context.Context, asset *domain.Asset) error {
	if f.CreateFunc != nil {
		return f.CreateFunc(ctx, asset)
	}
	return nil
}

func (f *fakeAssetRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	if f.GetByIDFunc != nil {
		return f.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (f *fakeAssetRepo) GetBySerialNumber(ctx context.Context, serial string) (*domain.Asset, error) {
	if f.GetBySerialNumberFunc != nil {
		return f.GetBySerialNumberFunc(ctx, serial)
	}
	return nil, nil
}

func (f *fakeAssetRepo) ListByStatus(ctx context.Context, status domain.AssetStatus) ([]domain.Asset, error) {
	if f.ListByStatusFunc != nil {
		return f.ListByStatusFunc(ctx, status)
	}
	return nil, nil
}

func (f *fakeAssetRepo) List(ctx context.Context) ([]domain.Asset, error) {
	if f.ListFunc != nil {
		return f.ListFunc(ctx)
	}
	return nil, nil
}

func (f *fakeAssetRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AssetStatus) error {
	if f.UpdateStatusFunc != nil {
		return f.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}

func (f *fakeAssetRepo) Update(ctx context.Context, asset *domain.Asset) error {
	if f.UpdateFunc != nil {
		return f.UpdateFunc(ctx, asset)
	}
	return nil
}

func newTestContext(method, target string, body []byte, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var reader *bytes.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	} else {
		reader = bytes.NewReader(nil)
	}

	c.Request = httptest.NewRequest(method, target, reader)
	if body != nil {
		c.Request.Header.Set("Content-Type", "application/json")
	}
	c.Params = params
	return c, w
}

func sampleAsset(id uuid.UUID) *domain.Asset {
	now := time.Now()
	return &domain.Asset{
		ID:           id,
		Name:         "MacBook Pro",
		Type:         "laptop",
		Category:     "hardware",
		SerialNumber: "SN-001",
		Status:       domain.AssetStatusAvailable,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func TestCreateAsset(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repo       *fakeAssetRepo
		wantStatus int
		checkBody  func(t *testing.T, body []byte)
	}{
		{
			name: "valid request succeeds",
			body: `{"name":"MacBook Pro","type":"laptop","category":"hardware","serial_number":"SN-001"}`,
			repo: &fakeAssetRepo{
				CreateFunc: func(_ context.Context, asset *domain.Asset) error {
					asset.CreatedAt = time.Now()
					asset.UpdatedAt = time.Now()
					return nil
				},
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body []byte) {
				var resp AssetResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.Name != "MacBook Pro" {
					t.Errorf("name = %q, want %q", resp.Name, "MacBook Pro")
				}
				if resp.SerialNumber != "SN-001" {
					t.Errorf("serial_number = %q, want %q", resp.SerialNumber, "SN-001")
				}
				if resp.Status != domain.AssetStatusAvailable {
					t.Errorf("status = %q, want %q", resp.Status, domain.AssetStatusAvailable)
				}
				if resp.ID == uuid.Nil {
					t.Error("expected non-nil id in response")
				}
			},
		},
		{
			name:       "missing required field fails validation",
			body:       `{"type":"laptop","category":"hardware","serial_number":"SN-001"}`,
			repo:       &fakeAssetRepo{},
			wantStatus: http.StatusBadRequest,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] == "" {
					t.Error("expected error message in response")
				}
			},
		},
		{
			name: "repository returns error",
			body: `{"name":"MacBook Pro","type":"laptop","category":"hardware","serial_number":"SN-001"}`,
			repo: &fakeAssetRepo{
				CreateFunc: func(_ context.Context, _ *domain.Asset) error {
					return errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "failed to create asset" {
					t.Errorf("error = %q, want %q", resp["error"], "failed to create asset")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewAssetHandler(tt.repo)
			c, w := newTestContext(http.MethodPost, "/assets", []byte(tt.body), nil)

			h.CreateAsset(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
			tt.checkBody(t, w.Body.Bytes())
		})
	}
}

func TestListAssets(t *testing.T) {
	assetID := uuid.New()
	assets := []domain.Asset{*sampleAsset(assetID)}

	tests := []struct {
		name       string
		target     string
		repo       *fakeAssetRepo
		wantStatus int
		checkBody  func(t *testing.T, body []byte)
	}{
		{
			name:   "no filter returns all assets",
			target: "/assets",
			repo: &fakeAssetRepo{
				ListFunc: func(_ context.Context) ([]domain.Asset, error) {
					return assets, nil
				},
				ListByStatusFunc: func(_ context.Context, _ domain.AssetStatus) ([]domain.Asset, error) {
					t.Error("ListByStatus should not be called when no status filter is provided")
					return nil, nil
				},
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp []AssetResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if len(resp) != 1 {
					t.Fatalf("len(assets) = %d, want 1", len(resp))
				}
				if resp[0].Name != "MacBook Pro" {
					t.Errorf("name = %q, want %q", resp[0].Name, "MacBook Pro")
				}
			},
		},
		{
			name:   "status filter calls ListByStatus",
			target: "/assets?status=available",
			repo: &fakeAssetRepo{
				ListFunc: func(_ context.Context) ([]domain.Asset, error) {
					t.Error("List should not be called when status filter is provided")
					return nil, nil
				},
				ListByStatusFunc: func(_ context.Context, status domain.AssetStatus) ([]domain.Asset, error) {
					if status != domain.AssetStatusAvailable {
						t.Errorf("status = %q, want %q", status, domain.AssetStatusAvailable)
					}
					return assets, nil
				},
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp []AssetResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if len(resp) != 1 {
					t.Fatalf("len(assets) = %d, want 1", len(resp))
				}
				if resp[0].Status != domain.AssetStatusAvailable {
					t.Errorf("status = %q, want %q", resp[0].Status, domain.AssetStatusAvailable)
				}
			},
		},
		{
			name:   "repository error",
			target: "/assets",
			repo: &fakeAssetRepo{
				ListFunc: func(_ context.Context) ([]domain.Asset, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "failed to list assets" {
					t.Errorf("error = %q, want %q", resp["error"], "failed to list assets")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewAssetHandler(tt.repo)
			c, w := newTestContext(http.MethodGet, tt.target, nil, nil)

			h.ListAssets(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
			tt.checkBody(t, w.Body.Bytes())
		})
	}
}

func TestGetAsset(t *testing.T) {
	assetID := uuid.New()

	tests := []struct {
		name       string
		idParam    string
		repo       *fakeAssetRepo
		wantStatus int
		checkBody  func(t *testing.T, body []byte)
	}{
		{
			name:    "valid ID found",
			idParam: assetID.String(),
			repo: &fakeAssetRepo{
				GetByIDFunc: func(_ context.Context, id uuid.UUID) (*domain.Asset, error) {
					if id != assetID {
						t.Errorf("id = %v, want %v", id, assetID)
					}
					return sampleAsset(assetID), nil
				},
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp AssetResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.ID != assetID {
					t.Errorf("id = %v, want %v", resp.ID, assetID)
				}
				if resp.Name != "MacBook Pro" {
					t.Errorf("name = %q, want %q", resp.Name, "MacBook Pro")
				}
			},
		},
		{
			name:       "invalid UUID format",
			idParam:    "not-a-uuid",
			repo:       &fakeAssetRepo{},
			wantStatus: http.StatusBadRequest,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "invalid asset id" {
					t.Errorf("error = %q, want %q", resp["error"], "invalid asset id")
				}
			},
		},
		{
			name:    "valid UUID but not found",
			idParam: assetID.String(),
			repo: &fakeAssetRepo{
				GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.Asset, error) {
					return nil, repository.ErrAssetNotFound
				},
			},
			wantStatus: http.StatusNotFound,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "asset not found" {
					t.Errorf("error = %q, want %q", resp["error"], "asset not found")
				}
			},
		},
		{
			name:    "other repository error",
			idParam: assetID.String(),
			repo: &fakeAssetRepo{
				GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.Asset, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "failed to get asset" {
					t.Errorf("error = %q, want %q", resp["error"], "failed to get asset")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewAssetHandler(tt.repo)
			c, w := newTestContext(http.MethodGet, "/assets/"+tt.idParam, nil, gin.Params{
				{Key: "id", Value: tt.idParam},
			})

			h.GetAsset(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
			tt.checkBody(t, w.Body.Bytes())
		})
	}
}

func TestUpdateAsset(t *testing.T) {
	assetID := uuid.New()
	validBody := `{"name":"Updated Laptop","type":"laptop","category":"hardware","serial_number":"SN-002","status":"assigned"}`

	tests := []struct {
		name       string
		idParam    string
		body       string
		repo       *fakeAssetRepo
		wantStatus int
		checkBody  func(t *testing.T, body []byte)
	}{
		{
			name:    "valid update succeeds",
			idParam: assetID.String(),
			body:    validBody,
			repo: &fakeAssetRepo{
				GetByIDFunc: func(_ context.Context, id uuid.UUID) (*domain.Asset, error) {
					return sampleAsset(id), nil
				},
				UpdateFunc: func(_ context.Context, asset *domain.Asset) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp AssetResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.Name != "Updated Laptop" {
					t.Errorf("name = %q, want %q", resp.Name, "Updated Laptop")
				}
				if resp.Status != domain.AssetStatusAssigned {
					t.Errorf("status = %q, want %q", resp.Status, domain.AssetStatusAssigned)
				}
				if resp.SerialNumber != "SN-002" {
					t.Errorf("serial_number = %q, want %q", resp.SerialNumber, "SN-002")
				}
			},
		},
		{
			name:       "invalid UUID",
			idParam:    "bad-id",
			body:       validBody,
			repo:       &fakeAssetRepo{},
			wantStatus: http.StatusBadRequest,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "invalid asset id" {
					t.Errorf("error = %q, want %q", resp["error"], "invalid asset id")
				}
			},
		},
		{
			name:    "asset not found",
			idParam: assetID.String(),
			body:    validBody,
			repo: &fakeAssetRepo{
				GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.Asset, error) {
					return nil, repository.ErrAssetNotFound
				},
			},
			wantStatus: http.StatusNotFound,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "asset not found" {
					t.Errorf("error = %q, want %q", resp["error"], "asset not found")
				}
			},
		},
		{
			name:    "invalid status value fails validation",
			idParam: assetID.String(),
			body:    `{"name":"Updated Laptop","type":"laptop","category":"hardware","serial_number":"SN-002","status":"bogus"}`,
			repo: &fakeAssetRepo{
				GetByIDFunc: func(_ context.Context, id uuid.UUID) (*domain.Asset, error) {
					return sampleAsset(id), nil
				},
			},
			wantStatus: http.StatusBadRequest,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] == "" {
					t.Error("expected validation error message in response")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewAssetHandler(tt.repo)
			c, w := newTestContext(http.MethodPut, "/assets/"+tt.idParam, []byte(tt.body), gin.Params{
				{Key: "id", Value: tt.idParam},
			})

			h.UpdateAsset(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
			tt.checkBody(t, w.Body.Bytes())
		})
	}
}

func TestRetireAsset(t *testing.T) {
	assetID := uuid.New()

	tests := []struct {
		name       string
		idParam    string
		repo       *fakeAssetRepo
		wantStatus int
		checkBody  func(t *testing.T, body []byte)
	}{
		{
			name:    "valid ID succeeds",
			idParam: assetID.String(),
			repo: &fakeAssetRepo{
				GetByIDFunc: func(_ context.Context, id uuid.UUID) (*domain.Asset, error) {
					return sampleAsset(id), nil
				},
				UpdateStatusFunc: func(_ context.Context, id uuid.UUID, status domain.AssetStatus) error {
					if id != assetID {
						t.Errorf("id = %v, want %v", id, assetID)
					}
					if status != domain.AssetStatusRetired {
						t.Errorf("status = %q, want %q", status, domain.AssetStatusRetired)
					}
					return nil
				},
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]any
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["message"] != "asset retired" {
					t.Errorf("message = %v, want %q", resp["message"], "asset retired")
				}
				if resp["id"] != assetID.String() {
					t.Errorf("id = %v, want %q", resp["id"], assetID.String())
				}
			},
		},
		{
			name:    "asset not found",
			idParam: assetID.String(),
			repo: &fakeAssetRepo{
				GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.Asset, error) {
					return nil, repository.ErrAssetNotFound
				},
			},
			wantStatus: http.StatusNotFound,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "asset not found" {
					t.Errorf("error = %q, want %q", resp["error"], "asset not found")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewAssetHandler(tt.repo)
			c, w := newTestContext(http.MethodPatch, "/assets/"+tt.idParam+"/retire", nil, gin.Params{
				{Key: "id", Value: tt.idParam},
			})

			h.RetireAsset(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
			tt.checkBody(t, w.Body.Bytes())
		})
	}
}
