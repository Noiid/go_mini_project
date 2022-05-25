package model

type Mahasiswa struct {
	ID       int    `db:"id"`
	NIM      string `db:"nim"`
	Name     string `db:"name"`
	Gender   string `db:"gender"`
	Religion string `db:"religion"`
}
