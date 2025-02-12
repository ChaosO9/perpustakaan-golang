package models

import (
	"database/sql"
	"fmt"
	"time"

	"perpustakaan-golang/db"

	uuid "github.com/google/uuid"
)

// BorrowedBookReport represents a single entry in the borrowed book report.
type BorrowedBookReport struct {
	IDBuku           uuid.UUID      `db:"id_buku" json:"id_buku"`
	Judul            string         `db:"judul" json:"judul"`
	Pengarang        string         `db:"pengarang" json:"pengarang"`
	IDAnggota        int64          `db:"id_anggota" json:"id_anggota"` // Changed to int64 to match database
	Nama             string         `db:"nama" json:"nama"`
	TanggalPinjam    time.Time      `db:"tanggal_pinjam" json:"tanggal_pinjam"`
	TanggalJatuhTempo time.Time      `db:"tanggal_jatuh_tempo" json:"tanggal_jatuh_tempo"`
	TanggalKembali   sql.NullTime   `db:"tanggal_kembali" json:"tanggal_kembali"`
}

type ReturnedBookReport struct {
	IDBuku           uuid.UUID      `db:"id_buku" json:"id_buku"`
	Judul            string         `db:"judul" json:"judul"`
	Pengarang        string         `db:"pengarang" json:"pengarang"`
	IDAnggota        int64          `db:"id_anggota" json:"id_anggota"` // Changed to int64
	Nama             string         `db:"nama" json:"nama"`
	TanggalPinjam    time.Time      `db:"tanggal_pinjam" json:"tanggal_pinjam"`
	TanggalJatuhTempo time.Time      `db:"tanggal_jatuh_tempo" json:"tanggal_jatuh_tempo"`
	TanggalKembali   time.Time      `db:"tanggal_kembali" json:"tanggal_kembali"`
	Denda            sql.NullInt64  `db:"denda" json:"denda"`
}

type BorrowedBooksReportResponse struct {
    Report []BorrowedBookReportSwagger `json:"report"`
}


type ReturnedBooksReportResponse struct {
    Report []ReturnedBookReportSwagger `json:"report"`
}

type BorrowedBookReportSwagger struct {
    IDBuku           uuid.UUID `json:"id_buku"`
    Judul            string    `json:"judul"`
    Pengarang        string    `json:"pengarang"`
    IDAnggota        int64 `json:"id_anggota"`
    Nama             string    `json:"nama"`
    TanggalPinjam    time.Time `json:"tanggal_pinjam"`
    TanggalJatuhTempo time.Time `json:"tanggal_jatuh_tempo"`
    TanggalKembali   string `json:"tanggal_kembali"`  
}

type ReturnedBookReportSwagger struct {
	IDBuku           uuid.UUID `json:"id_buku"`
	Judul            string    `json:"judul"`
	Pengarang        string    `json:"pengarang"`
	IDAnggota        int64 `json:"id_anggota"`
	Nama             string    `json:"nama"`
	TanggalPinjam    time.Time `json:"tanggal_pinjam"`
	TanggalJatuhTempo time.Time `json:"tanggal_jatuh_tempo"`
	TanggalKembali   time.Time `json:"tanggal_kembali"`
	Denda            int64     `json:"denda"` // Use int64 for Swagger
}

// GetBorrowedBookReports retrieves borrowed book report data.
func GetBorrowedBookReports(startDate, endDate time.Time) ([]BorrowedBookReport, error) {
	query := `
		SELECT mb.id AS id_buku, mb.judul, mb.pengarang, u.id AS id_anggota, u.nama, tb.tanggal_pinjam, tb.tanggal_jatuh_tempo, tb.tanggal_kembali
		FROM public.transaksi_buku tb
		JOIN public.master_buku mb ON tb.id_buku = mb.id
		JOIN public."user" u ON tb.id_anggota = u.id
		WHERE tb.tanggal_pinjam BETWEEN $1 AND $2` 

	var reports []BorrowedBookReport
	err := db.GetDB().Select(&reports, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get borrowed book reports: %w", err)
	}
	return reports, nil
}


func GetReturnedBookReports(startDate, endDate time.Time) ([]ReturnedBookReport, error) {
	query := `
		SELECT mb.id AS id_buku, mb.judul, mb.pengarang, u.id AS id_anggota, u.nama, tb.tanggal_pinjam, tb.tanggal_jatuh_tempo, tb.tanggal_kembali, tb.denda
		FROM public.transaksi_buku tb
		JOIN public.master_buku mb ON tb.id_buku = mb.id
		JOIN public."user" u ON tb.id_anggota = u.id
		WHERE tb.tanggal_kembali BETWEEN $1 AND $2 AND tb.status = 1` 

	var reports []ReturnedBookReport
	err := db.GetDB().Select(&reports, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get returned book reports: %w", err)
	}
	return reports, nil
}