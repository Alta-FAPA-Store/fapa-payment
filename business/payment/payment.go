package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type CustomerDetails struct {
	FirstName    string
	LastName     string
	Email        string
	Phone        string
	Address      string
	TotalPayment int
}

var snapC snap.Client
var coreC coreapi.Client

func InitializeSnapClient() {
	snapC.New("SB-Mid-server-0qpp_T4NqLWf8ifdV4kJoKhl", midtrans.Sandbox)
}

func InitializeCoreClient() {
	coreC.New("SB-Mid-server-0qpp_T4NqLWf8ifdV4kJoKhl", midtrans.Sandbox)
}

func createTransaction(createPaymentSpec CreatePaymentSpec) string {
	resp, err := snapC.CreateTransaction(GenerateSnapReq(createPaymentSpec))
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}

	return resp.RedirectURL
}

func checkTransaction(orderId string) (*coreapi.TransactionStatusResponse, error) {
	transactionStatusResp, err := coreC.CheckTransaction(orderId)

	if err != nil {
		return nil, errors.New(err.GetMessage())
	}

	return transactionStatusResp, nil
}

func GenerateSnapReq(createPaymentSpec CreatePaymentSpec) *snap.Request {
	// Initiate Customer address
	custAddress := &midtrans.CustomerAddress{
		FName:   createPaymentSpec.FirstName,
		LName:   createPaymentSpec.LastName,
		Phone:   createPaymentSpec.Phone,
		Address: createPaymentSpec.Address,
	}

	// Initiate Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "STORE-" + strconv.Itoa(createPaymentSpec.TransactionId),
			GrossAmt: int64(createPaymentSpec.TotalPrice),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    createPaymentSpec.FirstName,
			LName:    createPaymentSpec.LastName,
			Email:    createPaymentSpec.Email,
			Phone:    createPaymentSpec.Phone,
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items:           &createPaymentSpec.Items,
	}
	return snapReq
}

func updateTransaction(transactionId string, status string) error {
	postBody, _ := json.Marshal(map[string]string{
		"status": status,
	})

	id := strings.Split(transactionId, "-")
	requestBody := bytes.NewBuffer(postBody)

	req, _ := http.NewRequest("PUT", "http://127.0.0.1:8080/v1/transaction/"+id[1], requestBody)
	req.Header.Add("Content-Type", "application/json")

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
