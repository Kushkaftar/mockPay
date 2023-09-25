package postback

import (
	"log"
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/models"
	"mockPay/pkg/client"
)

type Postback struct {
	client     *client.Client
	repository postgres_db.Postback
}

func NewPostback(repository postgres_db.Postback) *Postback {
	return &Postback{
		client:     client.NewClient(),
		repository: repository,
	}
}

func (p *Postback) SendPostback(transactoin models.Transaction) {

	postbackData := models.Postback{
		MerchantID: transactoin.MerchantID,
	}

	if err := p.repository.GetPostback(&postbackData); err != nil {
		log.Printf("didn't receive postback, error - %s", err)
		return
	}

	// get url
	url, err := queryReplase(transactoin, postbackData.PostbackUrl)
	if err != nil {
		log.Printf("error parse url, error - %s", err)
	}
	log.Printf("url - %s", *url)

	// postback enabled
	if !postbackData.IsEnabled {
		log.Println("postback disabled")
		return
	}

	// send postback
	resp, err := p.client.Send(postbackData.PostbackMethod, postbackData.PostbackUrl, nil)
	if err != nil {
		log.Printf("error sending postback, error - %s", err)
		return
	}

	log.Printf("postback response - %+v", resp)
}
