package payment

import (
	"errors"
	"go-hexagonal/business"
	"go-hexagonal/util/validator"

	"github.com/midtrans/midtrans-go"
)

type service struct {
	repository Repository
}

//NewService Construct pet service object
func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

type CreatePaymentSpec struct {
	FirstName     string                 `validate:"required"`
	LastName      string                 `validate:"required"`
	Email         string                 `validate:"required"`
	Phone         string                 `validate:"required"`
	Address       string                 `validate:"required"`
	TransactionId int                    `validate:"required"`
	TotalPrice    int                    `validate:"required"`
	Items         []midtrans.ItemDetails `validate:"required"`
}

func (s *service) CreatePayment(createPaymentSpec CreatePaymentSpec) (string, error) {
	err := validator.GetValidator().Struct(createPaymentSpec)
	if err != nil {
		return "", business.ErrInvalidSpec
	}

	InitializeSnapClient()
	paymentUrl := createTransaction(createPaymentSpec)

	return paymentUrl, nil
}

func (s *service) Notification(orderId string) error {
	InitializeCoreClient()
	transactionStatusResp, err := checkTransaction(orderId)

	if err != nil {
		return errors.New(err.Error())
	} else {
		var status string
		if transactionStatusResp != nil {
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					status = "challenge"
				} else if transactionStatusResp.FraudStatus == "accept" {
					status = "success"
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				status = "success"
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				status = "failed"
			} else if transactionStatusResp.TransactionStatus == "pending" {
				status = "pending"
			}
		}

		err := updateTransaction(transactionStatusResp.OrderID, status)

		if err != nil {
			return err
		}
	}

	return nil
}
