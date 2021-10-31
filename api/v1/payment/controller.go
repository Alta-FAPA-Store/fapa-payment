package payment

import (
	"encoding/json"
	"go-hexagonal/api/common"
	"go-hexagonal/api/v1/payment/request"
	"go-hexagonal/business/payment"

	echo "github.com/labstack/echo/v4"
)

//Controller Get item API controller
type Controller struct {
	service payment.Service
}

//NewController Construct item API controller
func NewController(service payment.Service) *Controller {
	return &Controller{
		service,
	}
}

func (controller *Controller) CreatePayment(c echo.Context) error {
	createPaymentRequest := new(request.CreatePaymentRequest)

	err := json.NewDecoder(c.Request().Body).Decode(&createPaymentRequest)

	if err != nil {
		return c.JSON(common.NewBadRequestResponse())
	}

	urlSnap, err := controller.service.CreatePayment(*createPaymentRequest.ToCreatePaymentSpec())

	if err != nil {
		return c.JSON(common.NewErrorBusinessResponse(err))
	}

	return c.JSON(common.NewSuccessResponse(urlSnap))
}

func (controller *Controller) Notification(c echo.Context) error {
	var notificationPayload map[string]interface{}

	err := json.NewDecoder(c.Request().Body).Decode(&notificationPayload)

	if err != nil {
		return c.JSON(common.NewBadRequestResponse())
	}

	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		return c.JSON(common.NewBadRequestResponse())
	}

	err = controller.service.Notification(orderId)

	if err != nil {
		return c.JSON(common.NewErrorBusinessResponse(err))
	}

	return c.JSON(common.NewSuccessResponse(orderId))
}
