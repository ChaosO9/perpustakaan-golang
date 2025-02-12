package forms

import (
	"encoding/json"
	"mime/multipart"

	"github.com/go-playground/validator/v10"
)

// CreateBookForm represents the form data for creating a book.
type CreateBookForm struct {
	Judul                         string    `form:"judul" json:"judul" binding:"required,min=3,max=255"`
	Pengarang                     string    `form:"pengarang" json:"pengarang" binding:"required,min=3,max=255"`
	Penerbit                      string    `form:"penerbit" json:"penerbit" binding:"required,min=3,max=255"`
	ISBN                          string    `form:"isbn" json:"isbn" binding:"required,min=10,max=17"`
	TahunTerbit                   int64     `form:"tahun_terbit" json:"tahun_terbit" binding:"required,numeric"`
	Kategori                      string    `form:"kategori" json:"kategori" binding:"required,min=3,max=255"`
	Deskripsi                     string    `form:"deskripsi" json:"deskripsi" binding:"required,min=3,max=1000"`
	Foto			 *multipart.FileHeader 	`form:"foto" binding:"omitempty"`
	FotoFileName  					string 		`json:"filename"`
	JumlahEksemplar               int64     `form:"jumlah_eksemplar" json:"jumlah_eksemplar" binding:"required,numeric,min=1"`
}

type UpdateAvailability struct {
    Judul                       string         `json:"judul"`
    Pengarang                   string         `json:"pengarang"`
    Penerbit                    string         `json:"penerbit"`
    Isbn                        string         `json:"isbn"`
    TahunTerbit                 int64          `json:"tahun_terbit"`
    Kategori                    string         `json:"kategori"`
    Deskripsi                   string         `json:"deskripsi"`
    Foto                        string         `json:"foto"`
    JumlahEksemplar             int64          `json:"jumlah_eksemplar"`
    JumlahKetersediaanEksemplar int64          `json:"jumlah_ketersediaan_eksemplar"`
}

// BookForm handles validation and error messages.  (You can adapt the ArticleForm functions here)
type BookForm struct{}

//Create handles form validation error.
func (f BookForm) Create(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Title":
				return "Title should be between 3 to 255 characters"
			case "Author":
				return "Author should be between 3 to 255 characters"
			case "Publisher":
				return "Publisher should be between 3 to 255 characters"
			case "ISBN":
				return "Invalid ISBN"
			case "YearPublished":
				return "Invalid year published"
			case "Category":
				return "Category should be between 3 to 255 characters"
			case "Description":
				return "Description should be between 3 to 1000 characters"
			case "NumberOfCopies":
				return "Number of copies must be a positive number"
			case "NumberOfAvailableCopies":
				return "Number of available copies must be a positive number"
			default:
				return "Invalid Request"
			}
		}
	default:
		return "Invalid Request"
	}
	return "Something went wrong, please try again later"
}


//Update handles form validation error on update.
func (f BookForm) Update(err error) string {
	return f.Create(err) // Reuse Create for simplicity
}
