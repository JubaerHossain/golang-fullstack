package repository

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/auths/entity"
)

// AuthRepository defines methods for auth data access
type AuthRepository interface {
	GetSignIn(req *http.Request, newUser *entity.LoginUser) (*entity.LoginUserResponse, error)
	GetRefreshToken(req *http.Request, refreshToken *entity.RefreshToken) (*entity.LoginUserResponse, error)
}
