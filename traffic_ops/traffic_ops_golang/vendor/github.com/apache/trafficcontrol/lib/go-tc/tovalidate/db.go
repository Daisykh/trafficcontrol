package tovalidate

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/apache/trafficcontrol/lib/go-log"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
)

// DBExistsRule checks if value is in column in this table
// TODO: DBExistsRule ?   what about DBUniqueRule?
type DBExistsRule struct {
	db      *sqlx.DB
	table   string
	column  string
	message string
}

// NewDBExistsRule is a validation rule that checks if the given value is in the column in this table
func NewDBExistsRule(db *sqlx.DB, table string, column string) *DBExistsRule {
	return &DBExistsRule{
		db:      db,
		table:   table,
		column:  column,
		message: fmt.Sprintf("No rows with value in %s.%s", table, column),
	}
}

// Validate checks if the given value is valid or not.
func (r *DBExistsRule) Validate(value interface{}) error {
	if r.db == nil {
		return nil
	}
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	query := `SELECT COUNT(*) FROM ` + r.table + ` WHERE ` + r.column + `= $1`
	row := r.db.QueryRow(query, value)
	var cnt int
	err := row.Scan(&cnt)
	log.Debugln("**** QUERY **** ", query)
	log.Debugf(" value %d err %++v", cnt, err)
	if err != nil {
		return errors.New(r.message)
	}
	return nil
}

// Error sets the error message for the rule.
func (r *DBExistsRule) Error(message string) *DBExistsRule {
	r.message = message
	return r
}

// DBUniqueRule checks if value is in column in this table
// TODO: DBUniqueRule ?   what about DBUniqueRule?
type DBUniqueRule struct {
	db      *sqlx.DB
	table   string
	column  string
	idCheck func(int) bool
	message string
}

// NewDBUniqueRule is a validation rule that checks if the given value is in the column in this table
func NewDBUniqueRule(db *sqlx.DB, table string, column string, idCheck func(int) bool) *DBUniqueRule {
	return &DBUniqueRule{
		db:      db,
		table:   table,
		column:  column,
		idCheck: idCheck,
		message: column + ` must be unique in ` + table,
	}
}

// Validate returns an error if the value already exists in the table in this column
func (r *DBUniqueRule) Validate(value interface{}) error {
	if r.db == nil {
		return nil
	}
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	query := `SELECT id FROM ` + r.table
	row := r.db.QueryRowx(query, map[string]interface{}{r.column: value})
	var id int
	err := row.Scan(&id)
	// ok if no rows found or only one belongs to row being updated
	if err == sql.ErrNoRows || r.idCheck(id) {
		return nil
	}
	return errors.New(r.message)
}

// Error sets the error message for the rule.
func (r *DBUniqueRule) Error(message string) *DBUniqueRule {
	r.message = message
	return r
}
