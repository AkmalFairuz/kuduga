package proofparser

import (
	"github.com/akmalfairuz/finance/module/pointer"
	"github.com/go-playground/assert/v2"
	"testing"
)

// Synthetic vouchers and electricity tokens; never use customer payment data.
func TestProofParser(t *testing.T) {
	assert.Equal(t, pointer.Get(Parse("google-play-idr-1", "Ref id: EXAMPLE0000001. KODE VOUCHER: EXAMPLEVOUCHER01")), "EXAMPLEVOUCHER01")
	assert.Equal(t, pointer.Get(Parse("token-pln", "0000-0000-0000-0000-0000/JOHN-DOE/R1/1300/13,20.")), "0000 0000 0000 0000 0000")
	assert.Equal(t, pointer.Get(Parse("token-pln", "0000 0000 0000 0000 0000/JOHN-DOE/R1/1300/13,20.")), "0000 0000 0000 0000 0000")
}
