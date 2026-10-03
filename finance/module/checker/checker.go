package checker

import (
	"encoding/json"
	"fmt"
	"github.com/akmalfairuz/finance/module/fetch"
)

func CheckGopay(phoneNumber string) (string, error) {
	return checkIdViaOrderKuota("gopay", phoneNumber)
}

func CheckOvo(phoneNumber string) (string, error) {
	return checkIdViaOrderKuota("ovo", phoneNumber)
}

func CheckShopeepay(phoneNumber string) (string, error) {
	return checkIdViaOrderKuota("shopeepay", phoneNumber)
}

func CheckDana(phoneNumber string) (string, error) {
	return checkIdViaOrderKuota("dana", phoneNumber)
}

func checkIdViaOrderKuota(id string, phoneNumber string) (string, error) {
	resp, err := fetch.PostForm("https://checker.orderkuota.com/api/checkid/"+id, map[string]string{
		"phoneNumber": phoneNumber,
	})
	if err != nil {
		return "", err
	}
	if resp.Status >= 400 {
		return "", fmt.Errorf("error: %d", resp.Status)
	}

	var resp2 struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(resp.Body, &resp2); err != nil {
		return "", err
	}

	if resp2.Status != "success" {
		return "", fmt.Errorf("error status: %s", resp2.Status)
	}

	if resp2.Message == phoneNumber {
		return "", nil
	}

	return resp2.Message, nil
}
