package middlewares

import (
    "net/http"
    "product-tracker/utils"
    "strings"
    "errors"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString, err := extractToken(c)
        if err != nil {
            abortWithError(c, http.StatusUnauthorized, err.Error())
            return
        }

        userID, err := utils.ExtractUserIDFromToken(tokenString)
        if err != nil {
            abortWithError(c, http.StatusUnauthorized, err.Error())
            return
        }

        // Store the user ID in the context for use in handlers
        c.Set("userID", userID)

        c.Next()
    }
}

func extractToken(c *gin.Context) (string, error) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        return "", errors.New("Missing Authorization header")
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        return "", errors.New("Invalid Authorization header format")
    }

    return tokenString, nil
}

func abortWithError(c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{"error": message})
    c.Abort()
}