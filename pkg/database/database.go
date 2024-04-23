package database

import (
	"context"
	"fmt"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init() error {
	var err error
	iconfig := config.Get()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?TimeZone=Europe/Moscow",
		iconfig.PostgresUser,
		iconfig.PostgresPassword,
		iconfig.PostgresHost,
		iconfig.PostgresPort,
		iconfig.PostgresDatabase,
	)

	DB, err = pgxpool.New(context.Background(), dsn)

	return err
}
