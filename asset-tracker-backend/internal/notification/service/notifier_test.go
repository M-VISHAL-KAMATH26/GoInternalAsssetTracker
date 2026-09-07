package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"asset-backend/internal/notification/domain"
	"asset-backend/internal/shared/mq"

	"github.com/google/uuid"
)

// fakeNotificationRepository is a hand-rolled fake, same pattern used
// across every other handler test in this project.
type fakeNotificationRepository struct {
	created []*domain.Notification
	markedSentIDs []uuid.UUID
	createErr error
	markSentErr error
}

func (f *fakeNotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	if f.createErr != nil {
		return f.createErr
	}
	n.CreatedAt = time.Now()
	f.created = append(f.created, n)
	return nil
}

func (f *fakeNotificationRepository) ListByRecipient(ctx context.Context, recipientID uuid.UUID) ([]domain.Notification, error) {
	return nil, nil
}

func (f *fakeNotificationRepository) MarkSent(ctx context.Context, id uuid.UUID, sentAt time.Time) error {
	if f.markSentErr != nil {
		return f.markSentErr
	}
	f.markedSentIDs = append(f.markedSentIDs, id)
	return nil
}

func TestNotifierHandle(t *testing.T) {
	employeeID := uuid.New()

	tests := []struct {
		name              string
		event             mq.ApprovalDecidedEvent
		wantNotifications int
		wantTypes         []domain.NotificationType
		wantErr           bool
	}{
		{
			name: "approved produces employee and procurement notifications",
			event: mq.ApprovalDecidedEvent{
				RequestID:  uuid.New(),
				EmployeeID: employeeID,
				Decision:   "approved",
				AssetID:    uuid.New().String(),
			},
			wantNotifications: 2,
			wantTypes:         []domain.NotificationType{domain.NotificationTypeEmployeeApproved, domain.NotificationTypeProcurementHandoff},
		},
		{
			name: "rejected produces only employee notification",
			event: mq.ApprovalDecidedEvent{
				RequestID:  uuid.New(),
				EmployeeID: employeeID,
				Decision:   "rejected",
				Comment:    "not eligible",
			},
			wantNotifications: 1,
			wantTypes:         []domain.NotificationType{domain.NotificationTypeEmployeeRejected},
		},
		{
			name: "unknown decision returns error",
			event: mq.ApprovalDecidedEvent{
				RequestID:  uuid.New(),
				EmployeeID: employeeID,
				Decision:   "maybe",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeNotificationRepository{}
			notifier := NewNotifier(repo)

			body, err := json.Marshal(tt.event)
			if err != nil {
				t.Fatalf("failed to marshal test event: %v", err)
			}

			err = notifier.Handle(context.Background(), body)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(repo.created) != tt.wantNotifications {
				t.Fatalf("expected %d notifications, got %d", tt.wantNotifications, len(repo.created))
			}

			for i, wantType := range tt.wantTypes {
				if repo.created[i].Type != wantType {
					t.Errorf("notification %d: expected type %s, got %s", i, wantType, repo.created[i].Type)
				}
			}

			if len(repo.markedSentIDs) != tt.wantNotifications {
				t.Errorf("expected %d notifications marked sent, got %d", tt.wantNotifications, len(repo.markedSentIDs))
			}
		})
	}
}