package entity

import (
	"database/sql"
	"time"
)

type AvailableDriver struct {
	DriverID       string         `db:"driver_id"`
	Status         string         `db:"status"`
	LastSeenAt     time.Time      `db:"last_seen_at"`
	City           sql.NullString `db:"city"`
	Province       sql.NullString `db:"province"`
	JenisKendaraan sql.NullString `db:"jenis_kendaraan"`
	Nopol          sql.NullString `db:"nopol"`
}

type DriverInfo struct {
	DriverID       string `json:"driver_id" db:"driver_id"`
	FullName       string `json:"full_name" db:"full_name"`
	JenisKendaraan string `json:"jenis_kendaraan" db:"jenis_kendaraan"`
	Nopol          string `json:"nopol" db:"nopol"`
	City           string `json:"city" db:"city"`
}
