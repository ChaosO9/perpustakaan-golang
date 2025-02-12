package models

import (
	"database/sql"
	"errors"
	"fmt"
	"perpustakaan-golang/db"
	"perpustakaan-golang/forms"
	"time"

	// Assuming this is your database connection

	"golang.org/x/crypto/bcrypt"
)

// User represents a library member.
type User struct {
	ID            int64          `db:"id, primarykey" json:"id"`
	Email         string         `db:"email" json:"email"`
	Password      string         `db:"password" json:"password"` //Consider removing if not using password in the model
	Nama          string         `db:"nama" json:"nama"`
	UpdatedAt     int64          `db:"updated_at" json:"updated_at"`
	CreatedAt     int64          `db:"created_at" json:"created_at"`
	Alamat        string         `db:"alamat" json:"alamat"`
	NomorTelepon  string         `db:"nomor_telepon" json:"nomor_telepon"`
	TanggalLahir  time.Time      `db:"tanggal_lahir" json:"tanggal_lahir"`
	TanggalJoin   time.Time      `db:"tanggal_join" json:"tanggal_join"`
	StatusAnggota int16          `db:"status_anggota" json:"status_anggota"`
	Foto          sql.NullString `db:"foto" json:"foto"`
}

type UserSwagger struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Nama          string    `json:"nama"`
	UpdatedAt     int64     `json:"updated_at"`
	CreatedAt     int64     `json:"created_at"`
	Alamat        string    `json:"alamat"`
	NomorTelepon  string    `json:"nomor_telepon"`
	TanggalLahir  time.Time `json:"tanggal_lahir"`
	TanggalJoin   time.Time `json:"tanggal_join"`
	StatusAnggota int16     `json:"status_anggota"`
	Foto          string    `json:"foto"` // Use string for Swagger
}

type LoginResponse struct {
	Message string `json:"message"`
	User    UserSwagger   `json:"user"`
	Token   Token  `json:"token"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	User    UserSwagger   `json:"user"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

// UserModel handles database operations for members.
type UserModel struct{}
var authModel = new(AuthModel)

// All retrieves all members with pagination.  (Adapt from BookModel.All)
func (m UserModel) All(page, limit int) ([]User, int, error) {
	offset := (page - 1) * limit
	members := []User{}

	err := db.GetDB().Select(&members, `SELECT * FROM public."user" LIMIT ? OFFSET ?`, limit, offset) // Corrected SQL
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all users: %w", err)
	}

	var total int
	err = db.GetDB().Get(&total, "SELECT COUNT(*) FROM public.\"user\"") // Corrected SQL
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return members, total, nil
}


// One retrieves a single member by ID.
func (m UserModel) One(id int64) (User, error) { // Changed to int64
	var member User
	err := db.GetDB().Get(&member, "SELECT * FROM public.\"user\" WHERE id = ?", id) // Corrected SQL
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return member, nil
}

// Update modifies an existing member. (Adapt from BookModel.Update)
func (m UserModel) Update(id int64, form forms.UpdateMemberForm) error {
    var foto sql.NullString
    if form.FotoFileName != "" {
        foto = sql.NullString{String: form.FotoFileName, Valid: true}
    }

	_, err := db.GetDB().Exec(`UPDATE public."user" SET nama = ?, alamat = ?, nomor_telepon = ?, email = ?, tanggal_lahir = ?, status_anggota = ?, foto = ?, updated_at = ? WHERE id = ?`, //Corrected SQL
		form.Nama, form.Alamat, form.NomorTelepon, form.Email, form.TanggalLahir, form.StatusAnggota, foto, time.Now().Unix(), id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete removes a member from the database.
func (m UserModel) Delete(id int64) error { // Changed to int64
	_, err := db.GetDB().Exec("DELETE FROM public.\"user\" WHERE id = ?", id) // Corrected SQL
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (m UserModel) Login(form forms.LoginForm) (User, Token, error) {
	var user User
	err := db.GetDB().Get(&user, `SELECT * FROM public."user" WHERE LOWER(email) = ? LIMIT 1`, form.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, Token{}, errors.New("invalid email or password")
		}
		return User{}, Token{}, fmt.Errorf("database error during login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		return User{}, Token{}, errors.New("invalid email or password")
	}

	tokenDetails, err := authModel.CreateToken(user.ID)
	if err != nil {
		return User{}, Token{}, fmt.Errorf("failed to create token: %w", err)
	}

	//Save token details
	saveErr := authModel.CreateAuth(user.ID, tokenDetails)
	if saveErr != nil {
		return user, Token{}, fmt.Errorf("failed to save token: %w", saveErr)
	}

	token := Token{AccessToken: tokenDetails.AccessToken, RefreshToken: tokenDetails.RefreshToken}
	return user, token, nil
}

func (m UserModel) Register(form forms.RegisterForm) (User, error) {
	// Check if email already exists
	var count int
	err := db.GetDB().Get(&count, `SELECT COUNT(*) FROM public."user" WHERE LOWER(email) = $1`, form.Email)
	if err != nil {
		return User{}, fmt.Errorf("database error checking email: %w", err)
	}
	if count > 0 {
		return User{}, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("password hashing failed: %w", err)
	}

    var foto sql.NullString
    if form.FotoFileName != "" {
        foto.String = form.FotoFileName
        foto.Valid = true
    }

	now := time.Now()
	user := User{
		Email:         form.Email,
		Password:      string(hashedPassword),
		Nama:          form.Nama,
		UpdatedAt:     now.Unix(),
		CreatedAt:     now.Unix(),
		Alamat:        form.Alamat,
		NomorTelepon:  form.NomorTelepon,
		TanggalLahir:  form.TanggalLahir,
		TanggalJoin:   now,
		StatusAnggota: 1,
		Foto:          foto,
	}


	//Note that this uses RETURNING id to get the newly created ID
	err = db.GetDB().QueryRowx(`INSERT INTO public."user"(email, password, nama, updated_at, created_at, alamat, nomor_telepon, tanggal_lahir, tanggal_join, status_anggota, foto) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		user.Email, user.Password, user.Nama, user.UpdatedAt, user.CreatedAt, user.Alamat, user.NomorTelepon, user.TanggalLahir, user.TanggalJoin, user.StatusAnggota, user.Foto).Scan(&user.ID)
	if err != nil {
		return User{}, fmt.Errorf("database error registering user: %w", err)
	}

	return user, nil
}
