package cookie

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

type firefoxDump struct {
	Dump
}

func (d *firefoxDump) Run() ([]*http.Cookie, error) {
	db, err := sql.Open("sqlite3", d.file)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`SELECT name, value, host, isHttpOnly FROM moz_cookies`)
	if err != nil {
		return nil, err
	}

	var cookies []*http.Cookie
	for rows.Next() {
		var name, value, host string
		var isHttpOnly int
		err = rows.Scan(&name, &value, &host, &isHttpOnly)
		if err != nil {
			return nil, err
		}

		var isHttpOk bool
		if isHttpOnly == 1 {
			isHttpOk = true
		}

		cookies = append(cookies, &http.Cookie{
			Name:     name,
			Value:    value,
			Domain:   host,
			HttpOnly: isHttpOk,
		})
	}

	return cookies, nil
}
