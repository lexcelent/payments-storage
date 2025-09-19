package get

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/lexcelent/payments-storage/internal/models"
	"github.com/lexcelent/payments-storage/internal/storage"
)

type Response struct {
	Id            int       `json:"id"`
	PaymentDate   time.Time `json:"payment_date"`
	Amount        float32   `json:"amount"`
	EmailShop     string    `json:"email_shop"`
	EmailCustomer string    `json:"email_customer"`
}

type PaymentGetter interface {
	Payment(id int) (models.Payment, error)
}

func New(log *slog.Logger, paymentAdder PaymentGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.payments.get.New"

		log.Info("get payment", slog.String("op", op))

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			log.Error("method not allowed", "req_method", r.Method)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("no params in query")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("cannot get param from query: id")
			return
		}

		payment, err := paymentAdder.Payment(id)
		if err != nil {
			if errors.Is(err, storage.ErrPaymentNotExists) {
				w.WriteHeader(http.StatusNotFound)
				log.Error("payment not found", slog.Int("id", id))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.Error("error during get payment", slog.String("err", err.Error()))
			return
		}

		jsonData, err := json.Marshal(payment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
