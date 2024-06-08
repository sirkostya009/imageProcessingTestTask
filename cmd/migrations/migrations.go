package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	. "imageProcessingTestTask/db"
	"io/fs"
	"os"
)

var db *pgx.Conn

func main() {
	var err error
	db, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(context.Background(), Schema)
	if err != nil {
		panic(err)
	}

	err = fs.WalkDir(Migrations, ".", func(path string, dir fs.DirEntry, _ error) error {
		if !dir.IsDir() {
			file, err := Migrations.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Println("running migration", path)

			_, err = db.Exec(context.Background(), string(file))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
