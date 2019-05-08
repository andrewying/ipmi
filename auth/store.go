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
	Dsn          string
	SelectQuery  string
	pSelectQuery *sql.Stmt
	IndexQuery   string
	pIndexQuery  *sql.Stmt
	UpdateQuery  string
	pUpdateQuery *sql.Stmt
	DeleteQuery  string
	pDeleteQuery *sql.Stmt
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

	if m.pSelectQuery == nil {
		stmt, err := db.Prepare(m.SelectQuery)
		if err != nil {
			log.Printf("[ERROR] Failed to prepare SQL query: %s\n", err)
			return "", err
		}

		m.pSelectQuery = stmt
	}

	rows, err := m.pSelectQuery.Query(identity)
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

func (m *MysqlKeyStore) GetAll() (map[string]string, error) {
	db, err := sql.Open("mysql", m.Dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to MySQL server: %s\n", err)
		return nil, err
	}

	if m.IndexQuery == "" {
		return nil, ErrMethodNotImplemented
	}

	if m.pIndexQuery == nil {
		stmt, err := db.Prepare(m.IndexQuery)
		if err != nil {
			log.Printf("[ERROR] Failed to prepare SQL query: %s\n", err)
			return nil, err
		}

		m.pIndexQuery = stmt
	}

	rows, err := m.pIndexQuery.Query()
	if err != nil {
		log.Printf("[ERROR] Failed to execute SQL query: %s\n", err)
		return nil, err
	}

	if res, _ := rows.Columns(); len(res) != 2 {
		return nil, ErrSQLColumns
	}

	var scanErr error
	keys := map[string]string{}

	for rows.Next() {
		var identity, key string
		scanErr = rows.Scan(&identity, &key)
		if scanErr != nil {
			break
		}

		keys[key] = identity
	}

	if scanErr != nil {
		log.Printf("[ERROR] Failed to retrieve row: %s\n", scanErr)
		return nil, scanErr
	}

	return keys, nil
}

func (m *MysqlKeyStore) Update(identity, publicKey string) error {
	db, err := sql.Open("mysql", m.Dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to MySQL server: %s\n", err)
		return err
	}

	if m.UpdateQuery == "" {
		return ErrMethodNotImplemented
	}

	if m.pUpdateQuery == nil {
		stmt, err := db.Prepare(m.UpdateQuery)
		if err != nil {
			log.Printf("[ERROR] Failed to prepare SQL query: %s\n", err)
			return err
		}

		m.pUpdateQuery = stmt
	}

	_, err = m.pUpdateQuery.Exec(identity, publicKey)
	if err != nil {
		log.Printf("[ERROR] Failed to execute SQL query: %s\n", err)
		return err
	}

	return nil
}

func (m *MysqlKeyStore) Delete(identity ...interface{}) error {
	db, err := sql.Open("mysql", m.Dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to MySQL server: %s\n", err)
		return err
	}

	if m.DeleteQuery == "" {
		return ErrMethodNotImplemented
	}

	if m.pDeleteQuery == nil {
		stmt, err := db.Prepare(m.DeleteQuery)
		if err != nil {
			log.Printf("[ERROR] Failed to prepare SQL query: %s\n", err)
			return err
		}

		m.pDeleteQuery = stmt
	}

	_, err = m.pDeleteQuery.Exec(identity)
	if err != nil {
		log.Printf("[ERROR] Failed to execute SQL query: %s\n", err)
		return err
	}

	return nil
}
