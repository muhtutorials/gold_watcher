package repository

import (
	"database/sql"
	"errors"
	"time"
)

type SQLiteRepository struct {
	Conn *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{Conn: db}
}

func (repo *SQLiteRepository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS holdings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount REAL NOT NULL,
			purchase_date INTEGER NOT NULL,
			purchase_price INTEGER NOT NULL
		)
	`
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) InsertHolding(h Holdings) (*Holdings, error) {
	stmt := "INSERT INTO holdings (amount, purchase_date, purchase_price) VALUES (?, ?, ?)"

	result, err := repo.Conn.Exec(stmt, h.Amount, h.PurchaseDate.Unix(), h.PurchasePrice)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	h.ID = id

	return &h, nil
}

func (repo *SQLiteRepository) AllHoldings() ([]Holdings, error) {
	query := "SELECT id, amount, purchase_date, purchase_price FROM holdings ORDER BY purchase_date"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Holdings
	for rows.Next() {
		var h Holdings
		var unixTime int64
		err := rows.Scan(&h.ID, &h.Amount, &unixTime, &h.PurchasePrice)
		if err != nil {
			return nil, err
		}
		h.PurchaseDate = time.Unix(unixTime, 0)
		all = append(all, h)
	}

	return all, nil
}

func (repo *SQLiteRepository) GetHoldingByID(id int64) (*Holdings, error) {
	query := "SELECT id, amount, purchase_date, purchase_price FROM holdings WHERE id = ?"
	row := repo.Conn.QueryRow(query, id)

	var h Holdings
	var unixTime int64
	err := row.Scan(&h.ID, &h.Amount, &unixTime, &h.PurchasePrice)
	if err != nil {
		return nil, err
	}
	h.PurchaseDate = time.Unix(unixTime, 0)

	return &h, nil
}

func (repo *SQLiteRepository) UpdateHolding(id int64, h Holdings) error {
	if id == 0 {
		return errors.New("invalid updated id")
	}

	stmt := "UPDATE holdings SET amount = ?, purchase_date = ?, purchase_price = ? WHERE id = ?"
	result, err := repo.Conn.Exec(stmt, h.Amount, h.PurchaseDate.Unix(), h.PurchasePrice, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errUpdateFailed
	}

	return nil
}

func (repo *SQLiteRepository) DeleteHolding(id int64) error {
	stmt := "DELETE FROM holdings WHERE id = ?"

	result, err := repo.Conn.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errDeleteFailed
	}

	return nil
}
