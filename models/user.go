package models

type User struct {
	Name     *string `json:"name" validate:"required"`
	Email    *string `json:"email" validate:"email,required"`
	Password *string `json:"password" validate:"required,min=4"`
	Phone    *string `json:"phone"`
	BookID   *int16  `json:"book_id"`
	UserID   *int16  `json:"user_id"`
}

type Book struct {
	BookID    *int16  `json:"book_id"`
	Author    *string `json:"author" validate:"required"`
	Publisher *string `json:"publisher"`
	Title     *string `json:"title" validate:"required"`
	UserID    *int16  `json:"user_id" validate:"required"`
}
