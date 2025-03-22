package usecase

import (
	"errors"

	"github.com/tbtec/tremligeiro-login/internal/dto"
)

// var (
// 	ErrorInvalidCredentials = xerrors.NewBusinessError("TLL-LOGIN-001", "Invalid Credentials")
// )

type UscLogin struct {
}

func NewUseCaseLogin() *UscLogin {
	return &UscLogin{}
}

func (u *UscLogin) Login(loginRequest dto.LoginRequest) (dto.Login, error) {

	if loginRequest.DocumentNumber == "123456" && loginRequest.Password == "123456" {
		return dto.Login{
			AccessToken: "1234567890",
		}, nil
	}

	return dto.Login{}, errors.New("invalid credentials")
}
