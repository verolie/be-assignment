package model

type DetailTransaction struct {
	Trx_id       string  `json:"trx_id"`
	Apv_code     string  `json:"apv_code"`
	Trx_typ      string  `json:"trx_typ"`
	Amt          float64 `json:"amt"`
	Status       string  `json:"status"`
	Desc         string  `json:"desc"`
	Loc_acct     string  `json:"loc_acct"`
	Sender_pan   string  `json:"sender_pan"`
	Receiver_pan string  `json:"receiver_pan"`
}