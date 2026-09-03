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

	"asset-backend/internal/request/domain"
	"asset-backend/internal/request/repository"
	"asset-backend/internal/shared/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// fakeRequestRepo is a hand-rolled test double for RequestRepository.
type fakeRequestRepo struct {
	CreateFunc       func(ctx context.Context, request *domain.AssetRequest) error
	GetByIDFunc      func(ctx context.Context, id uuid.UUID) (*domain.AssetRequest, error)
	ListByEmployeeFunc func(ctx context.Context, employeeID uuid.UUID) ([]domain.AssetRequest, error)
	ListByStatusFunc func(ctx context.Context, status domain.RequestStatus) ([]domain.AssetRequest, error)
	UpdateStatusFunc func(ctx context.Context, id uuid.UUID, status domain.RequestStatus) error
}

func (f *fakeRequestRepo) Create(ctx context.Context, request *domain.AssetRequest) error {
	if f.CreateFunc != nil {
		return f.CreateFunc(ctx, request)
	}
	return nil
}

func (f *fakeRequestRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.AssetRequest, error) {
	if f.GetByIDFunc != nil {
		return f.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (f *fakeRequestRepo) ListByEmployee(ctx context.Context, employeeID uuid.UUID) ([]domain.AssetRequest, error) {
	if f.ListByEmployeeFunc != nil {
		return f.ListByEmployeeFunc(ctx, employeeID)
	}
	return nil, nil
}

func (f *fakeRequestRepo) ListByStatus(ctx context.Context, status domain.RequestStatus) ([]domain.AssetRequest, error) {
	if f.ListByStatusFunc != nil {
		return f.ListByStatusFunc(ctx, status)
	}
	return nil, nil
}

func (f *fakeRequestRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RequestStatus) error {
	if f.UpdateStatusFunc != nil {
		return f.UpdateStatusFunc(ctx, id, status)
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

func setAuthContext(c *gin.Context, employeeID uuid.UUID, role string) {
	c.Set(middleware.ContextKeyEmployeeID, employeeID.String())
	c.Set(middleware.ContextKeyRole, role)
}

func sampleRequest(id, employeeID uuid.UUID) *domain.AssetRequest {
	now := time.Now()
	return &domain.AssetRequest{
		ID:            id,
		EmployeeID:    employeeID,
		AssetType:     "laptop",
		Category:      "hardware",
		Justification: "Need for development work",
		Status:        domain.RequestStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func TestCreateRequest(t *testing.T) {
	employeeID := uuid.New()

	tests := []struct {
		name       string
		body       string
		employeeID uuid.UUID
		setAuth    bool
		repo       *fakeRequestRepo
		wantStatus int
		checkBody  func(t *testing.T, body []byte)
	}{
		{
			name:       "valid request succeeds",
			body:       `{"asset_type":"laptop","category":"hardware","justification":"Need for dev work"}`,
			employeeID: employeeID,
			setAuth:    true,
			repo: &fakeRequestRepo{
				CreateFunc: func(_ context.Context, req *domain.AssetRequest) error {
					if req.EmployeeID != employeeID {
						t.Errorf("employeeID = %v, want %v", req.EmployeeID, employeeID)
					}
					req.CreatedAt = time.Now()
					req.UpdatedAt = time.Now()
					return nil
				},
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body []byte) {
				var resp RequestResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.EmployeeID != employeeID {
					t.Errorf("employee_id = %v, want %v", resp.EmployeeID, employeeID)
				}
				if resp.AssetType != "laptop" {
					t.Errorf("asset_type = %q, want %q", resp.AssetType, "laptop")
				}
				if resp.Status != domain.RequestStatusPending {
					t.Errorf("status = %q, want %q", resp.Status, domain.RequestStatusPending)
				}
				if resp.ID == uuid.Nil {
					t.Error("expected non-nil id in response")
				}
			},
		},
		{
			name:       "missing auth returns unauthorized",
			body:       `{"asset_type":"laptop","category":"hardware"}`,
			setAuth:    false,
			repo:       &fakeRequestRepo{},
			wantStatus: http.StatusUnauthorized,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "invalid employee identity" {
					t.Errorf("error = %q, want %q", resp["error"], "invalid employee identity")
				}
			},
		},
		{
			name:       "missing required field fails validation",
			body:       `{"category":"hardware"}`,
			employeeID: employeeID,
			setAuth:    true,
			repo:       &fakeRequestRepo{},
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
		{
			name:       "repository returns error",
			body:       `{"asset_type":"laptop","category":"hardware"}`,
			employeeID: employeeID,
			setAuth:    true,
			repo: &fakeRequestRepo{
				CreateFunc: func(_ context.Context, _ *domain.AssetRequest) error {
					return errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "failed to create request" {
					t.Errorf("error = %q, want %q", resp["error"], "failed to create request")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewRequestHandler(tt.repo)
			c, w := newTestContext(http.MethodPost, "/requests", []byte(tt.body), nil)
			if tt.setAuth {
				setAuthContext(c, tt.employeeID, "employee")
			}

			h.CreateRequest(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
			tt.checkBody(t, w.Body.Bytes())
		})
	}
}

func TestListMyRequests(t *testing.T) {
	employeeID := uuid.New()
	requestID := uuid.New()
	requests := []domain.AssetRequest{*sampleRequest(requestID, employeeID)}

	tests := []struct {
		name       string
		employeeID uuid.UUID
		setAuth    bool
		repo       *fakeRequestRepo
		wantStatus int
		checkBody  func(t *testing.T, body []byte)
	}{
		{
			name:       "returns caller's requests",
			employeeID: employeeID,
			setAuth:    true,
			repo: &fakeRequestRepo{
				ListByEmployeeFunc: func(_ context.Context, id uuid.UUID) ([]domain.AssetRequest, error) {
					if id != employeeID {
						t.Errorf("employeeID = %v, want %v", id, employeeID)
					}
					return requests, nil
				},
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp []RequestResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if len(resp) != 1 {
					t.Fatalf("len(requests) = %d, want 1", len(resp))
				}
				if resp[0].EmployeeID != employeeID {
					t.Errorf("employee_id = %v, want %v", resp[0].EmployeeID, employeeID)
				}
				if resp[0].AssetType != "laptop" {
					t.Errorf("asset_type = %q, want %q", resp[0].AssetType, "laptop")
				}
			},
		},
		{
			name:       "missing auth returns unauthorized",
			setAuth:    false,
			repo:       &fakeRequestRepo{},
			wantStatus: http.StatusUnauthorized,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "invalid employee identity" {
					t.Errorf("error = %q, want %q", resp["error"], "invalid employee identity")
				}
			},
		},
		{
			name:       "repository error",
			employeeID: employeeID,
			setAuth:    true,
			repo: &fakeRequestRepo{
				ListByEmployeeFunc: func(_ context.Context, _ uuid.UUID) ([]domain.AssetRequest, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "failed to list requests" {
					t.Errorf("error = %q, want %q", resp["error"], "failed to list requests")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewRequestHandler(tt.repo)
			c, w := newTestContext(http.MethodGet, "/requests", nil, nil)
			if tt.setAuth {
				setAuthContext(c, tt.employeeID, "employee")
			}

			h.ListMyRequests(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
			tt.checkBody(t, w.Body.Bytes())
		})
	}
}

func TestGetRequest(t *testing.T) {
	requestID := uuid.New()
	ownerID := uuid.New()
	otherEmployeeID := uuid.New()

	tests := []struct {
		name       string
		idParam    string
		employeeID uuid.UUID
		role       string
		setAuth    bool
		repo       *fakeRequestRepo
		wantStatus int
		checkBody  func(t *testing.T, body []byte)
	}{
		{
			name:       "employee views own request",
			idParam:    requestID.String(),
			employeeID: ownerID,
			role:       "employee",
			setAuth:    true,
			repo: &fakeRequestRepo{
				GetByIDFunc: func(_ context.Context, id uuid.UUID) (*domain.AssetRequest, error) {
					if id != requestID {
						t.Errorf("id = %v, want %v", id, requestID)
					}
					return sampleRequest(requestID, ownerID), nil
				},
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp RequestResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.ID != requestID {
					t.Errorf("id = %v, want %v", resp.ID, requestID)
				}
				if resp.EmployeeID != ownerID {
					t.Errorf("employee_id = %v, want %v", resp.EmployeeID, ownerID)
				}
			},
		},
		{
			name:       "manager views another employee's request",
			idParam:    requestID.String(),
			employeeID: uuid.New(),
			role:       "manager",
			setAuth:    true,
			repo: &fakeRequestRepo{
				GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.AssetRequest, error) {
					return sampleRequest(requestID, ownerID), nil
				},
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				var resp RequestResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.ID != requestID {
					t.Errorf("id = %v, want %v", resp.ID, requestID)
				}
			},
		},
		{
			name:       "invalid UUID format",
			idParam:    "not-a-uuid",
			employeeID: ownerID,
			role:       "employee",
			setAuth:    true,
			repo:       &fakeRequestRepo{},
			wantStatus: http.StatusBadRequest,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "invalid request id" {
					t.Errorf("error = %q, want %q", resp["error"], "invalid request id")
				}
			},
		},
		{
			name:       "missing auth returns unauthorized",
			idParam:    requestID.String(),
			setAuth:    false,
			repo:       &fakeRequestRepo{},
			wantStatus: http.StatusUnauthorized,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "invalid employee identity" {
					t.Errorf("error = %q, want %q", resp["error"], "invalid employee identity")
				}
			},
		},
		{
			name:       "request not found",
			idParam:    requestID.String(),
			employeeID: ownerID,
			role:       "employee",
			setAuth:    true,
			repo: &fakeRequestRepo{
				GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.AssetRequest, error) {
					return nil, repository.ErrRequestNotFound
				},
			},
			wantStatus: http.StatusNotFound,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "request not found" {
					t.Errorf("error = %q, want %q", resp["error"], "request not found")
				}
			},
		},
		{
			name:       "other repository error",
			idParam:    requestID.String(),
			employeeID: ownerID,
			role:       "employee",
			setAuth:    true,
			repo: &fakeRequestRepo{
				GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.AssetRequest, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "failed to get request" {
					t.Errorf("error = %q, want %q", resp["error"], "failed to get request")
				}
			},
		},
		{
			name:       "employee cannot view another employee's request",
			idParam:    requestID.String(),
			employeeID: otherEmployeeID,
			role:       "employee",
			setAuth:    true,
			repo: &fakeRequestRepo{
				GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.AssetRequest, error) {
					return sampleRequest(requestID, ownerID), nil
				},
			},
			wantStatus: http.StatusNotFound,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]string
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["error"] != "request not found" {
					t.Errorf("error = %q, want %q", resp["error"], "request not found")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewRequestHandler(tt.repo)
			c, w := newTestContext(http.MethodGet, "/requests/"+tt.idParam, nil, gin.Params{
				{Key: "id", Value: tt.idParam},
			})
			if tt.setAuth {
				setAuthContext(c, tt.employeeID, tt.role)
			}

			h.GetRequest(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
			tt.checkBody(t, w.Body.Bytes())
		})
	}
}
