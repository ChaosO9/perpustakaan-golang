package models

import (
	"database/sql"
	"fmt"
	"log"
	"perpustakaan-golang/db"
	"perpustakaan-golang/forms"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Book represents a book in the database.
type Book struct {
	ID                            uuid.UUID      `db:"id, primarykey" json:"id"`
	Judul                         string         `db:"judul" json:"judul"`
	Pengarang                     string         `db:"pengarang" json:"pengarang"`
	Penerbit                      string         `db:"penerbit" json:"penerbit"`
	ISBN                          string         `db:"isbn" json:"isbn"`
	TahunTerbit                   int64          `db:"tahun_terbit" json:"tahun_terbit"`
	Kategori                      string         `db:"kategori" json:"kategori"`
	Deskripsi                     string         `db:"deskripsi" json:"deskripsi"`
	Foto                          sql.NullString `db:"foto" json:"foto"`
	CreatedAt                     time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt                     time.Time      `db:"updated_at" json:"updated_at"`
	JumlahEksemplar               int64          `db:"jumlah_eksemplar" json:"jumlah_eksemplar"`
	JumlahKetersediaanEksemplar int64          `db:"jumlah_ketersediaan_eksemplar" json:"jumlah_ketersediaan_eksemplar"`
}

type BookSwagger struct {
    ID                            uuid.UUID `json:"id"`
    Judul                         string    `json:"judul"`
    Pengarang                     string    `json:"pengarang"`
    Penerbit                      string    `json:"penerbit"`
    ISBN                          string    `json:"isbn"`
    TahunTerbit                   int64     `json:"tahun_terbit"`
    Kategori                      string    `json:"kategori"`
    Deskripsi                     string    `json:"deskripsi"`
    Foto                          string    `json:"foto"`
    CreatedAt                     time.Time `json:"created_at"`
    UpdatedAt                     time.Time `json:"updated_at"`
    JumlahEksemplar               int64     `json:"jumlah_eksemplar"`
    JumlahKetersediaanEksemplar int64     `json:"jumlah_ketersediaan_eksemplar"`
}

type BookCreatedResponse struct {
    Message string    `json:"message"`
    ID      uuid.UUID `json:"id"` 
}

type AllBooksResponse struct {
	Books []BookSwagger `json:"books"` 
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

type BookResponse struct { 
	Book BookSwagger `json:"book"` 
}

type BookUpdatedResponse struct {
    Message string `json:"message"`
}


type BookDeletedResponse struct {
    Message string `json:"message"`
}

// BookModel handles database operations for books.
type BookModel struct{}

func (b Book) ToBookSwagger() BookSwagger {  
    return BookSwagger{
        ID:                            b.ID,
        Judul:                         b.Judul,
        Pengarang:                     b.Pengarang,
        Penerbit:                      b.Penerbit,
        ISBN:                          b.ISBN,
        TahunTerbit:                   b.TahunTerbit,
        Kategori:                      b.Kategori,
        Deskripsi:                     b.Deskripsi,
        Foto:                          b.Foto.String, 
        CreatedAt:                     b.CreatedAt,
        UpdatedAt:                     b.UpdatedAt,
        JumlahEksemplar:               b.JumlahEksemplar,
        JumlahKetersediaanEksemplar: b.JumlahKetersediaanEksemplar,
    }
}

func ToBooksSwagger(books []Book) []BookSwagger {
    booksSwagger := make([]BookSwagger, len(books))
    for i, b := range books {
        booksSwagger[i] = b.ToBookSwagger()
    }
    return booksSwagger
}

func (m BookModel) UpdateAvailability(id uuid.UUID, update forms.UpdateAvailability) error {
    var judul, pengarang, penerbit, isbn, kategori, deskripsi, foto sql.NullString
    var tahunTerbit, jumlahEksemplar, jumlahKetersediaanEksemplar sql.NullInt64
    if update.Judul != "" {
        judul = sql.NullString{String: update.Judul, Valid: true}
    }
    if update.Pengarang != "" {
        pengarang = sql.NullString{String: update.Pengarang, Valid: true}
    }
    if update.Penerbit != "" {
        penerbit = sql.NullString{String: update.Penerbit, Valid: true}
    }
    if update.Isbn != "" {
        isbn = sql.NullString{String: update.Isbn, Valid: true}
    }
    if update.Kategori != "" {
        kategori = sql.NullString{String: update.Kategori, Valid: true}
    }
    if update.Deskripsi != "" {
        deskripsi = sql.NullString{String: update.Deskripsi, Valid: true}
    }
    if update.Foto != "" {
        foto = sql.NullString{String: update.Foto, Valid: true}
    }
    if update.TahunTerbit != 0 {
        tahunTerbit = sql.NullInt64{Int64: update.TahunTerbit, Valid: true}
    }
    if update.JumlahEksemplar != 0 {
        jumlahEksemplar = sql.NullInt64{Int64: update.JumlahEksemplar, Valid: true}
    }
    if update.JumlahKetersediaanEksemplar != 0 {
        jumlahKetersediaanEksemplar = sql.NullInt64{Int64: update.JumlahKetersediaanEksemplar, Valid: true}
    }

    query, args, err := sqlx.In(`UPDATE public.master_buku SET judul = ?, pengarang = ?, penerbit = ?, isbn = ?, tahun_terbit = ?, kategori = ?, deskripsi = ?, foto = ?, jumlah_eksemplar = ?, jumlah_ketersediaan_eksemplar = ?, updated_at = ? WHERE id = ?`,
        judul, pengarang, penerbit, isbn, tahunTerbit, kategori, deskripsi, foto, jumlahEksemplar, jumlahKetersediaanEksemplar, time.Now(), id)
    if err != nil {
        return fmt.Errorf("failed to build query for BookModel.Update: %w", err)
    }

    _, err = db.GetDB().Exec(query, args...)
    if err != nil {
        log.Printf("Error updating book: %v", err) // Add logging for better debugging
        return fmt.Errorf("failed to update book: %w", err)
    }

    return nil
}

// Create adds a new book to the database.
func (m BookModel) Create(form forms.CreateBookForm) (bookID uuid.UUID, err error) {
	bookID = uuid.New()
	now := time.Now()
	currentDate := now.Format("2006-01-02")

    var foto sql.NullString
    if form.FotoFileName != "" {
        foto.String = form.FotoFileName
        foto.Valid = true
    }

	_, err = db.GetDB().Exec("INSERT INTO public.master_buku(id, judul, pengarang, penerbit, isbn, tahun_terbit, kategori, deskripsi, foto, jumlah_eksemplar, jumlah_ketersediaan_eksemplar, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		bookID, form.Judul, form.Pengarang, form.Penerbit, form.ISBN, form.TahunTerbit, form.Kategori, form.Deskripsi, foto, form.JumlahEksemplar, form.JumlahEksemplar, currentDate, currentDate)
	return bookID, err
}


func (m BookModel) All(page, limit int) ([]Book, int, error) {
    if page <= 0 {
        page = 1 
    }
    offset := (page - 1) * limit

    books := []Book{}
    query := `SELECT id, judul, pengarang, penerbit, isbn, tahun_terbit, kategori, deskripsi, foto, created_at, updated_at, jumlah_eksemplar, jumlah_ketersediaan_eksemplar FROM public.master_buku LIMIT $1 OFFSET $2`
    rows, err := db.GetDB().Query(query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get all books: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var book Book
        err := rows.Scan(&book.ID, &book.Judul, &book.Pengarang, &book.Penerbit, &book.ISBN, &book.TahunTerbit, &book.Kategori, &book.Deskripsi, &book.Foto, &book.CreatedAt, &book.UpdatedAt, &book.JumlahEksemplar, &book.JumlahKetersediaanEksemplar)
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan book row: %w", err)
        }
        books = append(books, book)
    }

    if err := rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("error iterating rows: %w", err)
    }

    var total int
    err = db.GetDB().Get(&total, "SELECT COUNT(*) FROM public.master_buku") 
    if err != nil {
        return nil, 0, err
    }

    return books, total, nil
}




// One retrieves a single book by ID.
func (m BookModel) One(id uuid.UUID) (Book, error) {
	var book Book
	err := db.GetDB().Get(&book, "SELECT * FROM public.master_buku WHERE id = $1", id)
	if err != nil {
		return Book{}, fmt.Errorf("failed to get book: %w", err)
	}
	return book, nil
}

// Update modifies an existing book.
func (m BookModel) Update(id uuid.UUID, form forms.CreateBookForm) error {
    var foto sql.NullString
    if form.FotoFileName != "" {
        foto.String = form.FotoFileName
    }
    update := forms.UpdateAvailability{
        Judul:                       form.Judul,
        Pengarang:                   form.Pengarang,
        Penerbit:                    form.Penerbit,
        Isbn:                        form.ISBN,
        TahunTerbit:                 form.TahunTerbit,
        Kategori:                    form.Kategori,
        Deskripsi:                   form.Deskripsi,
        Foto:                        foto.String,
        JumlahEksemplar:             form.JumlahEksemplar,
        JumlahKetersediaanEksemplar: form.JumlahEksemplar,
    }
    return m.UpdateAvailability(id, update)
}

// Delete removes a book from the database.
func (m BookModel) Delete(id uuid.UUID) error {
	_, err := db.GetDB().Exec("DELETE FROM public.master_buku WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}
	return nil
}

