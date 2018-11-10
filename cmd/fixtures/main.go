package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/joho/godotenv"

	"github.com/shjp/shjp-dao"
)

// Order matters here..
var files = []string{
	"groups",
	"users",
	"roles",
	"events",
	"announcements",
	"comments",
	"groups_users",
	"groups_events",
	"groups_announcements",
	"users_events",
}

func insert(tx *pg.Tx, table string) error {
	file, err := os.Open(fmt.Sprintf("cmd/fixtures/data/%s.csv", table))
	if err != nil {
		log.Fatalf("Error reading file %s: %s", table, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		query := fmt.Sprintf(`
			INSERT INTO %s
			VALUES (%s)`,
			table, scanner.Text())

		log.Printf("query: %s\n\n", query)
		if _, err = tx.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	envVars, err := godotenv.Read()
	if err != nil {
		panic(err)
	}

	addr := envVars["ADDR"]
	user := envVars["USER"]
	dbName := envVars["DB"]
	password := envVars["PASSWORD"]

	db := dao.Init(&pg.Options{
		Addr:     addr,
		Password: password,
		User:     user,
		Database: dbName,
	})

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %s", err)
		os.Exit(1)
	}
	defer tx.Rollback()

	for _, name := range files {
		if err = insert(tx, name); err != nil {
			log.Printf("Error executing query: %s", err)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing the changes: %s", err)
		os.Exit(1)
	}
}
