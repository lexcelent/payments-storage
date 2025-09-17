package add

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	PaymentDate   time.Time `json:"paymentDate"`
	Amount        float32   `json:"amount"`
	EmailShop     string    `json:"emailShop"`
	EmailCustomer string    `json:"emailCustomer"`
}

type PaymentAdder interface {
	PaymentAdd(paymentDate time.Time, amount float32, emailShop string, emailCustomer string) (int64, error)
}

func New(log *slog.Logger, paymentAdder PaymentAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.payments.add.New"

		log.Info("add payment", slog.String("op", op))

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			log.Error("method not allowed", "req_method", r.Method)
			return
		}

		// parsing request
		// TODO: validate request
		var payment Request

		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &payment)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("invalid request")
			return
		}

		id, err := paymentAdder.PaymentAdd(payment.PaymentDate, payment.Amount, payment.EmailShop, payment.EmailCustomer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("error during add payment", slog.String("err", err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("payment added", slog.Int64("id", id))
	}
}
