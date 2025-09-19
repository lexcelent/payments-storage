package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lexcelent/payments-storage/internal/models"
	"github.com/lexcelent/payments-storage/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s := &Storage{db: db}

	// TODO: this is probably very bad decision to create table here
	s.createTable()

	return s, nil
}

func (s *Storage) createTable() error {

	// Сохранил на всякий случай
	// DROP TABLE IF EXISTS payments;

	query := `
		CREATE TABLE IF NOT EXISTS payments(
			id INTEGER PRIMARY KEY,
			payment_date DATETIME,
			amount FLOAT,
			email_shop TEXT,
			email_customer TEXT
		);
	`

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Payment returns payment by id
func (s *Storage) Payment(id int) (models.Payment, error) {
	const op = "storage.sqlite.Payment"

	stmt, err := s.db.Prepare("SELECT id, payment_date, amount, email_shop, email_customer FROM payments WHERE id = ?")
	if err != nil {
		return models.Payment{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRow(id)

	var payment models.Payment
	err = row.Scan(&payment.Id, &payment.PaymentDate, &payment.Amount, &payment.EmailCustomer, &payment.EmailShop)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Payment{}, storage.ErrPaymentNotExists
		}

		return models.Payment{}, fmt.Errorf("%s: %w", op, err)
	}

	return payment, nil
}

// PaymentAdd add new payment info into db
func (s *Storage) PaymentAdd(paymentDate time.Time, amount float32, emailShop string, emailCustomer string) (int64, error) {
	const op = "storage.sqlite.PaymentAdd"

	stmt, err := s.db.Prepare("INSERT INTO payments(payment_date, amount, email_shop, email_customer) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(paymentDate, amount, emailShop, emailCustomer)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
