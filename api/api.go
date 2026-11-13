package api

import "github.com/jackc/pgx/v5/pgxpool"

func New(pool *pgxpool.Pool) *API {
	return &API{pool: pool}

}
