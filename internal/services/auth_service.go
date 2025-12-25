package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/khusa-mahal/backend/internal/models"
	"github.com/khusa-mahal/backend/internal/repository/mongodb"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     *mongodb.UserRepository
	otpRepo      *mongodb.OTPRepository
	emailService *EmailService
	jwtSecret    []byte
}

func NewAuthService(userRepo *mongodb.UserRepository, otpRepo *mongodb.OTPRepository, emailService *EmailService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		emailService: emailService,
		jwtSecret:    []byte(os.Getenv("JWT_SECRET")),
	}
}

// Register creates a temp user (unverified) and sends OTP
func (s *AuthService) Register(ctx context.Context, input models.User, password string) error {
	// 1. Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, input.Email)
	if existingUser != nil {
		if existingUser.IsVerified {
			return errors.New("user already exists with this email")
		}
		// If exists but not verified, we can resend OTP or overwrite.
		// For simplicity, we'll overwrite/ignore and just send new OTP logic purely.
		// Actually, standard is to not allow duplicate keys. Let's assume we proceed.
	}

	// 2. Hash Password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	input.PasswordHash = string(hashedBytes)
	input.IsVerified = false

	// 3. Save User (if new)
	if existingUser == nil {
		if err := s.userRepo.Create(ctx, &input); err != nil {
			return err
		}
	}

	// 4. Generate OTP
	otpCode, err := s.generateOTP(6)
	if err != nil {
		return err
	}

	// 5. Save OTP
	otp := &models.OTP{
		Email:     input.Email,
		Code:      otpCode,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := s.otpRepo.Save(ctx, otp); err != nil {
		return err
	}


	// 6. Send Email (async - don't block signup)
	go func() {
		if err := s.emailService.SendOTP(input.Email, otpCode); err != nil {
			// Log error but don't fail registration
			fmt.Printf("Warning: Failed to send OTP email to %s: %v\n", input.Email, err)
		}
	}()

	return nil
}

// VerifyOTP checks code and activates user
func (s *AuthService) VerifyOTP(ctx context.Context, email, code string) (string, *models.User, error) {
	// 1. Find Valid OTP
	_, err := s.otpRepo.FindValidOTP(ctx, email, code)
	if err != nil {
		return "", nil, errors.New("invalid or expired OTP")
	}

	// 2. Mark User Verified
	if err := s.userRepo.UpdateVerification(ctx, email); err != nil {
		return "", nil, err
	}

	// 3. Cleanup OTPs
	_ = s.otpRepo.DeleteByEmail(ctx, email)

	// 4. Generate Token (Auto Login) - get user first to get ID
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}

	token, err := s.GenerateToken(user)
	return token, user, err
}

// Login validates user and returns token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return "", nil, errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := s.GenerateToken(user)
	return token, user, err
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"userId": user.ID.Hex(),
		"email":  user.Email,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) generateOTP(length int) (string, error) {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}
