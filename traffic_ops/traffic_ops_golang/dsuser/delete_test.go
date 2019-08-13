package dsuser

/*
 * LICENSED to the Apache Software Foundation (ASF) under one
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
	"context"
	"regexp"
	"testing"

	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDelete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM deliveryservice_tmuser WHERE deliveryservice = $1 AND tm_user_id = $2 RETURNING tm_user_id`)).WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"tm_user_id"}).AddRow(2))

	mock.ExpectCommit()

	dbCtx := context.TODO()
	tx, err := db.BeginTxx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	didDelete, err := deleteDSUser(tx.Tx, 1, 2)
	if err != nil || didDelete != true {
		t.Errorf("deleteDSUser expected: true - error nil, actual: %t - %v", didDelete, err)
	}
	if err = db.Close(); err != nil {
		t.Errorf("Error '%s' was not expected while closing the database", err)
	}

}
