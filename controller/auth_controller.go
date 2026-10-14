package controller

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"toko_kue/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	DB *sql.DB
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{
		DB: db,
	}
}

func (a *AuthController) Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal enkripsi password"})
		return
	}

	var query string
	if input.Role == "admin" {
		query = "INSERT INTO admin (nama, password) VALUES (?, ?)"
	} else if input.Role == "petugas" {
		query = "INSERT INTO petugas (nama, password) VALUES (?, ?)"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role harus 'admin' atau 'petugas'"})
		return
	}

	res, err := a.DB.Exec(query, input.Nama, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil", "id": id})
}

func (a *AuthController) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	var storedPassword string
	var query string

	if input.Role == "admin" {
		query = "SELECT id_admin, password FROM admin WHERE nama = ?"
	} else if input.Role == "petugas" {
		query = "SELECT id_petugas, password FROM petugas WHERE nama = ?"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid"})
		return
	}

	err := a.DB.QueryRow(query, input.Nama).Scan(&id, &storedPassword)
	if err != nil {
		log.Println("Error SQL (Nama tidak ditemukan):", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nama tidak ditemukan"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(input.Password))
	if err != nil {
		log.Println("Error Bcrypt (Password salah):", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	secretEnv := os.Getenv("JWT_SECRET")
	if secretEnv == "" {
		secretEnv = "secret_key_default"
	}

	jwtSecret := []byte(secretEnv)

	claims := jwt.MapClaims{
		"user_id": id,
		"role":    input.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login berhasil", "token": tokenString})

}
