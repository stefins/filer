package filer

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type File struct {
	Name    string `db:"NAME"`
	File_ID string `db:"FILE_ID"`
}

func (file *File) Insert(db *sqlx.DB) int {
	tx := db.MustBegin()
	defer tx.Commit()
	_, err := tx.Exec("INSERT INTO FILE(NAME,FILE_ID) VALUES ($1,$2)", file.Name, file.File_ID)
	if err != nil {
		log.Println("File Already Exists")
		return 0
	}
	return 1
}

func InitDB() *sqlx.DB {
	schema := `CREATE TABLE IF NOT EXISTS FILE (NAME text NOT NULL UNIQUE,FILE_ID text NOT NULL UNIQUE);`
	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot Ping the DB")
	}
	db.MustExec(schema)
	return db
}

func Search(term string, db *sqlx.DB) []File {
	var files []File
	err := db.Select(&files, "SELECT * FROM FILE WHERE NAME LIKE $1", "%"+term+"%")
	if err != nil {
		log.Println("Cannot find any record")
	}
	n := len(files)
	if n > 10 {
		return files[0:10]
	}
	return files[0:n]
}
