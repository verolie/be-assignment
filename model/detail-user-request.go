package model

import (
	_ "github.com/shopspring/decimal"
)

type User struct {
	Acct_num      string `json:"acct_num"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	DetailAccount DetailAccount
}

type DetailAccount struct {
	Loc_acct       string  `json:"loc_acct"`
	Prin_pan       string  `json:"prin_pan"`
	Acct_typ       string  `json:"acct_typ"`
	Actv_typ       string  `json:"actv_typ"`
	Blnc_amt       float64 `json:"blnc_amt"`
	Loan_amt       float64 `json:"loan_amt"`
	Cycc_day       int     `json:"cycc_day"`
	Min_loan_pymnt float64 `json:"min_loan_pymnt"`
}
