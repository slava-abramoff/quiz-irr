package services

import (
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	secret string
}

func NewAuthServce(s string) *authService {
	return &authService{secret: s}
}

func (a *authService) HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	return string(hash)
}

func (a *authService) ComparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *authService) MakeAccessToken(admin *models.Admin) (string, error) {
	claims := jwt.MapClaims{
		"id":      admin.ID,
		"email":   admin.Email,
		"is_root": admin.IsRoot,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(a.secret))
}

func (a *authService) MakeRefreshToken(admin *models.Admin) (string, error) {
	claims := jwt.MapClaims{
		"id":  admin.ID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(a.secret))
}

func (a *authService) RefreshAccessToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", err
	}

	adminID := uint(claims["id"].(float64))

	newClaims := jwt.MapClaims{
		"id":  adminID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	return newToken.SignedString([]byte(a.secret))
}

func (a *authService) GetPayload(tokenString string) (*dto.TokenPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	payload := &dto.TokenPayload{
		ID:     uint(claims["id"].(float64)),
		Email:  claims["email"].(string),
		IsRoot: claims["is_root"].(bool),
	}

	return payload, nil
}
