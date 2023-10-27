package models

type InvalidToken struct {
	Id        *int    `json:"id"`
	Token     *string `json:"token" validate:"required"`
	UpdatedAt *string `json:"updated_at" validate:"required"`
}
