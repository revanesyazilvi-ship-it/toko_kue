package repository

import (
	"database/sql"
	"toko_kue/models"
)

type OrderRepository interface {
	GetAll() ([]models.Pesanan, error)
	Create(pesanan *models.Pesanan) error
	UpdateStatus(id int, status string) error
	Delete(id int) error
	GetHistoryPelanggan() ([]models.HistoryPelanggan, error)
	GetStatistikPelanggan() (models.StatistikPelanggan, error)
}
type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetAll() ([]models.Pesanan, error) {
	rows, err := r.db.Query(`
		SELECT 
			id_pesanan,
			id_petugas,
			id_produk,
			qty,
			tanggal_pesanan,
			alamat,
			COALESCE(nama_pelanggan, ''),
			COALESCE(nomor_hp, ''),
			total_pesanan,
			COALESCE(metode_pembayaran, ''),
			status
		FROM pesanan
		ORDER BY tanggal_pesanan DESC
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pesanan []models.Pesanan

	for rows.Next() {
		var p models.Pesanan

		if err := rows.Scan(
			&p.IDPesanan,
			&p.IDPetugas,
			&p.IDProduk,
			&p.Qty,
			&p.TanggalPesanan,
			&p.Alamat,
			&p.NamaPelanggan,
			&p.NomorHP,
			&p.TotalPesanan,
			&p.MetodePembayaran,
			&p.Status,
		); err != nil {
			return nil, err
		}
		pesanan = append(pesanan, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pesanan, nil
}

func (r *orderRepository) Create(pesanan *models.Pesanan) error {
	query := `
		INSERT INTO pesanan (
			id_petugas,
			id_produk,
			qty,
			tanggal_pesanan,
			alamat,
			nama_pelanggan,
			nomor_hp,
			total_pesanan,
			metode_pembayaran,
			status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		pesanan.IDPetugas,
		pesanan.IDProduk,
		pesanan.Qty,
		pesanan.TanggalPesanan,
		pesanan.Alamat,
		pesanan.NamaPelanggan,
		pesanan.NomorHP,
		pesanan.TotalPesanan,
		pesanan.MetodePembayaran,
		pesanan.Status,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err == nil {
		pesanan.IDPesanan = int(id)
	}

	return nil
}

func (r *orderRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE pesanan SET status = ? WHERE id_pesanan = ?`

	_, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) GetHistoryPelanggan() ([]models.HistoryPelanggan, error) {
    query := `
        SELECT 
            COALESCE(nama_pelanggan, 'Tanpa Nama') as nama_pelanggan,
            COALESCE(nomor_hp, '-') as kontak,
            COALESCE(alamat, '-') as alamat,
            COUNT(id_pesanan) as total_transaksi,
            SUM(total_pesanan) as total_belanja
        FROM pesanan
        GROUP BY nama_pelanggan, nomor_hp, alamat
        ORDER BY total_belanja DESC
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []models.HistoryPelanggan
    for rows.Next() {
        var h models.HistoryPelanggan
        if err := rows.Scan(&h.NamaPelanggan, &h.Kontak, &h.Alamat, &h.TotalTransaksi, &h.TotalBelanja); err != nil {
            return nil, err
        }
        list = append(list, h)
    }

	if err := rows.Err(); err != nil {
    	return nil, err
	}
    return list, nil
}

func (r *orderRepository) GetStatistikPelanggan() (models.StatistikPelanggan, error) {
    query := `
        SELECT 
            COUNT(DISTINCT nama_pelanggan) as total_pelanggan,
            COUNT(id_pesanan) as total_transaksi,
            COALESCE(SUM(total_pesanan), 0) as total_pembelian
        FROM pesanan
    `
    var s models.StatistikPelanggan
    err := r.db.QueryRow(query).Scan(&s.TotalPelanggan, &s.TotalTransaksi, &s.TotalPembelian)
    if err != nil {
        return s, err
    }
    return s, nil
}

func (r *orderRepository) Delete(id int) error {
	query := `DELETE FROM pesanan WHERE id_pesanan = ?`

	result, err := r.db.Exec(query, id)
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