package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"asset-backend/internal/notification/domain"
	"asset-backend/internal/notification/repository"
	"asset-backend/internal/shared/mq"

	"github.com/google/uuid"
)

// procurementRecipientID is a fixed, well-known UUID representing the
// procurement team as a notification recipient — not a real employee.
var procurementRecipientID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type Notifier struct {
	repo repository.NotificationRepository
}

func NewNotifier(repo repository.NotificationRepository) *Notifier {
	return &Notifier{repo: repo}
}

// Handle decodes a raw ApprovalDecidedEvent message and produces the
// appropriate notification(s): always one for the employee, plus one
// for procurement when the decision was an approval.
func (n *Notifier) Handle(ctx context.Context, body []byte) error {
	var event mq.ApprovalDecidedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	var employeeMsg string
	var employeeType domain.NotificationType

	switch event.Decision {
	case "approved":
		employeeType = domain.NotificationTypeEmployeeApproved
		employeeMsg = fmt.Sprintf("Your equipment request %s has been approved. Asset ID: %s", event.RequestID, event.AssetID)
	case "rejected":
		employeeType = domain.NotificationTypeEmployeeRejected
		employeeMsg = fmt.Sprintf("Your equipment request %s has been rejected. Reason: %s", event.RequestID, event.Comment)
	default:
		return fmt.Errorf("unknown decision type: %s", event.Decision)
	}

	if err := n.send(ctx, event.EmployeeID, employeeType, employeeMsg); err != nil {
		return err
	}

	if event.Decision == "approved" {
		procurementMsg := fmt.Sprintf("Asset %s has been reserved for request %s — prepare for handoff.", event.AssetID, event.RequestID)
		if err := n.send(ctx, procurementRecipientID, domain.NotificationTypeProcurementHandoff, procurementMsg); err != nil {
			return err
		}
	}

	return nil
}

// send persists the notification and "delivers" it — for this phase,
// delivery is a structured log line; swapping in real email/Slack
// later only touches this one function.
func (n *Notifier) send(ctx context.Context, recipientID uuid.UUID, notifType domain.NotificationType, message string) error {
	notification := &domain.Notification{
		ID:          uuid.New(),
		RecipientID: recipientID,
		Type:        notifType,
		Message:     message,
	}

	if err := n.repo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	// Simulated delivery.
	log.Printf("[NOTIFY] to=%s type=%s message=%q", recipientID, notifType, message)

	if err := n.repo.MarkSent(ctx, notification.ID, notification.CreatedAt); err != nil {
		return fmt.Errorf("failed to mark notification sent: %w", err)
	}

	return nil
}