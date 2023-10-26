package models

type User struct {
	Name     *string `json:"name" validate:"required"`
	Email    *string `json:"email" validate:"email,required"`
	Password *string `json:"password" validate:"required,min=4"`
	Phone    *string `json:"phone"`
	UserID   *int    `json:"user_id"`
	UserType *string `json:"user_type" validate:"oneof=ADMIN USER admin user"`
}
