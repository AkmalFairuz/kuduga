package errortype

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type VisibleError struct {
	code       int
	httpStatus int
	message    string
}

func (e VisibleError) Error() string {
	return fmt.Sprintf("[%d] %s", e.code, e.message)
}

func (e VisibleError) Code() int {
	return e.code
}

func (e VisibleError) Message() string {
	return e.message
}

func (e VisibleError) Print(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("X-Finance-Error-Code", strconv.Itoa(e.code))
	return ctx.Status(e.httpStatus).JSON(map[string]string{"message": e.message})
}

func NewErrorVisible(code int, httpStatus int, message string) VisibleError {
	return VisibleError{code: code, httpStatus: httpStatus, message: message}
}

var (
	CodeInvalidDestination = 200004
)

var (
	ErrBadRequest = NewErrorVisible(100000, http.StatusBadRequest, "Bad Request")

	ErrInvalidUsernameOrPassword = NewErrorVisible(100001, http.StatusUnauthorized, "Username atau password tidak valid")
	ErrTooManyLoginAttempt       = NewErrorVisible(100002, http.StatusTooManyRequests, "Terlalu banyak percobaan login")
	ErrAccountBlocked            = NewErrorVisible(100003, http.StatusUnauthorized, "Akun diblokir")
	ErrInvalidPin                = NewErrorVisible(100004, http.StatusUnauthorized, "PIN salah")
	ErrInvalidOTP                = NewErrorVisible(100005, http.StatusBadRequest, "OTP tidak valid")
	ErrInvalidOTPCode            = NewErrorVisible(100006, http.StatusUnauthorized, "OTP salah")
	ErrTooManyOTPRequest         = NewErrorVisible(100007, http.StatusTooManyRequests, "Terlalu banyak permintaan kode OTP. Coba lagi nanti.")
	ErrEmailNotRegistered        = NewErrorVisible(100008, http.StatusUnauthorized, "Email tidak terdaftar")

	ErrPurchaseBillNotFound        = NewErrorVisible(200002, http.StatusBadRequest, "Tagihan belum tersedia")
	ErrPurchaseProductNotAvailable = NewErrorVisible(200003, http.StatusBadRequest, "Produk sedang tidak tersedia")
	ErrPurchaseInvalidDestination  = NewErrorVisible(CodeInvalidDestination, http.StatusBadRequest, "Nomor tujuan tidak valid")
	ErrPurchaseBalanceLimit        = NewErrorVisible(200006, http.StatusBadRequest, "Melebihi limit penggunaan saldo")
	ErrPurchaseProductNotFound     = NewErrorVisible(200007, http.StatusNotFound, "Produk tidak ditemukan")
	ErrPurchaseBillExpired         = NewErrorVisible(200008, http.StatusBadRequest, "Tagihan kadaluarsa")
	ErrPurchaseNoAccess            = NewErrorVisible(200009, http.StatusForbidden, "Anda tidak memiliki akses ke pembelian ini")
	ErrPurchaseNotFound            = NewErrorVisible(200010, http.StatusNotFound, "Pembelian tidak ditemukan")
	ErrPurchaseDuplicate           = NewErrorVisible(200011, http.StatusBadRequest, "Anda baru saja membeli produk ini dengan tujuan yang sama pada 3 menit terakhir. Untuk mencegah terjadinya pembelian duplikat, pembelian ini tidak dapat dilakukan sekarang.")
	ErrPurchasePriceChange         = NewErrorVisible(200012, http.StatusBadRequest, "Harga produk berubah.")
	ErrPurchaseKycRequired         = NewErrorVisible(200013, http.StatusForbidden, "Akun anda harus terverifikasi untuk membeli produk ini. Lakukan verifikasi akun di Profil > Pengaturan Pribadi > Verifikasi Akun.")

	ErrTooManyDeposit = NewErrorVisible(210001, http.StatusBadRequest, "Anda memiliki deposit yang belum dibayar. Lakukan pembatalan atau pembayaran pada deposit yang belum dibayar.")

	ErrBalanceNotEnough      = NewErrorVisible(300001, http.StatusBadRequest, "Saldo tidak cukup")
	ErrUserTransactionLocked = NewErrorVisible(300002, http.StatusForbidden, "Kami tidak bisa memproses transaksi anda untuk sementara. Hubungi kami di Pusat Bantuan atau kirim email ke support@example.invalid untuk informasi lebih lanjut.")
	ErrReferralNotFound      = NewErrorVisible(300003, http.StatusNotFound, "Kode referral tidak ditemukan")
	ErrReferralAlreadyUsed   = NewErrorVisible(300004, http.StatusBadRequest, "Anda sudah menggunakan kode referral")
	ErrReferralSelf          = NewErrorVisible(300005, http.StatusBadRequest, "Anda tidak bisa menggunakan kode referral anda sendiri")
	ErrReferralUnable        = NewErrorVisible(300006, http.StatusBadRequest, "Anda tidak bisa menggunakan kode referral")
)
