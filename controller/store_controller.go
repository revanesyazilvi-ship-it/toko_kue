package controller

import (
	"database/sql"
	"net/http"
	"time"

	"toko_kue/models"
	"toko_kue/repository"

	"github.com/gin-gonic/gin"
)

type StoreController struct {
	DB        *sql.DB
	storeRepo *repository.StoreRepository 
}

func NewStoreController(db *sql.DB, storeRepo *repository.StoreRepository) *StoreController {
	return &StoreController{
		DB: db,
		storeRepo: storeRepo,
	}
}

func (s *StoreController) GetProduk(c *gin.Context) {
	rows, err := s.DB.Query("SELECT id_produk, nama_produk, deskripsi, harga, stok FROM produk")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	list := make([]models.Produk, 0)
	for rows.Next() {
		var p models.Produk
		if err := rows.Scan(&p.IDProduk, &p.NamaProduk, &p.Deskripsi, &p.Harga, &p.Stok); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (s *StoreController) CreateProduk(c *gin.Context) {
	var input models.Produk
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := s.DB.Exec("INSERT INTO produk (nama_produk, deskripsi, harga, stok) VALUES (?, ?, ?, ?)",
		input.NamaProduk, input.Deskripsi, input.Harga, input.Stok)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	input.IDProduk = int(id)
	c.JSON(http.StatusCreated, gin.H{"message": "Produk berhasil ditambahkan", "data": input})
}

func (s *StoreController) CreatePesanan(c *gin.Context) {
	var input models.Pesanan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.IDPetugas == nil || *input.IDPetugas <= 0 {
		var exists bool
			err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM petugas WHERE id_petugas = ?)", input.IDPetugas).Scan(&exists)
			if err != nil || !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id_petugas tidak ditemukan pada data petugas"})
			return
		}

	}else{
		input.IDPetugas = nil
	}

	input.TanggalPesanan = time.Now()

	res, err := s.DB.Exec(`INSERT INTO pesanan (id_petugas, id_produk, qty, tanggal_pesanan, alamat, nama_pelanggan, nomor_hp, total_pesanan, metode_pembayaran, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.IDPetugas,
		input.IDProduk,
		input.Qty,
		input.TanggalPesanan,
		input.Alamat,
		input.NamaPelanggan,
		input.NomorHP,
		input.TotalPesanan,
		input.MetodePembayaran,
		input.Status,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	input.IDPesanan = int(id)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pesanan berhasil dibuat", 
		"data": input,
	})
}

func (s *StoreController) GetPesanan(c *gin.Context) {
	rows, err := s.DB.Query(`
		SELECT
			p.id_pesanan,
			p.id_petugas,
			p.id_produk,
			p.tanggal_pesanan,
			p.alamat,
			p.total_pesanan,
			pr.nama_produk
		FROM pesanan p
		LEFT JOIN produk pr ON p.id_produk = pr.id_produk
		ORDER BY p.id_pesanan DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	list := make([]gin.H, 0)

	for rows.Next() {
		var (
			idPesanan      int
			idPetugas      int
			idProduk       int
			tanggalPesanan time.Time
			alamat         string
			totalPesanan   int
			namaProduk     string
		)

		err := rows.Scan(
			&idPesanan,
			&idPetugas,
			&idProduk,
			&tanggalPesanan,
			&alamat,
			&totalPesanan,
			&namaProduk,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		list = append(list, gin.H{
			"id_pesanan":      idPesanan,
			"id_petugas":      idPetugas,
			"id_produk":       idProduk,
			"tanggal_pesanan": tanggalPesanan,
			"alamat":          alamat,
			"total_pesanan":   totalPesanan,
			"nama_produk":     namaProduk,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

func (s *StoreController) ProcessPembayaran(c *gin.Context) {
	var input models.Pembayaran
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi DB"})
		return
	}

	res, err := tx.Exec("INSERT INTO pembayaran (id_pesanan, metode_pembayaran, total) VALUES (?, ?, ?)",
		input.IDPesanan, input.MetodePembayaran, input.Total)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	idPembayaran, _ := res.LastInsertId()

	_, err = tx.Exec("INSERT INTO history_penjualan (id_pembayaran) VALUES (?)", idPembayaran)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pembayaran selesai dan dicatat ke History"})
}

func (s *StoreController) GetHistoryPembayaran(c *gin.Context) {
	histories, err := s.storeRepo.GetHistoryPembayaran(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": histories,
	})
}

func (s *StoreController) UpdateProduk(c *gin.Context) {
	id := c.Param("id")
	var req models.Produk
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := s.DB.Exec("UPDATE produk SET nama_produk = ?, harga = ? WHERE id_produk = ?", req.NamaProduk, req.Harga, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diperbarui"})
}

func (s *StoreController) DeleteProduk(c *gin.Context) {
	id := c.Param("id")
	_, err := s.DB.Exec("DELETE FROM produk WHERE id_produk = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}

func (s *StoreController) GetStore(c *gin.Context) {
    toko, err := s.storeRepo.GetStoreInfo(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": toko})
}

func (s *StoreController) UpdateStore(c *gin.Context) {
	var input models.Toko

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data pengaturan tidak valid",
		})
		return
	}

	if input.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID toko tidak boleh kosong",
		})
		return
	}

	err := s.storeRepo.UpdateStore(
		c.Request.Context(),
		input,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Data toko tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pengaturan toko berhasil diperbarui",
	})
}

func (s *StoreController) GetProdukFavoritHandler(c *gin.Context) {
    produk, err := s.storeRepo.GetProdukFavorit(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": produk})
}

func (s *StoreController) DeletePesanan(c *gin.Context) {
	id := c.Param("id")

	// Cek apakah pesanan ada
	var exists bool

	err := s.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pesanan WHERE id_pesanan = ?)",
		id,
	).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Pesanan tidak ditemukan",
		})
		return
	}

	// Hapus pesanan
	_, err = s.DB.Exec(
		"DELETE FROM pesanan WHERE id_pesanan = ?",
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pesanan berhasil dihapus",
	})
}