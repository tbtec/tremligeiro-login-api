package usecase

import (
	"errors"

	"github.com/tbtec/tremligeiro-login/internal/dto"

	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	//"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go/aws/session"
)

// var (
// 	ErrorInvalidCredentials = xerrors.NewBusinessError("TLL-LOGIN-001", "Invalid Credentials")
// )

type UscLogin struct {
}

func NewUseCaseLogin() *UscLogin {
	return &UscLogin{}
}

var (
	userPoolID   = "us-east-1_ocwLNaBHg"        // User Pool ID
	clientID     = "3neaabp5n7tdsu43mt86d22meh" // App Client ID
	region       = "us-east-1"
	clientSecret = "1afe8p6bl46q788ph4f3edetivp4lk78a5s223ibkn6at34vgnme"
)

func generateSecretHash(clientID, clientSecret, username string) string {
	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func authenticateUser(username, password, secretHash string, sess *session.Session) (*cognitoidentityprovider.AdminInitiateAuthOutput, error) {

	svc := cognitoidentityprovider.New(sess)

	authParams := map[string]*string{
		"USERNAME":    aws.String(username),
		"PASSWORD":    aws.String(password),
		"SECRET_HASH": aws.String(secretHash),
	}

	authInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:       aws.String("ADMIN_NO_SRP_AUTH"),
		UserPoolId:     aws.String(userPoolID),
		ClientId:       aws.String(clientID),
		AuthParameters: authParams,
	}

	authOutput, err := svc.AdminInitiateAuth(authInput)
	if err != nil {
		log.Fatalf("Erro na autenticação: %s\n", err)
		return authOutput, errors.New("invalid credentials")
	}

	return authOutput, nil
}

func validateChallenge(username, password, secretHash string, sess *session.Session, authOutput *cognitoidentityprovider.AdminInitiateAuthOutput) (dto.Login, error) {

	if *authOutput.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		cognitoClient := cognitoidentityprovider.New(sess)

		input := &cognitoidentityprovider.RespondToAuthChallengeInput{
			Session:       aws.String(*authOutput.Session),
			ClientId:      aws.String(clientID),
			ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
			ChallengeResponses: map[string]*string{
				"NEW_PASSWORD": aws.String(password),
				"USERNAME":     aws.String(username),
				"SECRET_HASH":  aws.String(secretHash),
			},
		}

		_, err := cognitoClient.RespondToAuthChallenge(input)
		if err != nil {
			return dto.Login{}, err
		}

		return dto.Login{
			AccessToken: "Primeiro login efetuado",
		}, nil
	}

	return dto.Login{
		AccessToken: "Challenge não encontrado",
	}, nil

}

func (u *UscLogin) Login(loginRequest dto.LoginRequest) (dto.Login, error) {

	username := loginRequest.DocumentNumber
	password := loginRequest.Password

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatal("Erro ao criar sessão AWS:", err)
	}

	secretHash := generateSecretHash(clientID, clientSecret, username)

	authOutput, err := authenticateUser(username, password, secretHash, sess)

	if err != nil {
		log.Fatalf("Erro na autenticação: %s\n", err)
		return dto.Login{}, err
	}

	if authOutput.Session != nil {
		return validateChallenge(username, password, secretHash, sess, authOutput)
	} else {
		slog.Info("Nenhuma sessão encontrada na resposta")
	}

	return dto.Login{
		AccessToken: *authOutput.AuthenticationResult.AccessToken,
	}, nil

}
