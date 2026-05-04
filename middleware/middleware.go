package middleware

import (
	"fmt"
	"inventorybook/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthValid(c *gin.Context) {
	var tokenString string

	// Ambil token dari query atau form
	tokenString = c.Query("auth")
	if tokenString == "" {
		tokenString = c.PostForm("auth")
	}

	// Jika token kosong
	if tokenString == "" {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"content": "Token tidak ditemukan, silakan login kembali",
		})
		c.Abort()
		return
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Validasi signing method harus HMAC
		if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("invalid token signing method: %v", token.Header["alg"])
		}

		return []byte(models.SECRET), nil
	})

	// Validasi hasil token
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"content": "Token expired atau tidak valid",
		})
		c.Abort()
		return
	}

	// Token harus valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token verified:", claims)
		c.Next()
		return
	}

	// Fallback invalid
	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"content": "Token tidak valid",
	})
	c.Abort()
}