package users

import (
	"context"
	"database/sql"
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, user User) error
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user User) error {
	//validasi nama
	user.Name = strings.TrimSpace(user.Name)
	nameLength := utf8.RuneCountInString(user.Name)
	if nameLength < 2 || nameLength > 100 {
		return errors.New("Nama user harus terdiri dari 2-100 karakter")
	}

	//vlaidasi email
	user.Email = strings.TrimSpace(user.Email)
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return errors.New("format email tidak valid")
	}

	emailLength := utf8.RuneCountInString(user.Email)
	if emailLength < 5 || emailLength > 150 {
		return errors.New("email harus terdiri dari 5-150 karakter")
	}

	//validasi password
	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password tidak boleh kosong")
	}

	passwordLength := utf8.RuneCountInString(user.Password)
	if passwordLength < 8 {
		return errors.New("password minimal 8 karakter")
	}

	if len(user.Password) > 72 {
		return errors.New("password maksimal 72 byte")
	}

	//validasi phone
	user.Phone = strings.TrimSpace(user.Phone)

	if user.Phone == "" {
		return errors.New("nomor telepon wajib diisi")
	}

	if user.Phone != "" && len(user.Phone) < 10 {
		return errors.New("nomor telepon minimal 10 karakter")
	}

	if len(user.Phone) > 15 {
		return errors.New("nomor telepon maksimal 15 karakter")
	}

	//validasi role
	user.Role = strings.TrimSpace(user.Role)

	if user.Role == "" {
		user.Role = "customer"
	}

	switch user.Role {
	case "customer", "staff", "admin":
		// Role valid
	default:
		return errors.New("role harus customer, staff, atau admin")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Simpan user ke database
	return s.repo.CreateUser(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	// Validasi ID
	if id <= 0 {
		return errors.New("ID user tidak valid")
	}

	// Cek apakah user tersedia
	_, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user tidak ditemukan")
		}

		return err
	}

	// Hapus user
	return s.repo.DeleteUser(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user User) error {
	// Validasi ID
	if user.ID <= 0 {
		return errors.New("ID user tidak valid")
	}

	// Cek user
	_, err := s.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user tidak ditemukan")
		}
		return err
	}

	// Validasi nama
	user.Name = strings.TrimSpace(user.Name)

	nameLength := utf8.RuneCountInString(user.Name)
	if nameLength < 2 || nameLength > 100 {
		return errors.New("nama user harus terdiri dari 2-100 karakter")
	}

	// Validasi email
	user.Email = strings.TrimSpace(user.Email)

	address, err := mail.ParseAddress(user.Email)
	if err != nil || address.Address != user.Email {
		return errors.New("format email tidak valid")
	}

	emailLength := utf8.RuneCountInString(user.Email)
	if emailLength < 5 || emailLength > 150 {
		return errors.New("email harus terdiri dari 5-150 karakter")
	}

	// Cek email sudah digunakan
	exists, err := s.repo.IsEmailUsed(ctx, user.Email, user.ID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("email sudah digunakan")
	}

	// Validasi phone
	user.Phone = strings.TrimSpace(user.Phone)

	if user.Phone == "" {
		return errors.New("nomor telepon wajib diisi")
	}

	if len(user.Phone) < 10 || len(user.Phone) > 15 {
		return errors.New("nomor telepon harus terdiri dari 10-15 karakter")
	}

	// Update user
	return s.repo.UpdateUser(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*User, error) {
	// Validasi ID
	if id <= 0 {
		return nil, errors.New("ID tidak valid")
	}

	// Ambil user dari repository
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user tidak ditemukan")
		}

		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// Validasi email
	email = strings.TrimSpace(email)

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return nil, errors.New("format email tidak valid")
	}

	emailLength := utf8.RuneCountInString(email)
	if emailLength < 5 || emailLength > 150 {
		return nil, errors.New("email harus terdiri dari 5-150 karakter")
	}

	// Ambil user berdasarkan email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user tidak ditemukan")
		}

		return nil, err
	}

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]User, error) {
	// Ambil semua user dari repository
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
