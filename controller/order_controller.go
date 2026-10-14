package controller

import (
		"net/http"
		"strconv"
		"time"
		"toko_kue/models"
		"toko_kue/repository"

		"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderRepo repository.OrderRepository
}

func NewOrderController(orderRepo repository.OrderRepository) *OrderController {
	return &OrderController{orderRepo: orderRepo}
}

func (oc *OrderController) GetOrders(c *gin.Context) {
	pesanan, err := oc.orderRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pesanan)
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var pesanan models.Pesanan
	if err := c.ShouldBindJSON(&pesanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if pesanan.TanggalPesanan.IsZero() {
		pesanan.TanggalPesanan = time.Now()
	}

	if err := oc.orderRepo.Create(&pesanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pesanan)
}

func (oc *OrderController) UpdateStatusPesanan(c *gin.Context) {
    idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
        return
    }

    var input struct {
        Status string `json:"status" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	err = oc.orderRepo.UpdateStatus(id, input.Status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status pesanan"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Status pesanan berhasil diperbarui",
    })
}

func (oc *OrderController) DeleteOrder(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID pesanan tidak valid",
		})
		return
	}

	err = oc.orderRepo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus pesanan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pesanan berhasil dihapus",
	})
}

func (oc *OrderController) GetHistoryPelangganHandler(c *gin.Context) {
    history, err := oc.orderRepo.GetHistoryPelanggan()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    statistik, err := oc.orderRepo.GetStatistikPelanggan()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{
        "statistik": statistik,
        "history":   history,
    })
}