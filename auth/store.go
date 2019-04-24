/*
 * Adsisto
 * Copyright (c) 2019 Andrew Ying
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of version 3 of the GNU General Public License as published by the
 * Free Software Foundation. In addition, this program is also subject to certain
 * additional terms available at <SUPPLEMENT.md>.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
 * A PARTICULAR PURPOSE.  See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package auth

import (
	"database/sql"
	"errors"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlKeyStore struct {
	AuthorisedKeysInterface
	Dsn      string
	RawQuery string
	query    *sql.Stmt
}

var (
	ErrSQLColumns = errors.New("invalid columns returned by SQL query")
)

func (m *MysqlKeyStore) New(config map[string]string) {
	store := reflect.ValueOf(m)
	for key, item := range config {
		value := reflect.Indirect(store).FieldByName(key)
		if !value.IsValid() || value.Type().Name() != "string" || !value.CanSet() {
			continue
		}

		value.SetString(item)
	}
}

func (m *MysqlKeyStore) Get(identity ...interface{}) (string, error) {
	db, err := sql.Open("mysql", m.Dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to MySQL server: %s\n", err)
		return "", err
	}

	if m.query == nil {
		stmt, err := db.Prepare(m.RawQuery)
		if err != nil {
			log.Printf("[ERROR] Failed to prepare SQL query: %s\n", err)
			return "", err
		}

		m.query = stmt
	}

	rows, err := m.query.Query(identity)
	if err != nil {
		log.Printf("[ERROR] Failed to execute SQL query: %s\n", err)
		return "", err
	}

	if res, _ := rows.Columns(); len(res) != 1 {
		return "", ErrSQLColumns
	}

	rows.Next()

	var key string
	err = rows.Scan(&key)
	if err != nil {
		return "", err
	}

	return key, nil
}
