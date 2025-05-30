package services

import (
	"errors"
	"time"

	"KnowledgeHub/config"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims представляє claims для JWT токена
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type JWTService struct {
	config *config.Config
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		config: cfg,
	}
}

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token has expired")
	ErrInvalidClaims = errors.New("invalid token claims")
)

func (j *JWTService) GenerateTokenPair(userID uint, username, email string) (*TokenPair, error) {
	accessToken, expiresAt, err := j.generateAccessToken(userID, username, email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.generateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (j *JWTService) generateAccessToken(userID uint, username, email string) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(j.config.JWT.AccessTokenTTL) * time.Second)

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "KnowledgeHub",
			Subject:   "access_token",
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.config.JWT.SigningAlgorithm), claims)
	tokenString, err := token.SignedString([]byte(j.config.JWT.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

func (j *JWTService) generateRefreshToken(userID uint) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(j.config.JWT.RefreshTokenTTL) * time.Second)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "KnowledgeHub",
		Subject:   "refresh_token",
		ID:        string(rune(userID)),
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.config.JWT.SigningAlgorithm), claims)
	tokenString, err := token.SignedString([]byte(j.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != j.config.JWT.SigningAlgorithm {
			return nil, ErrInvalidToken
		}
		return []byte(j.config.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	// Перевіряємо, що це access токен
	if claims.Subject != "access_token" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshToken валідує refresh токен
func (j *JWTService) ValidateRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {

		if token.Method.Alg() != j.config.JWT.SigningAlgorithm {
			return nil, ErrInvalidToken
		}
		return []byte(j.config.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	if claims.Subject != "refresh_token" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (j *JWTService) RefreshTokens(refreshToken string, userID uint, username, email string) (*TokenPair, error) {

	_, err := j.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return j.GenerateTokenPair(userID, username, email)
}

func (j *JWTService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrInvalidToken
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", ErrInvalidToken
	}

	return authHeader[len(bearerPrefix):], nil
}
