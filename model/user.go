package model

type User struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}
