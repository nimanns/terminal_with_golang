package utils

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
    "time"
)

var jwtKey = []byte("secret_key")

func HashPassword(password string) string {
    bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func GenerateToken(userID uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    })
    return token.SignedString(jwtKey)
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
            c.Abort()
            return
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
            c.Abort()
            return
        }

        userID := uint(claims["user_id"].(float64))
        c.Set("userID", userID)
        c.Next()
    }
}
