package responsetype

import "github.com/akmalfairuz/finance/server/model"

type OTPInfo struct {
	ID    int64  `json:"otpId"`
	Token string `json:"token"`
}

func NewOTPInfoFromModel(otp model.OTP, rawTok string) OTPInfo {
	return OTPInfo{
		ID:    otp.ID,
		Token: rawTok,
	}
}
