package repository

import (
	"database/sql"
	"errors"
	"toko_kue/models"

	_ "github.com/go-sql-driver/mysql"
)

type KueRepository struct {
		DB *sql.DB
}

func NewUserRepository(db *sql.DB) *KueRepository {
	return &KueRepository{DB: db}
}

func (r *KueRepository) CreateUser(input models.RegisterInput, hashedPassword string) (int64, error) {
	var query string
	if input.Role == "admin" {
		query = "INSERT INTO admin (nama, pasword) VALUES (?, ?)"
	} else if input.Role == "petugas" {
		query = "INSERT INTO petugas (nama, password) VALUES (?, ?)"
	} else {
		return 0, errors.New("role tidak valid")
	}

	res, err := r.DB.Exec(query, input.Nama, hashedPassword)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *KueRepository) GetUserByNama(nama string, role string) (int, string, error) {
	var id int
	var storedPassword string
	var query string

	if role == "admin" {
		query = "SELECT id_admin, pasword FROM admin WHERE nama = ?"
	} else if role == "petugas" {
		query = "SELECT id_petugas, password FROM petugas WHERE nama = ?"
	} else {
		return 0, "", errors.New("role tidak valid")
	}

	err := r.DB.QueryRow(query, nama).Scan(&id, &storedPassword)
	if err != nil {
		return 0, "", err
	}

	return id, storedPassword, nil
}
