package repository

import (
	"context"
	"database/sql"

	"toko_kue/models"
)

type StoreRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) GetHistoryPembayaran(ctx context.Context) ([]models.HistoryResponse, error) {
	query := `
		SELECT 
			h.id_penjualan,
			h.id_pembayaran,
			IFNULL(p.id_pesanan, 0),
			IFNULL(p.metode_pembayaran, '-'),
			IFNULL(p.total, 0)
		FROM history_penjualan h
		LEFT JOIN pembayaran p 
			ON h.id_pembayaran = p.id_pembayaran
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []models.HistoryResponse

	for rows.Next() {
		var h models.HistoryResponse

		err := rows.Scan(
			&h.IDPenjualan,
			&h.IDPembayaran,
			&h.IDPesanan,
			&h.MetodePembayaran,
			&h.Total,
		)

		if err != nil {
			return nil, err
		}

		histories = append(histories, h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *StoreRepository) GetProduk(ctx context.Context) ([]models.Produk, error) {
	query := `
		SELECT
			id_produk,
			nama_produk,
			deskripsi,
			harga,
			stok,
			IFNULL(kategori, ''),
			IFNULL(gambar, ''),
			IFNULL(terjual, 0)
		FROM produk
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listProduk []models.Produk

	for rows.Next() {
		var p models.Produk

		err := rows.Scan(
			&p.IDProduk,
			&p.NamaProduk,
			&p.Deskripsi,
			&p.Harga,
			&p.Stok,
			&p.Kategori,
			&p.Gambar,
			&p.Terjual,
		)

		if err != nil {
			return nil, err
		}

		listProduk = append(listProduk, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listProduk, nil
}

func (r *StoreRepository) CreateProduk(ctx context.Context, produk models.Produk) error {
	query := `
		INSERT INTO produk (
			nama_produk,
			deskripsi,
			harga,
			stok,
			kategori,
			gambar,
			terjual
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		produk.NamaProduk,
		produk.Deskripsi,
		produk.Harga,
		produk.Stok,
		produk.Kategori,
		produk.Gambar,
		produk.Terjual,
	)

	return err
}

func (r *StoreRepository) UpdateProduk(ctx context.Context, id int, produk models.Produk) error {
	query := `
		UPDATE produk
		SET
			nama_produk = ?,
			deskripsi = ?,
			harga = ?,
			stok = ?,
			kategori = ?
		WHERE id_produk = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		produk.NamaProduk,
		produk.Deskripsi,
		produk.Harga,
		produk.Stok,
		produk.Kategori,
		id,
	)

	return err
}

func (r *StoreRepository) DeleteProduk(ctx context.Context, id int) error {
	query := `
		DELETE FROM produk
		WHERE id_produk = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

func (r *StoreRepository) GetStoreInfo(ctx context.Context) (models.Toko, error) {
	query := `
		SELECT
			id,
			nama_toko,
			email_toko,
			nomor_telepon,
			alamat_toko
		FROM toko
		ORDER BY id ASC
		LIMIT 1
	`

	var toko models.Toko

	err := r.db.QueryRowContext(ctx, query).Scan(
		&toko.ID,
		&toko.NamaToko,
		&toko.EmailToko,
		&toko.NomorTelepon,
		&toko.AlamatToko,
	)

	if err != nil {
		return toko, err
	}

	return toko, nil
}

func (r *StoreRepository) UpdateStore(ctx context.Context, toko models.Toko) error {
	query := `
		UPDATE toko
		SET nama_toko = ?,
			email_toko = ?,
			nomor_telepon = ?,
			alamat_toko = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		toko.NamaToko,
		toko.EmailToko,
		toko.NomorTelepon,
		toko.AlamatToko,
		toko.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *StoreRepository) GetProdukFavorit(ctx context.Context) ([]models.Produk, error) {
	query := `
		SELECT
			id_produk,
			nama_produk,
			deskripsi,
			harga,
			stok,
			IFNULL(kategori, ''),
			IFNULL(gambar, ''),
			IFNULL(terjual, 0)
		FROM produk
		WHERE is_favorite = 1
		LIMIT 4
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Produk

	for rows.Next() {
		var p models.Produk

		err := rows.Scan(
			&p.IDProduk,
			&p.NamaProduk,
			&p.Deskripsi,
			&p.Harga,
			&p.Stok,
			&p.Kategori,
			&p.Gambar,
			&p.Terjual,
		)

		if err != nil {
			return nil, err
		}

		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
