package payment

type Service interface {
	CreatePayment(createPaymentSpec CreatePaymentSpec) (string, error)
	Notification(orderId string) error
}

type Repository interface {
}
