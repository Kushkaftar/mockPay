package models

// transaction type
const (
	_ = iota
	PurchaseType
	ReccurentType
	RefundType
	FormType
)

// transaction status
const (
	_ = iota
	NewStatus
	ProcessingStatus
	ComplitedStatus
	RejectedStatus
)
