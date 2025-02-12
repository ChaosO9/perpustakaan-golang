package forms

import (
	"encoding/json"
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
)

type UserForm struct{}

// UpdateMemberForm represents the form data for updating a member.
type UpdateMemberForm struct {
	Nama          string    `form:"nama" json:"nama"`
	Alamat        string    `form:"alamat" json:"alamat"`
	NomorTelepon  string    `form:"nomor_telepon" json:"nomor_telepon"`
	Email         string    `form:"email" json:"email" binding:"email"`
	TanggalLahir  time.Time `form:"tanggal_lahir" json:"tanggal_lahir"`
	StatusAnggota int16     `form:"status_anggota" json:"status_anggota"`
	Foto          string    `form:"foto" json:"foto"`
	FotoFileName  string 	`json:"filename"`
}

type RegisterForm struct {
	Email         string    `form:"email" json:"email" binding:"required,email"`
	Password      string    `form:"password" json:"password" binding:"required,min=8"` // Add password validation
	Nama          string    `form:"nama" json:"nama" binding:"required"`
	Alamat        string    `form:"alamat" json:"alamat" binding:"required"`
	NomorTelepon  string    `form:"nomor_telepon" json:"nomor_telepon" binding:"required"`
	TanggalLahir  time.Time `form:"tanggal_lahir" json:"tanggal_lahir" binding:"required"`
	Foto		  *multipart.FileHeader 	`form:"foto" binding:"omitempty"`
	FotoFileName  string 	`json:"filename"`
}

// RegisterFormSwagger Swagger-friendly version (without file upload)
type RegisterFormSwagger struct {
	Email         string    `form:"email" json:"email" binding:"required,email"`
	Password      string    `form:"password" json:"password" binding:"required,min=8"`
	Nama          string    `form:"nama" json:"nama" binding:"required"`
	Alamat        string    `form:"alamat" json:"alamat" binding:"required"`
	NomorTelepon  string    `form:"nomor_telepon" json:"nomor_telepon" binding:"required"`
	TanggalLahir  time.Time `form:"tanggal_lahir" json:"tanggal_lahir" binding:"required"`
	StatusAnggota int16     `form:"status_anggota" json:"status_anggota" binding:"required"`
	// Foto field omitted
	FotoFileName string `json:"filename"`
}


// LoginForm represents the form data for member login.
type LoginForm struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (f UserForm) Name(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter your name"
		}
		return errMsg[0]
	case "min", "max":
		return "Your name should be between 3 to 20 characters"
	case "fullName":
		return "Name should not include any special characters or numbers"
	default:
		return "Something went wrong, please try again later"
	}
}

//Email ...
func (f UserForm) Email(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter your email"
		}
		return errMsg[0]
	case "min", "max", "email":
		return "Please enter a valid email"
	default:
		return "Something went wrong, please try again later"
	}
}

//Password ...
func (f UserForm) Password(tag string) (message string) {
	switch tag {
	case "required":
		return "Please enter your password"
	case "min", "max":
		return "Your password should be between 3 and 50 characters"
	case "eqfield":
		return "Your passwords does not match"
	default:
		return "Something went wrong, please try again later"
	}
}

//Signin ...
func (f UserForm) Login(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Email" {
				return f.Email(err.Tag())
			}
			if err.Field() == "Password" {
				return f.Password(err.Tag())
			}
		}

	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}

//Register ...
func (f UserForm) Register(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				return f.Name(err.Tag())
			}

			if err.Field() == "Email" {
				return f.Email(err.Tag())
			}

			if err.Field() == "Password" {
				return f.Password(err.Tag())
			}

		}
	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}
