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
	KeysStoreInterface
	Dsn          string
	SelectQuery  string
	pSelectQuery *sql.Stmt
	IndexQuery   string
	pIndexQuery  *sql.Stmt
	InsertQuery  string
	pInsertQuery *sql.Stmt
	UpdateQuery  string
	pUpdateQuery *sql.Stmt
	DeleteQuery  string
	pDeleteQuery *sql.Stmt
}

var (
	ErrInvalidInput = errors.New("invalid input provided")
	ErrSQLColumns   = errors.New("invalid columns returned by SQL query")
)

func prepareQuery(db *sql.DB, raw string, prepared *sql.Stmt) error {
	if raw == "" {
		return ErrMethodNotImplemented
	}

	var err error

	if prepared == nil {
		prepared, err = db.Prepare(raw)
		if err != nil {
			log.Printf("[ERROR] Failed to prepare SQL query: %s\n", err)
		}
	}

	return err
}

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

func (m *MysqlKeyStore) Get(identity ...interface{}) (KeyInstance, error) {
	key := KeyInstance{}

	if identity == nil {
		return key, ErrInvalidInput
	}

	db, err := sql.Open("mysql", m.Dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to MySQL server: %s\n", err)
		return key, err
	}

	if err = prepareQuery(db, m.SelectQuery, m.pSelectQuery); err != nil {
		return key, err
	}

	rows, err := m.pSelectQuery.Query(identity)
	if err != nil {
		log.Printf("[ERROR] Failed to execute SQL query: %s\n", err)
		return key, err
	}

	if res, _ := rows.Columns(); len(res) != 1 {
		return key, ErrSQLColumns
	}

	if !rows.Next() {
		return key, ErrKeyNotFound
	}

	var publicKey string
	var accessLevel int

	err = rows.Scan(&publicKey, &accessLevel)
	if err != nil {
		return key, err
	}

	key.Key = publicKey
	key.AccessLevel = accessLevel
	return key, nil
}

func (m *MysqlKeyStore) GetAll() (interface{}, error) {
	db, err := sql.Open("mysql", m.Dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to MySQL server: %s\n", err)
		return nil, err
	}

	if err = prepareQuery(db, m.IndexQuery, m.pIndexQuery); err != nil {
		return "", err
	}

	rows, err := m.pIndexQuery.Query()
	if err != nil {
		log.Printf("[ERROR] Failed to execute SQL query: %s\n", err)
		return nil, err
	}

	if res, _ := rows.Columns(); len(res) != 3 {
		return nil, ErrSQLColumns
	}

	var scanErr error
	keys := map[string]interface{}{}

	for rows.Next() {
		var (
			identity    string
			key         string
			accessLevel int
		)
		scanErr = rows.Scan(&identity, &key, &accessLevel)
		if scanErr != nil {
			break
		}

		keys[identity] = map[string]interface{}{
			"key":         key,
			"accessLevel": accessLevel,
		}
	}

	if scanErr != nil {
		log.Printf("[ERROR] Failed to retrieve row: %s\n", scanErr)
		return nil, scanErr
	}

	return keys, nil
}

func (m *MysqlKeyStore) Insert(values ...string) error {
	if len(values) < 2 {
		return ErrInvalidInput
	}

	db, err := sql.Open("mysql", m.Dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to MySQL server: %s\n", err)
		return err
	}

	if err = prepareQuery(db, m.InsertQuery, m.pInsertQuery); err != nil {
		return err
	}

	_, err = m.pInsertQuery.Exec(values)
	if err != nil {
		log.Printf("[ERROR] Failed to execute SQL query: %s\n", err)
		return err
	}

	return nil
}

func (m *MysqlKeyStore) Update(values ...string) error {
	if len(values) < 2 {
		return ErrInvalidInput
	}

	db, err := sql.Open("mysql", m.Dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to MySQL server: %s\n", err)
		return err
	}

	if err = prepareQuery(db, m.UpdateQuery, m.pUpdateQuery); err != nil {
		return err
	}

	_, err = m.pUpdateQuery.Exec(values)
	if err != nil {
		log.Printf("[ERROR] Failed to execute SQL query: %s\n", err)
		return err
	}

	return nil
}

func (m *MysqlKeyStore) Delete(identity ...interface{}) error {
	if identity == nil {
		return ErrInvalidInput
	}

	db, err := sql.Open("mysql", m.Dsn)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to MySQL server: %s\n", err)
		return err
	}

	if err = prepareQuery(db, m.DeleteQuery, m.pDeleteQuery); err != nil {
		return err
	}

	_, err = m.pDeleteQuery.Exec(identity)
	if err != nil {
		log.Printf("[ERROR] Failed to execute SQL query: %s\n", err)
		return err
	}

	return nil
}
