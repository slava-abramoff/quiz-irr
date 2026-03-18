package services

import (
	"fmt"
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
		"id":       admin.ID,
		"email":    admin.Email,
		"is_root":  admin.IsRoot,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
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
		return "", fmt.Errorf("invalid token")
	}

	idRaw, ok := claims["id"]
	if !ok {
		return "", fmt.Errorf("missing id claim")
	}
	idFloat, ok := idRaw.(float64)
	if !ok {
		return "", fmt.Errorf("invalid id claim type")
	}

	// Backward-compatible: older refresh tokens might not include these claims.
	email, _ := claims["email"].(string)
	isRoot, _ := claims["is_root"].(bool)

	newClaims := jwt.MapClaims{
		"id":      uint(idFloat),
		"email":   email,
		"is_root": isRoot,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
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
		return nil, fmt.Errorf("invalid token")
	}

	idRaw, ok := claims["id"]
	if !ok {
		return nil, fmt.Errorf("missing id claim")
	}
	idFloat, ok := idRaw.(float64)
	if !ok {
		return nil, fmt.Errorf("invalid id claim type")
	}

	// Backward-compatible: older tokens might not include all claims.
	email, _ := claims["email"].(string)
	isRoot, _ := claims["is_root"].(bool)

	payload := &dto.TokenPayload{
		ID:     uint(idFloat),
		Email:  email,
		IsRoot: isRoot,
	}

	return payload, nil
}
