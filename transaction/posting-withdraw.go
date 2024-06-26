package transaction

import (
	"context"
	"example/transaction/model"
	"example/transaction/prisma/db"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func WithdrawProcess(c *gin.Context) {
	var payment model.DetailTransaction
	var err error
	status = "Success"

	email, exists := c.Get("email")
	if !exists {
        ResponseErrorDetail(CreateErrorResp("Invalid Email" , err.Error()))
        return
    }

	if err = c.BindJSON(&payment); err != nil {
        c.JSON(http.StatusBadRequest,  ResponseErrorDetail(CreateErrorResp("Invalid request body" , err.Error())))
        return
    }

	client := db.NewClient()
    if err = client.Prisma.Connect(); err != nil {
        log.Fatalf("Error connecting to database: %v", err)
		return
    }
    defer client.Prisma.Disconnect()

	user, err := client.User.FindUnique(
        db.User.Email.Equals(email.(string)),
    ).Exec(context.Background())
	if  ((err != nil && err.Error() != "ErrNotFound") || user == nil) {
		c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Email Not Match" , err.Error())))
        return 
    }

	if(CheckAccount(client, payment.Loc_acct)){
		//try to update account
		if (actvTyp != "W") {
			UpdateAccount(c, client, payment);
			if status == "Failed" {
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error Account Already Write Off" , err.Error())))
        	return
		}

		//insert when payment success
    	_, err = client.TransactionDetail.CreateOne(
        	db.TransactionDetail.TrxID.Set(payment.Trx_id),
			db.TransactionDetail.Timestamps.Set(time.Now()),
			db.TransactionDetail.ReceiverPan.Set(""),
			db.TransactionDetail.SenderPan.Set(payment.Sender_pan),
			db.TransactionDetail.ApvCode.Set(payment.Apv_code),
			db.TransactionDetail.TrxTyp.Set(payment.Trx_typ),
			db.TransactionDetail.Amt.Set(payment.Amt),
			db.TransactionDetail.Status.Set(payment.Status),
			db.TransactionDetail.Desc.Set(payment.Desc),
			db.TransactionDetail.AcctDetail.Link(db.AccountDetail.LocAcct.Equals(payment.Loc_acct)),
    	).Exec(context.Background())
   		if err != nil {
        	c.JSON(http.StatusInternalServerError,   ResponseErrorDetail(CreateErrorResp("Error inserting payment data" , err.Error())))
        	return
    	}
	}else {
		c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error Cannot Find Data" , err.Error())))
        return
	}

	c.JSON(http.StatusOK, ResponseDataDetail("withdraw success"))
}
