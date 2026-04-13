/**
* @description:
* @author
* @date 2026-03-25 07:57:59
* @version 1.0
*
* Change Logs:
* Date           Author       Notes
*
 */

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// SQLStore implements the Store interface and provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries:  New(connPool),
		connPool: connPool,
	}
}
