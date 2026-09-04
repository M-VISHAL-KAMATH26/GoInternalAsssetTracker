package mq

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ApprovalDecidedEvent is the message published whenever a manager
// approves or rejects an asset request. Both the publisher (Request
// Service) and any future consumer (Phase 8) import this same struct,
// so the event schema can't silently drift between them.
type ApprovalDecidedEvent struct {
	RequestID  uuid.UUID `json:"request_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
	Decision   string    `json:"decision"` // "approved" or "rejected"
	AssetID    string    `json:"asset_id,omitempty"`
	Comment    string    `json:"comment,omitempty"`
	DecidedAt  time.Time `json:"decided_at"`
}

func (e ApprovalDecidedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}