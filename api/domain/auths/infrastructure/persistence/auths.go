package persistence

import (
	"context"
	"fmt"
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/auths/entity"
	"github.com/JubaerHossain/cn-api/domain/auths/repository"
	userEntity "github.com/JubaerHossain/cn-api/domain/users/entity"
	"github.com/JubaerHossain/rootx/pkg/auth"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	utilQuery "github.com/JubaerHossain/rootx/pkg/query"
)

type AuthRepositoryImpl struct {
	app *app.App
}

// NewAuthRepository returns a new instance of AuthRepositoryImpl
func NewAuthRepository(app *app.App) repository.AuthRepository {
	return &AuthRepositoryImpl{
		app: app,
	}
}

// GetSignIn returns a new auth
func (r *AuthRepositoryImpl) GetSignIn(req *http.Request, loginUser *entity.LoginUser) (*entity.LoginUserResponse, error) {
	user := &userEntity.User{}
	if err := r.app.DB.QueryRow(context.Background(), "SELECT id, name, email, role, password FROM users WHERE email = $1", loginUser.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password); err != nil {
		return &entity.LoginUserResponse{}, fmt.Errorf("user not found")
	}

	if err := utilQuery.ComparePassword(user.Password, loginUser.Password); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	accessToken, refreshToken, err := auth.CreateTokens(user.ID, "", r.app, 24) // Call the CreateTokens function
	if err != nil {
		return nil, err
	}

	// Return the response
	return &entity.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetRefreshToken returns a new auth
func (r *AuthRepositoryImpl) GetRefreshToken(req *http.Request, reqRefreshToken *entity.RefreshToken) (*entity.LoginUserResponse, error) {
	// Implement logic to get refresh-token
	claims, err := auth.ValidateToken(reqRefreshToken.RefreshToken)
	if err != nil {
		return nil, err
	}

	userID := claims["sub"].(float64)
	user := &userEntity.User{}
	if err := r.app.DB.QueryRow(context.Background(), "SELECT id, role FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Role); err != nil {
		return &entity.LoginUserResponse{}, fmt.Errorf("user not found")
	}

	accessToken, refreshToken, err := auth.CreateTokens(user.ID, "", r.app, 24) // Call the CreateTokens function
	if err != nil {
		return nil, err
	}

	// Return the response
	return &entity.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
