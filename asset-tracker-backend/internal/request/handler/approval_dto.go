package handler

type ApprovalActionRequest struct {
	Comment string `json:"comment" binding:"max=1000"`
}

type ApprovalResponse struct {
	RequestID string `json:"request_id"`
	Decision  string `json:"decision"`
	AssetID   string `json:"asset_id,omitempty"`
	Message   string `json:"message"`
}