package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, request RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, request LoginRequest) (*AuthResponse, error)
	GetProfile(ctx context.Context, userID int64) (*UserResponse, error)
}

type authService struct {
	userService users.UserService
	tokenConfig TokenConfig
}

func NewAuthService(
	userService users.UserService,
	tokenConfig TokenConfig,
) AuthService {
	return &authService{
		userService: userService,
		tokenConfig: tokenConfig,
	}
}

func (s *authService) Register(
	ctx context.Context,
	request RegisterRequest,
) (*AuthResponse, error) {
	if request.Name == "" || request.Email == "" ||
		request.Password == "" || request.Phone == "" {
		return nil, errors.New("semua field wajib diisi")
	}

	// Registrasi publik selalu membuat akun customer.
	user := users.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Phone:    request.Phone,
		Role:     "customer",
	}

	if err := s.userService.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Ambil user yang baru dibuat untuk mendapatkan ID.
	createdUser, err := s.userService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	return s.createAuthResponse(createdUser)
}

func (s *authService) Login(
	ctx context.Context,
	request LoginRequest,
) (*AuthResponse, error) {
	if request.Email == "" || request.Password == "" {
		return nil, errors.New("email dan password wajib diisi")
	}

	user, err := s.userService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	return s.createAuthResponse(user)
}

func (s *authService) GetProfile(
	ctx context.Context,
	userID int64,
) (*UserResponse, error) {
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  user.Role,
	}, nil
}

func (s *authService) createAuthResponse(
	user *users.User,
) (*AuthResponse, error) {
	if s.tokenConfig.Secret == "" {
		return nil, errors.New("JWT secret belum dikonfigurasi")
	}

	if s.tokenConfig.Expiration <= 0 {
		return nil, errors.New("durasi JWT tidak valid")
	}

	now := time.Now()

	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "driverent-api",
			Subject:   "user",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenConfig.Expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.tokenConfig.Secret))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: tokenString,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Role:  user.Role,
		},
	}, nil
}
