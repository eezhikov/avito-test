package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"mytest/models"
	"mytest/repository"
)

type balanceRepoPostgres struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewBalanceRepo(ctx context.Context, db *pgxpool.Pool) repository.BalanceRepository {
	return &balanceRepoPostgres{
		db:  db,
		ctx: ctx,
	}
}

func (b *balanceRepoPostgres) GetBalanceById(id int) (*models.Balance, error) {
	var balance sql.NullFloat64

	if err := b.db.QueryRow(b.ctx, "SELECT balance FROM balance WHERE id=$1", id).Scan(&balance); err != nil {
		return nil, err
	}
	return &models.Balance{Id: id, Balance: balance.Float64}, nil
}
func (b *balanceRepoPostgres) BalanceChange(quality float64, id int) error {

	tx, err := b.db.BeginTx(b.ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(b.ctx)
	var balance float64

	if err := tx.QueryRow(b.ctx, "SELECT balance FROM balance WHERE id=$1", id).Scan(&balance); err != nil {
		return err
	}
	newBalance := balance + quality
	if newBalance < 0 {
		return errors.New("balance is lower than quality, your change not made")
	}
	if _, err := tx.Exec(b.ctx, "UPDATE balance SET balance = $1 WHERE id = $2", newBalance, id); err != nil {
		return err
	}

	if err := tx.Commit(b.ctx); err != nil {
		return err
	}
	return err
}
func (b *balanceRepoPostgres) BalanceTransaction(quality float64, userIdFrom int, userIdTo int) error {
	tx, err := b.db.BeginTx(b.ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(b.ctx)

	rows, err := tx.Query(b.ctx, "SELECT id, balance FROM balance WHERE id=$1 OR id=$2", userIdFrom, userIdTo)
	if err != nil {
		fmt.Println(err)
	}
	userBalance := make(map[int]float64)

	for rows.Next() {
		var userIdRow int
		var userBalanceRow float64
		if err := rows.Scan(&userIdRow, &userBalanceRow); err != nil {
			fmt.Println(err)
			return err
		}
		userBalance[userIdRow] = userBalanceRow

	}
	if userBalance[userIdFrom]-quality < 0 {
		return errors.New("balance is lower than quality, your change not made")
	}

	if _, err = tx.Exec(b.ctx, "UPDATE balance SET balance = $1 WHERE id = $2", userBalance[userIdFrom]-quality, userIdFrom); err != nil {
		return err
	}

	if _, err = tx.Exec(b.ctx, "UPDATE balance SET balance = $1 WHERE id = $2", userBalance[userIdTo]+quality, userIdTo); err != nil {
		return err
	}

	if err = tx.Commit(b.ctx); err != nil {
		return err
	}
	return err
}
