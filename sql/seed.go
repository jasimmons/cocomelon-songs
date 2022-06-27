package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

type Songs struct {
	Songs []Song `json:"songs"`
}

type Song struct {
	Name      string `json:"name"`
	Season    int    `json:"season"`
	Episode   int    `json:"episode"`
	StartTime string `json:"start_time"`
}

func (db *DB) Close() error {
	if db == nil {
		return nil
	}

	return db.DB.Close()
}

func Connect(dsn string, maxOpenConns, maxIdleConns int) (*DB, error) {
	log.Print("dsn=" + dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(60 * time.Second)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}

func main() {
	dsn := mustDSNFromEnv()
	db, err := Connect(dsn, 10, 10)
	if err != nil {
		log.Fatal("failed to connect to db: " + err.Error())
	}

	songs, err := loadFromFile("data.json")
	if err != nil {
		log.Fatal("failed to load songs from file: " + err.Error())
	}

	err = seed(db.DB, songs)
	if err != nil {
		log.Fatal(err)
	}
}

func mustDSNFromEnv() string {
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		cfg, err := mysql.ParseDSN(dsn)
		if err != nil {
			log.Fatal(err)
		}
		return cfg.FormatDSN()
	}

	const (
		db = "cocomelon"
	)
	var (
		username string
		password string
		host     string
		port     string
	)

	username = envOrFatal("DB_USERNAME")
	password = envOrFatal("DB_PASSWORD")
	host = envOrFatal("DB_HOSTNAME")
	port = envOrFatal("DB_PORT")

	cfg := mysql.NewConfig()
	cfg.User = username
	cfg.Passwd = password
	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(host, port)
	cfg.DBName = db

	return cfg.FormatDSN()
}

func envOrFatal(envvar string) string {
	fromEnv := os.Getenv(envvar)
	if fromEnv == "" {
		log.Fatal(fmt.Sprintf("missing %s", envvar))
	}
	return fromEnv
}

func loadFromFile(filename string) (Songs, error) {
	var songs Songs

	f, err := os.Open(filename)
	if err != nil {
		return songs, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	b, err := io.ReadAll(reader)
	if err != nil {
		return songs, err
	}

	err = json.Unmarshal(b, &songs)
	return songs, err
}

func seed(db *sql.DB, songs Songs) error {
	for _, song := range songs.Songs {
		_, err := sq.Insert("songs").
			Columns("name", "season", "episode", "start_time").
			Values(song.Name, song.Season, song.Episode, song.StartTime).
			Suffix("ON DUPLICATE KEY UPDATE name = ?, season = ?, episode = ?, start_time = ?",
				song.Name, song.Season, song.Episode, song.StartTime).
			RunWith(db).Exec()
		if err != nil {
			return err
		}
	}
	return nil
}
