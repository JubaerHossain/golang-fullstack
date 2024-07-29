package entity

import (
	"github.com/JubaerHossain/rootx/pkg/core/entity"
)

type LoginUser struct {
	Email    string `json:"email" validate:"required,min=11,max=15"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LoginUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type AuthUser struct {
	ID    uint        `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
	Role  entity.Role `json:"role"`
}
