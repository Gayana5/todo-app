package todo

import "errors"

type User struct {
	Id         int    `json:"-" db:"id"`
	FirstName  string `json:"first_name" db:"first_name" binding:"required"`
	SecondName string `json:"second_name" db:"second_name" binding:"required"`
	Email      string `json:"email" db:"email" binding:"required,email"`
	Password   string `json:"password" db:"password_hash" binding:"required"`
}

type UpdateUserInput struct {
	FirstName  *string `json:"first_name"`
	SecondName *string `json:"second_name"`
}

func (i UpdateUserInput) Validate() error {
	if i.FirstName == nil && i.SecondName == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
