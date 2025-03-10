package utils

import (
    "errors"
    "fmt"
    "product-tracker/config"
    "strconv"

    "github.com/golang-jwt/jwt"
)

func ExtractUserIDFromToken(tokenString string) (uint, error) {
    // Ensure the configuration is loaded
    c := config.GetConfig()
    if c == nil {
        return 0, errors.New("configuration not loaded")
    }

    secretKey := c.Jwt.Secret
    if secretKey == "" {
        return 0, errors.New("JWT secret not found in config")
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(secretKey), nil
    })

    // Extract and print user_id even if the token is invalid
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        if userIDStr, ok := claims["user_id"].(string); ok {
            userID, err := strconv.ParseFloat(userIDStr, 64)
            if err != nil {
                fmt.Printf("Error converting user_id to float64: %v\n", err)
            } else {
                fmt.Printf("Extracted user_id: %v\n", userID)
            }
        } else if userID, ok := claims["user_id"].(float64); ok {
            fmt.Printf("Extracted user_id: %v\n", userID)
        } else {
            fmt.Printf("user_id not found in token claims %v\n", claims)
        }
    }

    if err != nil {
        return 0, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        if userIDStr, ok := claims["user_id"].(string); ok {
            userID, err := strconv.ParseFloat(userIDStr, 64)
            if err != nil {
                return 0, errors.New("error converting user_id to float64")
            }
            return uint(userID), nil
        } else if userID, ok := claims["user_id"].(float64); ok {
            return uint(userID), nil
        }
    }

    return 0, errors.New("invalid token")
}