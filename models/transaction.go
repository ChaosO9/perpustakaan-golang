package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"perpustakaan-golang/db"
	"perpustakaan-golang/utils"

	"github.com/google/uuid"
)

// Transaction represents a book transaction.
type Transaction struct {
	ID               uuid.UUID      `db:"id, primarykey" json:"id"`
	IDAnggota        int64      `db:"id_anggota" json:"id_anggota"`
	IDBuku           uuid.UUID      `db:"id_buku" json:"id_buku"`
	TanggalPinjam    time.Time      `db:"tanggal_pinjam" json:"tanggal_pinjam"`
	TanggalJatuhTempo time.Time      `db:"tanggal_jatuh_tempo" json:"tanggal_jatuh_tempo"`
	TanggalKembali   sql.NullTime   `db:"tanggal_kembali" json:"tanggal_kembali"`
	Denda            sql.NullInt64  `db:"denda" json:"denda"`
	Status           int16          `db:"status" json:"status"` // 0: borrowed, 1: returned, ...
}

type TransactionSwagger struct {
	ID               uuid.UUID `json:"id"`
	IDAnggota        int64 `json:"id_anggota"`
	IDBuku           uuid.UUID `json:"id_buku"`
	TanggalPinjam    time.Time `json:"tanggal_pinjam"`
	TanggalJatuhTempo time.Time `json:"tanggal_jatuh_tempo"`
	TanggalKembali   string    `json:"tanggal_kembali"` // Use string for Swagger
	Denda            int64     `json:"denda"`          // Use int64 (or string) for Swagger
	Status           int16     `json:"status"`
}

type GetTransactionsByMemberIDResponse struct {
	Transactions []TransactionSwagger `json:"transactions"` // Use the Swagger struct
}

type TransactionModel struct{}

//Borrow creates a new transaction.  You'll need to adjust the "loan duration" logic here.
func (tm TransactionModel) Borrow(idAnggota int64, idBuku uuid.UUID, loanDuration time.Duration) (uuid.UUID, error) {
	transactionID := uuid.New()
	now := time.Now()
	dueDate := now.Add(loanDuration)

	tx, err := db.GetDB().Beginx()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(`UPDATE public.master_buku 
                        SET jumlah_ketersediaan_eksemplar = jumlah_ketersediaan_eksemplar - 1
                        WHERE id = $1 AND jumlah_ketersediaan_eksemplar > 0 
                        RETURNING jumlah_ketersediaan_eksemplar`, idBuku)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to update book availability: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return uuid.Nil, errors.New("book is unavailable")
	}

	_, err = tx.Exec("INSERT INTO public.transaksi_buku (id, id_anggota, id_buku, tanggal_pinjam, tanggal_jatuh_tempo, status) VALUES ($1, $2, $3, $4, $5, $6)",
		transactionID, idAnggota, idBuku, now, dueDate, 0)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to borrow book: %w", err)
	}

	return transactionID, tx.Commit()
}


// Return updates a transaction to mark the book as returned.  Denda calculation needs implementation.
func (tm TransactionModel) Return(transactionID uuid.UUID, returnedAt time.Time) (int64, error) {
	tx, err := db.GetDB().Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if transaction has already been returned
	var count int
	err = tx.Get(&count, "SELECT COUNT(*) FROM public.transaksi_buku WHERE id = $1 AND status = 1", transactionID)
	if err != nil {
		return 0, fmt.Errorf("failed to check transaction status: %w", err)
	}
	if count > 0 {
		return 0, errors.New("transaction already returned")
	}

	transaction := Transaction{}
	err = tx.Get(&transaction, "SELECT * FROM public.transaksi_buku WHERE id = $1", transactionID)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction: %w", err)
	}

	denda := utils.CalculateFine(transaction.TanggalJatuhTempo, returnedAt)

	_, err = tx.Exec("UPDATE public.transaksi_buku SET tanggal_kembali = $1, denda = $2, status = $3 WHERE id = $4", returnedAt, denda, 1, transactionID)
	if err != nil {
		return 0, fmt.Errorf("failed to update transaction: %w", err)
	}

	_, err = tx.Exec(`UPDATE public.master_buku SET jumlah_ketersediaan_eksemplar = jumlah_ketersediaan_eksemplar + 1 WHERE id = $1`, transaction.IDBuku)
	if err != nil {
		return 0, fmt.Errorf("failed to update book availability: %w", err)
	}

	return denda, tx.Commit()
}


// GetTransactionByID retrieves a single transaction by ID.
func (tm TransactionModel) GetTransactionByID(id uuid.UUID) (Transaction, error) {
	transaction := Transaction{}
	err := db.GetDB().Get(&transaction, "SELECT * FROM public.transaksi_buku WHERE id = $1", id)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to get transaction by ID: %w", err)
	}
	return transaction, nil
}

// GetTransactionsByMemberID retrieves all transactions for a given member.
func (tm TransactionModel) GetTransactionsByMemberID(id uuid.UUID) ([]Transaction, error) {
	transactions := []Transaction{}
	err := db.GetDB().Select(&transactions, "SELECT * FROM public.transaksi_buku WHERE id_anggota = ?", id) //Corrected: err assignment
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by member ID: %w", err)
	}
	return transactions, nil
}
