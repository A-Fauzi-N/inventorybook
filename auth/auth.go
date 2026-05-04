package auth

import (
	"inventorybook/models"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

// Redirect root ke login
func HomeHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/login")
}

// Menampilkan halaman login
func LoginGetHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"content": "",
	})
}

// Proses login
func LoginPostHandler(c *gin.Context) {
	var credential models.Login

	// Binding form
	if err := c.ShouldBind(&credential); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"content": "Username/Password request tidak valid",
		})
		return
	}

	// Validasi username & password
	if credential.Username != models.USER || credential.Password != models.PASSWORD {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"content": "Username/Password invalid",
		})
		return
	}

	// JWT Claims
	claim := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		Issuer:    "book_inventory",
		IssuedAt:  time.Now().Unix(),
	}

	// Generate token
	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	token, err := sign.SignedString([]byte(models.SECRET))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"content": "Gagal membuat token",
		})
		return
	}

	// Redirect ke /books?auth=token
	q := url.Values{}
	q.Set("auth", token)

	location := url.URL{
		Path:     "/books",
		RawQuery: q.Encode(),
	}

	c.Redirect(http.StatusFound, location.RequestURI())
}