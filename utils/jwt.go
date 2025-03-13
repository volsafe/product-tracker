package utils

import (
	"errors"
	"fmt"
	"product-tracker/config"
	"time"

	"github.com/golang-jwt/jwt"
)

// Custom errors for better error handling
var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token has expired")
	ErrInvalidClaims        = errors.New("invalid token claims")
	ErrUserIDNotFound       = errors.New("user_id not found in token claims")
	ErrInvalidSigningMethod = errors.New("unexpected signing method")
	ErrConfigNotLoaded      = errors.New("configuration not loaded")
	ErrJWTSecretNotFound    = errors.New("JWT secret not found in config")
	ErrInvalidTokenFormat   = errors.New("invalid token format")
	ErrTokenNotBefore       = errors.New("token not yet valid")
)

// Claim keys for better maintainability
const (
	ClaimUserID = "user_id"
	ClaimExp    = "exp"
	ClaimIAT    = "iat"
	ClaimNBF    = "nbf"
	ClaimJTI    = "jti"
)

// TokenOptions contains options for token generation
type TokenOptions struct {
	ExpirationTime time.Duration
	NotBefore      time.Duration
	Issuer         string
	Audience       string
}

// DefaultTokenOptions returns default token options
func DefaultTokenOptions() TokenOptions {
	return TokenOptions{
		ExpirationTime: 24 * time.Hour,
		NotBefore:      0,
		Issuer:         "product-tracker",
		Audience:       "product-tracker-users",
	}
}

// TokenClaims represents the custom claims structure
type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	IAT      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
	NBF      int64  `json:"nbf,omitempty"`
	Issuer   string `json:"iss,omitempty"`
	Audience string `json:"aud,omitempty"`
	JTI      string `json:"jti,omitempty"`
}

// Valid implements the jwt.Claims interface
func (c *TokenClaims) Valid() error {
	now := time.Now().Unix()

	// Check if token has expired
	if c.Exp < now {
		return ErrTokenExpired
	}

	// Check if token was issued in the future
	if c.IAT > now {
		return errors.New("token issued in the future")
	}

	// Check if token is not yet valid
	if c.NBF > 0 && c.NBF > now {
		return ErrTokenNotBefore
	}

	return nil
}

// ExtractUserIDFromToken extracts the user ID from the JWT token and validates it
func ExtractUserIDFromToken(tokenString string) (uint, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*TokenClaims, error) {
	c := config.GetConfig()
	if c == nil {
		return nil, ErrConfigNotLoaded
	}

	secretKey := c.Jwt.Secret
	if secretKey == "" {
		return nil, ErrJWTSecretNotFound
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, ErrTokenExpired
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, ErrTokenNotBefore
			case ve.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return nil, ErrInvalidToken
			}
		}
		return nil, fmt.Errorf("token parsing error: %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// GenerateToken generates a new JWT token with the given user ID and options
func GenerateToken(userID uint, opts TokenOptions) (string, error) {
	c := config.GetConfig()
	if c == nil {
		return "", ErrConfigNotLoaded
	}

	secretKey := c.Jwt.Secret
	if secretKey == "" {
		return "", ErrJWTSecretNotFound
	}

	now := time.Now()
	claims := &TokenClaims{
		UserID:   userID,
		IAT:      now.Unix(),
		Exp:      now.Add(opts.ExpirationTime).Unix(),
		Issuer:   opts.Issuer,
		Audience: opts.Audience,
	}

	if opts.NotBefore > 0 {
		claims.NBF = now.Add(opts.NotBefore).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// RefreshToken generates a new token with the same claims but a new expiration time
func RefreshToken(oldToken string, opts TokenOptions) (string, error) {
	claims, err := ValidateToken(oldToken)
	if err != nil {
		return "", err
	}

	return GenerateToken(claims.UserID, opts)
}

// GetTokenExpiration returns the expiration time of a token
func GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(claims.Exp, 0), nil
}
