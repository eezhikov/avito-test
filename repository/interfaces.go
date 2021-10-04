package repository

import (
	"mytest/models"
)

type BalanceRepository interface {
	GetBalanceById(int) (*models.Balance, error)
	BalanceChange(float64, int) error
	BalanceTransaction(float64, int, int) error
}
