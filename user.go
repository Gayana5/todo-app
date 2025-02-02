package todo

type User struct {
	Id          int    `json:"-" db:"id"`
	First_Name  string `json:"first_name" binding:"required"`
	Second_Name string `json:"second_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
}
