package models

import "time"

type Admin struct {
	IDAdmin int    `json:"id_admin"`
	Nama    string `json:"nama"`
	Password string `json:"password"`
}

type Petugas struct {
	IDPetugas int    `json:"id_petugas"`
	Nama      string `json:"nama"`
	Password  string `json:"password"`
}

type Produk struct {
	IDProduk   int    `json:"id_produk"`
	NamaProduk string `json:"nama_produk"`
	Deskripsi  string `json:"deskripsi"`
	Harga      float64    `json:"harga"`
	Stok       int    `json:"stok"`
	Kategori  *string `json:"kategori"`
	Gambar    *string `json:"gambar"`
	Terjual    int    `json:"terjual"`
	IsFavorite bool	  `json:"is_favorite"` 
}

type Pesanan struct {
	IDPesanan        int       `json:"id_pesanan"`
    IDPetugas        *int       `json:"id_petugas"`
    IDProduk         int       `json:"id_produk"`
    Qty              int       `json:"qty"`
    TanggalPesanan   time.Time `json:"tanggal_pesanan"`
    Alamat           string    `json:"alamat"`
    NamaPelanggan    string    `json:"nama_pelanggan"`
    NomorHP          string    `json:"nomor_hp"`
    TotalPesanan     int       `json:"total_pesanan"`
    MetodePembayaran string    `json:"metode_pembayaran"`
    Status           string    `json:"status"`
}

type Pembayaran struct {
	IDPembayaran     int    `json:"id_pembayaran"`
	IDPesanan        int    `json:"id_pesanan"`
	MetodePembayaran string `json:"metode_pembayaran"`
	Total            float64 `json:"total"`
}

type RegisterInput struct {
	Nama     string `json:"nama" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type LoginInput struct {
	Nama     string `json:"nama" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type HistoryPembayaranResponse struct {
	IDHistory        int     `json:"id_history"`
	IDPembayaran     int     `json:"id_pembayaran"`
	IDPesanan        int     `json:"id_pesanan"`
	MetodePembayaran string  `json:"metode_pembayaran"`
	Total            float64 `json:"total"`
	Status           string  `json:"status"`
}

type HistoryResponse struct {
	IDPenjualan      int     `json:"id_penjualan"`
	IDPembayaran     int     `json:"id_pembayaran"`
	IDPesanan        int     `json:"id_pesanan"`
	MetodePembayaran string  `json:"metode_pembayaran"`
	Total            float64 `json:"total"`
}

type HistoryPelanggan struct {
    NamaPelanggan   string  `json:"nama_pelanggan"`
    Kontak          string  `json:"kontak"`
    Alamat          string  `json:"alamat"`
    TotalTransaksi  int     `json:"total_transaksi"`
    TotalBelanja    float64 `json:"total_belanja"`
}

type StatistikPelanggan struct {
    TotalPelanggan int     `json:"total_pelanggan"`
    TotalTransaksi int     `json:"total_transaksi"`
    TotalPembelian float64 `json:"total_pembelian"`
}

type Toko struct {
    ID          int    `json:"id"`
    NamaToko    string `json:"nama_toko"`
    EmailToko   string `json:"email_toko"`
    NomorTelepon string `json:"nomor_telepon"`
    AlamatToko  string `json:"alamat_toko"`
}