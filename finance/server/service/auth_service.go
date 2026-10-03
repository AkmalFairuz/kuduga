package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/akmalfairuz/finance/module/hash"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
	"time"
)

type AuthService struct {
	log                 logrus.FieldLogger
	tokenRepository     *repository.TokenRepository
	googleOauthClientId string
}

var MaxTokenExpireDuration = time.Hour * 24 * 30

func NewAuthService(log logrus.FieldLogger, tokenRepository *repository.TokenRepository) *AuthService {
	return &AuthService{
		log:             log,
		tokenRepository: tokenRepository,
	}
}

func (s *AuthService) SetGoogleOauthClientId(id string) {
	s.googleOauthClientId = id
}

func (s *AuthService) CreateToken(user model.User, rawToken string) (model.Token, error) {
	hashedTok := hash.Sha256Bytes([]byte(rawToken))
	tokenId, err := s.tokenRepository.CreateToken(user.ID, user.Role, hashedTok, MaxTokenExpireDuration)
	if err != nil {
		return model.Token{}, err
	}
	token, err := s.tokenRepository.GetToken(tokenId)
	if err != nil {
		return model.Token{}, err
	}
	return token, nil
}

func (s *AuthService) CheckToken(tokenID int64, tokenValue string) (model.Token, error) {
	token, err := s.tokenRepository.GetToken(tokenID)
	if err != nil {
		return model.Token{}, err
	}
	if !bytes.Equal(hash.Sha256Bytes([]byte(tokenValue)), token.HashedToken) {
		return model.Token{}, errors.New("invalid token")
	}
	if time.Now().Unix() >= token.ExpiredAt {
		return model.Token{}, errors.New("token expired")
	}
	return token, nil
}

func (s *AuthService) RevokeToken(tokenID int64) error {
	return s.tokenRepository.DeleteToken(tokenID)
}

func (s *AuthService) GetTokens(userId int64) ([]model.Token, error) {
	return s.tokenRepository.GetTokens(userId)
}

func (s *AuthService) SetDeviceFcmToken(tokenID int64, deviceFcmToken string) error {
	return s.tokenRepository.SetDeviceFcmToken(tokenID, deviceFcmToken)
}

func (s *AuthService) ExtendToken(token model.Token) error {
	return s.tokenRepository.SetTokenExpiredAt(token.ID, time.Now().Add(MaxTokenExpireDuration).Unix())
}

func (s *AuthService) ValidateGoogleToken(idt string) (*model.GoogleUserInfo, error) {
	tok, err := idtoken.Validate(context.Background(), idt, s.googleOauthClientId)
	if err != nil {
		return nil, err
	}

	if !tok.Claims["email_verified"].(bool) {
		return nil, errors.New("email not verified")
	}

	return &model.GoogleUserInfo{
		Email:         tok.Claims["email"].(string),
		EmailVerified: tok.Claims["email_verified"].(bool),
		ProfileURL:    tok.Claims["picture"].(string),
		DisplayName:   tok.Claims["name"].(string),
	}, nil
}

func (s *AuthService) GetActiveUserCount(from, to int64) (int64, error) {
	return s.tokenRepository.GetActiveUserCount(from, to)
}
