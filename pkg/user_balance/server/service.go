package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"mytest/pkg/user_balance/service"
	"mytest/repository/postgres"
	"net/http"
)

type ServiceBalance struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewServ(db *pgxpool.Pool, logger *zap.Logger) *ServiceBalance {

	return &ServiceBalance{
		db:     db,
		logger: logger,
	}
}

func (s *ServiceBalance) BalanceChange(ctx *gin.Context) {
	response := service.UserBalanceAddResponse{}
	req := service.UserBalanceAddRequest{}

	if err := ctx.BindJSON(&req); err != nil {
		s.logger.Error("can't parse json", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong parameters",
		})
		return
	}

	if req.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong id",
		})
		return
	}

	if req.Quality == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong quality",
		})
		return
	}

	balanceRepository := postgres.NewBalanceRepo(ctx, s.db)

	if err := balanceRepository.BalanceChange(req.Quality, req.Id); err != nil {
		response.Credited = false
		ctx.JSON(http.StatusInternalServerError, response)
		s.logger.Error("balance change error", zap.Error(err))
		return
	}
	response.Credited = true
	ctx.JSON(http.StatusOK, response)

}

func (s *ServiceBalance) MoneyTransaction(ctx *gin.Context) {
	response := service.UserBalanceTransactionResponse{}
	req := service.UserBalanceTransactionRequest{}

	if err := ctx.BindJSON(&req); err != nil {

		s.logger.Error("can't bind json", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong parameters",
		})
		return
	}

	if req.UserIdFrom <= 0 || req.UserIdTo <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong id",
		})
		return
	}
	if req.Quality <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong quality",
		})
		return
	}

	balanceRepository := postgres.NewBalanceRepo(ctx, s.db)

	if err := balanceRepository.BalanceTransaction(req.Quality, req.UserIdFrom, req.UserIdTo); err != nil {
		response.Transaction = false
		ctx.JSON(http.StatusOK, response)
		s.logger.Error("transaction error", zap.Error(err))
		return
	}

	response.Transaction = true
	ctx.JSON(http.StatusOK, response)

}

func (s *ServiceBalance) UserBalance(ctx *gin.Context) {

	req := service.UserBalanceRequest{}
	response := service.UserBalanceResponse{}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "wrong parameters",
		})
		s.logger.Error("can't bind json", zap.Error(err))
		return
	}

	if req.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong id",
		})
		return
	}

	balanceRepository := postgres.NewBalanceRepo(ctx, s.db)
	balance, err := balanceRepository.GetBalanceById(req.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		s.logger.Error("can't get balance", zap.Error(err))
		return
	}

	response.Balance = balance.Balance

	currencyParam := ctx.Request.URL.Query()
	if len(currencyParam) == 0 {
		ctx.JSON(http.StatusOK, response)
		//s.logger.Error("wrong parameters", zap.Error(err))
		return
	}
	if _, ok := currencyParam["currency"]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong currency parameter",
		})
		s.logger.Error("wrong currency parameter", zap.Error(err))
		return
	}
	respCurr, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": "try again later",
		})
		s.logger.Error("can't get currency from api", zap.Error(err))
		return
	}
	var resp service.Valute
	if err := json.NewDecoder(respCurr.Body).Decode(&resp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "try again later",
		})
		s.logger.Error("can't decode json", zap.Error(err))
		return
	}

	resultCurrencyName := currencyParam["currency"][0]
	if resp.ValuteMap[resultCurrencyName].Value == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong currency name",
		})
		s.logger.Error("can't find currency in currency names json api", zap.Error(err))
		return
	}
	ResultCurrencyValue := resp.ValuteMap[resultCurrencyName].Value
	response.Balance /= ResultCurrencyValue

	ctx.JSON(http.StatusOK, response)
}
