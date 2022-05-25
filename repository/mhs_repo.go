package repository

import (
	"database/sql"
	"jwt/model"
)

type MhsRepo struct {
	db *sql.DB
}

func NewMhsRepository(db *sql.DB) *MhsRepo {
	return &MhsRepo{db}
}

func (r *MhsRepo) FetchMhs() ([]model.Mahasiswa, error) {
	preparedStatement := `
	SELECT 
	    *
	FROM 
		mahasiswas`

	rows, err := r.db.Query(preparedStatement)
	if err != nil {
		return nil, err
	}

	var mahasiswas []model.Mahasiswa

	for rows.Next() {
		var mhs model.Mahasiswa
		err := rows.Scan(&mhs.ID, &mhs.NIM, &mhs.Name, &mhs.Gender, &mhs.Religion)
		if err != nil {
			return nil, err
		}
		mahasiswas = append(mahasiswas, mhs)
	}

	return mahasiswas, nil
}

func (r *MhsRepo) CreateMhs(mhs model.Mahasiswa) (int64, error) {
	preparedStatement := `INSERT INTO mahasiswas (nim, name, gender, religion)
		VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(preparedStatement, mhs.NIM, mhs.Name, mhs.Gender, mhs.Religion)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *MhsRepo) UpdateMhs(mhs model.Mahasiswa) error {
	preparedStatement := `
		UPDATE mahasiswas
		SET name = ?, gender = ?, religion = ?
		WHERE nim = ?
	`
	_, err := r.db.Exec(preparedStatement, mhs.Name, mhs.Gender, mhs.Religion, mhs.NIM)
	if err != nil {
		return err
	}

	return nil
}

func (r *MhsRepo) DeleteMhsByNIM(nim string) error {
	preparedStatement := `DELETE FROM mahasiswas WHERE nim = ?`
	_, err := r.db.Exec(preparedStatement, nim)
	return err
}
