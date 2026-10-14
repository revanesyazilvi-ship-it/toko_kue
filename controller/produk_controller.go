package controller

import (
	"net/http"
	"strconv"
	"time"
	"fmt"
	"toko_kue/models"
	"toko_kue/repository"

	"github.com/gin-gonic/gin"
)

type ProdukController struct {
	storeRepo *repository.StoreRepository
}

func NewProdukController(storeRepo *repository.StoreRepository) *ProdukController {
	return &ProdukController{storeRepo: storeRepo}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (p *ProdukController) GetProduk(c *gin.Context) {
	produkList, err := p.storeRepo.GetProduk(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": produkList})
}

func (p *ProdukController) CreateProduk(c *gin.Context) { 
	namaProduk := c.PostForm("nama_produk")
	deskripsi := c.PostForm("deskripsi")
	hargaStr := c.PostForm("harga")
	stokStr := c.PostForm("stok")
	kategori := c.PostForm("kategori")

	harga, err := strconv.Atoi(hargaStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Harga tidak valid", 
		})
		return
	}
	
	stok, err := strconv.Atoi(stokStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stok tidak valid",
		})
		return
	}
	
	imageURL := ""
	
	file, err := c.FormFile("gambar")

	if err == nil {
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		filePath := "uploads/" + fileName

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file gambar"})
			return
		}
		imageURL = "http://localhost:8081/" + filePath
	}
	input := models.Produk{
		NamaProduk: namaProduk,
		Deskripsi:  deskripsi,
		Harga:      float64(harga),
		Stok:       stok,
		Kategori:   stringPtr(kategori),
		Gambar:     stringPtr(imageURL),
}

err = p.storeRepo.CreateProduk(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah produk: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Produk berhasil ditambahkan", "data": input})
}

func (p *ProdukController) UpdateProduk(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID produk tidak valid"})
		return
	}

	namaProduk := c.PostForm("nama_produk")
	deskripsi := c.PostForm("deskripsi")
	hargaStr := c.PostForm("harga")
	stokStr := c.PostForm("stok")
	kategori := c.PostForm("kategori")

	harga, err := strconv.Atoi(hargaStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Harga tidak valid",
			"harga_yang_diterima": hargaStr,
		})
		return
	}

	stok, err := strconv.Atoi(stokStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stok tidak valid",
			"stok_yang_diterima": stokStr,
		})
		return
	}

	imageURL := ""
	file, err := c.FormFile("gambar")
	if err == nil {
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		filePath := "uploads/" + fileName

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file gambar baru"})
			return
		}
		imageURL = "http://localhost:8081/" + filePath
	}

	input := models.Produk{
		NamaProduk: namaProduk,
		Deskripsi:  deskripsi,
		Harga:      float64(harga),
		Stok:       stok,
		Kategori:   stringPtr(kategori),
	}

	if imageURL != "" {
		input.Gambar = stringPtr(imageURL)
	}

	err = p.storeRepo.UpdateProduk(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diperbarui"})
}

func (p *ProdukController) DeleteProduk(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID produk tidak valid"})
		return
	}

	err = p.storeRepo.DeleteProduk(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}