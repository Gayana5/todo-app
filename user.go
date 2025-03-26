package todo

type User struct {
	Id         int    `json:"-" db:"id"`
	FirstName  string `json:"first_name" db:"first_name" binding:"required"`
	SecondName string `json:"second_name" db:"second_name" binding:"required"`
	Email      string `json:"email" db:"email" binding:"required,email"`
	Password   string `json:"password" db:"password_hash" binding:"required"`
}
