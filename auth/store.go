/*
 * Copyright (c) Andrew Ying 2019.
 *
 * This file is part of the Intelligent Platform Management Interface (IPMI) software.
 * IPMI is licensed under the API Copyleft License. A copy of the license is available
 * at LICENSE.md.
 *
 * As far as the law allows, this software comes as is, without any warranty or
 * condition, and no contributor will be liable to anyone for any damages related
 * to this software or this license, under any kind of legal claim.
 */

package auth

import (
	"database/sql"
	"errors"
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
		value := store.FieldByName(key)
		if !value.IsValid() || value.Type().Name() != "string" || !value.CanSet() {
			continue
		}

		value.SetString(item)
	}
}

func (m *MysqlKeyStore) Get(identity ...interface{}) (string, error) {
	db, err := sql.Open("mysql", m.Dsn)
	if err != nil {
		return "", err
	}

	if m.query == nil {
		stmt, err := db.Prepare(m.RawQuery)
		if err != nil {
			return "", err
		}

		m.query = stmt
	}

	rows, err := m.query.Query(identity)
	if err != nil {
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
