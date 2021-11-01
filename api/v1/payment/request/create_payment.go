package request

import (
	"go-hexagonal/business/payment"

	"github.com/midtrans/midtrans-go"
)

type CreatePaymentRequest struct {
	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	Email         string        `json:"email"`
	Phone         string        `json:"phone"`
	Address       string        `json:"address"`
	TransactionId int           `json:"transaction_id"`
	TotalPrice    int           `json:"total_price"`
	Items         []ItemDetails `json:"items"`
}

type ItemDetails struct {
	Name     string `json:"product_name"`
	Price    int64  `json:"price"`
	Quantity int32  `json:"quantity"`
}

func (col *CreatePaymentRequest) ToCreatePaymentSpec() *payment.CreatePaymentSpec {
	var createPaymentRequest payment.CreatePaymentSpec

	createPaymentRequest.FirstName = col.FirstName
	createPaymentRequest.LastName = col.LastName
	createPaymentRequest.Email = col.Email
	createPaymentRequest.Phone = col.Phone
	createPaymentRequest.Address = col.Address
	createPaymentRequest.TransactionId = col.TransactionId
	createPaymentRequest.TotalPrice = col.TotalPrice

	for _, value := range col.Items {
		createPaymentRequest.Items = append(createPaymentRequest.Items, midtrans.ItemDetails{Name: value.Name, Price: value.Price, Qty: value.Quantity})
	}

	return &createPaymentRequest
}
