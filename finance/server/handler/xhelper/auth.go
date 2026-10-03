package xhelper

import (
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/service"
	"strconv"
	"strings"
)

func GetToken(authService *service.AuthService, rawToken string) *model.Token {
	if rawToken == "" {
		return nil
	}
	// Example rawToken =
	// Authorization: Bearer 51432-Xi4dB2VIKr3tiAZnq6iJsArm4eDg19y0
	tokenType, rawToken, found := strings.Cut(rawToken, " ")
	if !found || tokenType != "Bearer" {
		return nil
	}

	tokenId, tokenValue, found := strings.Cut(rawToken, "-")
	if !found {
		return nil
	}

	tokenId2, err := strconv.Atoi(tokenId)
	if err != nil {
		return nil
	}

	tok, err := authService.CheckToken(int64(tokenId2), tokenValue)
	if err != nil {
		return nil
	}

	return &tok
}
