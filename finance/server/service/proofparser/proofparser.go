package proofparser

import (
	"github.com/akmalfairuz/finance/module/pointer"
	"strings"
)

var parsers = map[string]func(proof string) *string{}

func registerParser(name string, parser func(proof string) *string) {
	parsers[name] = parser
}

func init() {
	registerParser("token-pln", func(proof string) *string {
		// proof = 00000000000000000000/JOHN DOE/R1M/67.30/900
		// return 0000 0000 0000 0000 0000 (by cut the string by "/")
		// return nil if proof is not valid
		tok, _, found := strings.Cut(proof, "/")
		if !found {
			return nil
		}
		if strings.Contains(tok, "-") {
			tok = strings.ReplaceAll(tok, "-", "")
		}
		if strings.Contains(tok, " ") {
			tok = strings.ReplaceAll(tok, " ", "")
		}
		if len(tok) != 20 {
			return nil
		}
		splitted := make([]string, 0)
		for i := 0; i < len(tok); i += 4 {
			splitted = append(splitted, tok[i:i+4])
		}
		return pointer.Make(strings.Join(splitted, " "))
	})

	registerParser("google-play-idr-1", func(proof string) *string {
		// proof = Ref id: EXAMPLE0000001. KODE VOUCHER: EXAMPLEVOUCHER01
		// return EXAMPLEVOUCHER01 by get the last 16 characters, but verify
		// if contains Ref id: and KODE VOUCHER: first
		// return nil if proof is not valid
		if !strings.Contains(proof, "Ref id:") || !strings.Contains(proof, "KODE VOUCHER:") {
			return nil
		}
		return pointer.Make(proof[len(proof)-16:])
	})
}

func Parse(parser string, proof string) *string {
	if parser == "" {
		return nil
	}
	if p, ok := parsers[parser]; ok {
		if parsed := p(proof); p != nil {
			return parsed
		}
	}
	return nil
}
