package feature_postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func CheckConnection() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgres://postgres:pass@localhost:5432/task_manager")
	if err != nil {
		panic(err)
	}

	if err := conn.Ping(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")
}
