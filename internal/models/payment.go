package models

import "time"

type Payment struct {
	Id            int
	PaymentDate   time.Time
	Amount        float32
	EmailShop     string
	EmailCustomer string
}
