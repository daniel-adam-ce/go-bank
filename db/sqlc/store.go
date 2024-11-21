package db

// import (
// 	"context"
// 	"database/sql"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// // Store provides all functinos to execute db queries and transactions
// type Store struct {
// 	*Queries
// 	db sql.DB
// }

// func NewStore(db *sql.DB) *Store {
// 	return &Store{
// 		db: db,
// 		Queries:  New(db),
// 	}
// }

// func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
// 		tx, err := store.execTx()
// }
