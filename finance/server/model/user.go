package model

import (
	"io"
	"strings"
)

const (
	RoleMember = iota
	RolePremium
	RoleAdmin = 95
)

type User struct {
	ID           int64   `db:"id"`
	Name         string  `db:"name"`
	Email        string  `db:"email"`
	DisplayName  string  `db:"displayName"`
	PasswordHash []byte  `db:"passwordHash"`
	Balance      int64   `db:"balance"`
	PinHash      *string `db:"pinHash"`
	KycStatus    bool    `db:"kycStatus"`
	Role         int     `db:"role"`
	// Locked if true, user cannot do any transaction
	Locked bool `db:"locked"`

	// ReferralCode user's referral code
	ReferralCode *string `db:"referralCode"`
	// ReferralUserID who user referred to
	ReferralUserID *int64 `db:"referralUserId"`
	// ReferralAt when user referred
	ReferralAt *int64 `db:"referralAt"`

	DeviceUniqueID *string `db:"deviceUniqueId"`

	CreatedAt int64 `db:"createdAt"`
}

type CreateUser struct {
	Name           string
	Email          string
	DisplayName    string
	Password       string
	DeviceUniqueID *string
}

const (
	UserKycStatusPending = iota
	UserKycStatusSuccess
	UserKycStatusFailed
)

type UserKyc struct {
	ID             int64  `db:"id"`
	UserID         int64  `db:"userId"`
	DocumentID     string `db:"documentId"`
	FullName       string `db:"fullName"`
	DocumentFileID string `db:"documentFileId"`
	Status         int    `db:"status"`
	CreatedAt      int64  `db:"createdAt"`
}

type CreateUserKyc struct {
	UserID         int64
	DocumentID     string
	FullName       string
	DocumentFile   io.Reader
	DocumentFileID string
}

func IsValidUsername(name string) bool {
	valid := len(name) >= 3 && len(name) <= 20 && strings.ToLower(name) == name
	if !valid {
		return false
	}
	// check if name contains only alphanumeric characters
	for _, c := range name {
		if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}
