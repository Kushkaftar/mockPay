package postback

import (
	"mockPay/internal/pkg/models"
	"net/url"
)

const (
	status          = "{status}"
	transactionUUID = "{transaction_uuid}"
	transactionType = "{transaction_type}"
)

var transactionTypeMap = map[int]string{
	models.PurchaseType:  "purchase",
	models.ReccurentType: "recurrent",
	models.RefundType:    "refund",
}

var transactionStatusMap = map[int]string{
	models.NewStatus:        "new",
	models.ProcessingStatus: "processing",
	models.ComplitedStatus:  "complite",
	models.RejectedStatus:   "rejected",
}

func queryReplase(transactoin models.Transaction, postbackUrl string) (*string, error) {
	u, err := url.Parse(postbackUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for key, value := range q {
		for i, v := range value {
			switch v {
			case status:
				q[key][i] = transactionStatusMap[transactoin.TransactionStatus]
			case transactionType:
				q[key][i] = transactionTypeMap[transactoin.TransactionType]
			case transactionUUID:
				q[key][i] = transactoin.UUID
			}
		}
	}
	u.RawQuery = q.Encode()

	pu := u.String()
	return &pu, nil
}
