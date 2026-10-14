package main

import (
	"log"
	"os"

	"toko_kue/config"
	"toko_kue/controller"
	"toko_kue/middleware"
	"toko_kue/policy"
	"toko_kue/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment OS")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret_key_default"
	}

	db := config.ConnectDB()
	defer db.Close()

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	r.Use(cors.New(cors.Config{
			AllowOrigins:  []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
   			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    		AllowCredentials: false,
}))

	storeRepo := repository.NewStoreRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderController := controller.NewOrderController(orderRepo)

	authCtrl := controller.NewAuthController(db)
	storeCtrl := controller.NewStoreController(db, storeRepo)
	produkCtrl := controller.NewProdukController(storeRepo)

	r.POST("/api/auth/register", authCtrl.Register)
	r.POST("/api/auth/login", authCtrl.Login)

	r.GET("/api/home/produk", produkCtrl.GetProduk)
	r.GET("/api/produk", produkCtrl.GetProduk)

	r.POST("/api/pesanan", orderController.CreateOrder)

	r.POST("/api/pembayaran", storeCtrl.ProcessPembayaran)
	r.GET("/api/history_penjualan", storeCtrl.GetHistoryPembayaran)
	r.GET("/api/history-pelanggan", orderController.GetHistoryPelangganHandler)

	r.GET("/api/store", storeCtrl.GetStore)
	r.PUT("/api/store", storeCtrl.UpdateStore)

	r.GET("/api/produk/favorit", storeCtrl.GetProdukFavoritHandler)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
			api.GET("/customer/history", storeCtrl.GetHistoryPembayaran)

			admin := api.Group("")
			admin.Use(middleware.RequireRole(policy.RoleAdmin))
		{
			admin.POST("/produk", produkCtrl.CreateProduk)
			admin.PUT("/produk/:id", produkCtrl.UpdateProduk)
			admin.DELETE("/produk/:id", produkCtrl.DeleteProduk)

			admin.GET("/pesanan", orderController.GetOrders)
admin.PUT("/pesanan/:id/status", orderController.UpdateStatusPesanan)
admin.DELETE("/pesanan/:id", orderController.DeleteOrder)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	for _, route := range r.Routes() {
	if route.Method == "DELETE" {
		log.Println("DELETE ROUTE:", route.Path)
	}
}

	r.Run("0.0.0.0:" + port)
}
