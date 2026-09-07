package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// This test assumes all services (user, inventory, request,
// notification-consumer) are already running locally, same as your
// manual Postman flow — it automates that flow rather than replacing
// the services themselves. Run with: go test ./tests/integration/... -v
// after starting every service in separate terminals.

const (
	requestServiceURL   = "http://localhost:8082"
	inventoryServiceURL = "http://localhost:8081"
)

func TestFullApprovalFlow(t *testing.T) {
	employeeToken := os.Getenv("TEST_EMPLOYEE_TOKEN")
	managerToken := os.Getenv("TEST_MANAGER_TOKEN")
	adminToken := os.Getenv("TEST_ADMIN_TOKEN")
	if employeeToken == "" || managerToken == "" || adminToken == "" {
		t.Skip("TEST_EMPLOYEE_TOKEN, TEST_MANAGER_TOKEN, TEST_ADMIN_TOKEN must be set — generate via go run ./cmd/gen-test-token")
	}

	// 1. Create request as employee.
	createBody, _ := json.Marshal(map[string]string{
		"asset_type": "laptop",
		"category":   "IT Equipment",
		"justification": "integration test",
	})
	req, _ := http.NewRequest("POST", requestServiceURL+"/requests", bytes.NewReader(createBody))
	req.Header.Set("Authorization", "Bearer "+employeeToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 creating request, got %d", resp.StatusCode)
	}

	var created struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}

	// 2. Approve as manager.
	approveReq, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/requests/%s/approve", requestServiceURL, created.ID), nil)
	approveReq.Header.Set("Authorization", "Bearer "+managerToken)

	approveResp, err := http.DefaultClient.Do(approveReq)
	if err != nil {
		t.Fatalf("approve request failed: %v", err)
	}
	defer approveResp.Body.Close()
	if approveResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 approving request, got %d", approveResp.StatusCode)
	}

	var approved struct {
		AssetID string `json:"asset_id"`
	}
	if err := json.NewDecoder(approveResp.Body).Decode(&approved); err != nil {
		t.Fatalf("failed to decode approve response: %v", err)
	}
	if approved.AssetID == "" {
		t.Fatal("expected asset_id in approve response, got empty")
	}

	// 3. Confirm the asset actually flipped to "assigned" in Inventory.
	getAssetReq, _ := http.NewRequest("GET", fmt.Sprintf("%s/assets/%s", inventoryServiceURL, approved.AssetID), nil)
	getAssetReq.Header.Set("Authorization", "Bearer "+adminToken)

	assetResp, err := http.DefaultClient.Do(getAssetReq)
	if err != nil {
		t.Fatalf("get asset failed: %v", err)
	}
	defer assetResp.Body.Close()

	var asset struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(assetResp.Body).Decode(&asset); err != nil {
		t.Fatalf("failed to decode asset response: %v", err)
	}
	if asset.Status != "assigned" {
		t.Fatalf("expected asset status 'assigned', got %q", asset.Status)
	}

	// 4. Confirm a notification row landed in notification_db — give
	// the async consumer a moment to process the queued event.
	time.Sleep(2 * time.Second)

	notifDSN := os.Getenv("NOTIFICATION_DB_DSN")
	if notifDSN == "" {
		notifDSN = "root:vishal123@tcp(127.0.0.1:3306)/notification_db?charset=utf8mb4&parseTime=True&loc=Local"
	}
	db, err := sql.Open("mysql", notifDSN)
	if err != nil {
		t.Fatalf("failed to open notification_db: %v", err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notifications WHERE message LIKE ?", "%"+created.ID+"%").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query notifications: %v", err)
	}
	if count == 0 {
		t.Fatal("expected at least one notification referencing this request, found none")
	}
}