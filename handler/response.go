package handler

import (
	"example/transaction/model"
	"time"
)

func ResponseDataDetail(data interface{}) model.DataStatusResponse {
	resp := model.DataStatusResponse{
		Status:  "Success",
		Message: "Success",
		Data:    data,
	}
	return resp
}

func ResponseErrorDetail(data interface{}) model.DataStatusResponse {
	resp := model.DataStatusResponse{
		Status:  "Error",
		Message: "Error",
		Data:    data,
	}
	return resp
}

func CreateErrorResp(desc string, err string) model.ErrorResponse {
	errorResponse := model.ErrorResponse{
   		Error_desc: desc,
   		Exception: err,
   		TimeStamp: time.Now(),
	}
	return errorResponse
}