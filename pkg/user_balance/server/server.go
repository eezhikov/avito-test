package server

import (
	"github.com/gin-gonic/gin"
	"mytest/pkg/user_balance/service"
)

func NewGin(srv service.Service) *gin.Engine {
	r := gin.Default()

	r.POST("/change", srv.BalanceChange)
	r.POST("/transaction", srv.MoneyTransaction)
	r.POST("/balance", srv.UserBalance)

	return r

}
