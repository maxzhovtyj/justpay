package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

func NewConn(dsn string) (*pgx.Conn, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	// postgres://username:password@localhost:5432/database_name
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
