package service

import "github.com/gin-gonic/gin"

type Service interface {
	BalanceChange(ctx *gin.Context)
	UserBalance(ctx *gin.Context)
	MoneyTransaction(ctx *gin.Context)
}
