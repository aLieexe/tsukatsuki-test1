package models

import (
	"context"
	"fmt"
	"go-webserver/utils"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func connectDb(t *testing.T) (*pgxpool.Pool, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		return nil, err
	}

	stringConnection := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v",
		utils.GetEnv("TESTPGUSER"), utils.GetEnv("TESTPGPASSWORD"), utils.GetEnv("TESTPGHOST"), utils.GetEnv("TESTPGPORT"), utils.GetEnv("TESTPGDATABASE"))
	config, err := pgxpool.ParseConfig(stringConnection)

	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	setupScript, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(context.Background(), string(setupScript))
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		defer db.Close()

		cleanupScript, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(context.Background(), string(cleanupScript))
		if err != nil {
			t.Fatal(err)
		}

	})

	return db, nil
}

func newTestDB(t *testing.T) *pgxpool.Pool {
	db, err := connectDb(t)
	if err != nil {
		t.Fatal(err)
	}

	return db
}
