package responses

import "github.com/akmalfairuz/finance/server/handler/responsetype"

func M(message string) *responsetype.MessageResponse {
	return &responsetype.MessageResponse{Message: message}
}
