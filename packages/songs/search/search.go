package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Song struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	Season    int    `json:"season"`
	Episode   int    `json:"episode"`
	StartTime string `json:"start_time"`
}

func Main(args map[string]interface{}) map[string]interface{} {
	var term string

	if t, ok := args["term"].(string); ok {
		term = t
	}

	if term == "" {
		return map[string]interface{}{
			"errors": "empty search term",
		}
	}

	db, err := dbFromEnv()
	if err != nil {
		return map[string]interface{}{
			"errors": err.Error(),
		}
	}

	songs, err := getSongsMatchingTerm(db, term)
	if err != nil {
		return map[string]interface{}{
			"errors": err.Error(),
		}
	}

	if len(songs) == 0 {
		return map[string]interface{}{
			"body": "[]",
		}
	}

	songsJson, err := json.Marshal(songs)
	if err != nil {
		return map[string]interface{}{
			"errors": err.Error(),
		}
	}

	return map[string]interface{}{
		"body": string(songsJson),
	}
}

func dbFromEnv() (*sql.DB, error) {
	user := os.Getenv("DB_USERNAME")
	if user == "" {
		// default DO mysql user
		user = "doadmin"
	}

	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		return nil, errors.New("missing DB_PASSWORD")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		return nil, errors.New("missing DB_HOST")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		// default DO mysql port
		port = "25060"
	}

	db := os.Getenv("DB_DATABASE")
	if db == "" {
		db = "cocomelon"
	}

	cfg := &mysql.Config{
		User:   user,
		Passwd: pass,
		Net:    "tcp",
		Addr:   net.JoinHostPort(host, port),
		DBName: db,
	}
	dsn := cfg.FormatDSN()

	return sql.Open("mysql", dsn)
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func getSongsMatchingTerm(db *sql.DB, term string) ([]Song, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := sq.Select("name", "season", "episode", "start_time").
		From("songs").
		Where(sq.Like{"name": fmt.Sprintf("%%%s%%", term)}).
		OrderBy("season ASC, episode ASC, start_time ASC")

	rows, err := query.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		song, err := scanSong(rows)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return songs, nil
}

func scanSong(row scannable) (Song, error) {
	var s Song
	err := row.Scan(
		&s.Name,
		&s.Season,
		&s.Episode,
		&s.StartTime,
	)
	return s, err
}
