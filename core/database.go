package core

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectDb() *pgxpool.Pool {
	// load dotenv file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro got when loading .env")

	}
	ctx := context.Background()
	// load env and get database url
	databaseUrl := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		log.Fatal(err)

	}

	// cheack connection wth ping
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Db isn not Connected!")
	}

	// create all table is not exsits
	// defined alla tebles quires as list
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id  SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			profilePicture TEXT, 
			createdAt TIMESTAMPTZ NOT NULL DEFAULT now()	
	)`,
		`CREATE TABLE IF NOT EXISTS room (
			id  TEXT UNIQUE PRIMARY KEY,
			name TEXT NOT NULL,
			streamUrl TEXT NOT NULL,
			subtitleUrl TEXT,
			passcode TEXT,
			adminID  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			createdAt TIMESTAMPTZ NOT NULL DEFAULT now()	
	)`,
		`CREATE TABLE IF NOT EXISTS roomUser (
			userId   INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			roomId  TEXT NOT NULL REFERENCES room(id) ON DELETE CASCADE,
			joinedAt TIMESTAMPTZ NOT NULL DEFAULT now(),
			PRIMARY KEY (userId, roomId)
		
	)`,
	}

	for _, s := range stmts {
		if _, err := pool.Exec(ctx, s); err != nil {
			log.Fatalf("Table creating error : %v", err)
		}
	}

	// show database is connected
	fmt.Println("Database has connected!")

	return pool

}
