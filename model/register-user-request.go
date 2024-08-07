package model

type RegisUser struct {
	Acct_num      string `json:"acct_num"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	Password      string `json:"password"`
	DetailAccount DetailAccount
}