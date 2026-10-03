package model

const (
	OTPTypeEmail = iota
	OTPTypeSMS
)

const (
	OTPScopeRegister = iota
	OTPScopeResetPassword
	OTPScopeUpdateEmail
	OTPScopeRecoveryPin
)

type OTP struct {
	ID    int64 `db:"id"`
	Type  int   `db:"type"`
	Scope int   `db:"scope"`
	// Contact email/phone number
	Contact      string `db:"contact"`
	Code         string `db:"code"`
	HashedToken  []byte `db:"hashedToken"`
	FailAttempt  int64  `db:"failAttempt"`
	RetryAttempt int64  `db:"retryAttempt"`
	LastRetryAt  int64  `db:"lastRetryAt"`
	UsedAt       *int64 `db:"usedAt"`
	CreatedAt    int64  `db:"createdAt"`
	ExpiredAt    int64  `db:"expiredAt"`
}
