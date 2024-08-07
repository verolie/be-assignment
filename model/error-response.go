package model

import "time"

type ErrorResponse struct {
	Error_desc      string  	`json:"error_desc"`
	Exception 		string		`json:"exception"`
	TimeStamp 		time.Time 	`json:"timestamp"`
}
