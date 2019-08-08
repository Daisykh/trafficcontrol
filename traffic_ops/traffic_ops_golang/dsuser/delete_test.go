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
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
)

func getTestDSUser() []tc.DeliveryServiceNullable {
	dsinfo := []tc.DeliveryServiceNullable{}
	testDSInfo := tc.DeliveryServiceNullable{}

	dsinfo = append(dsinfo, testDSInfo)

	return dsinfo
}

func TestDelete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// db := sql.Open(mockDB, "sqlmock")
	// defer db.Close()

	mock.ExpectBegin()
	testDSs := getTestDSUser()
	rows := sqlmock.NewRows([]string{"deliveryservice", "tm_user_id"}).AddRow(1, 2)

	mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM deliveryservice_tmuser WHERE deliveryservice = $1 AND tm_user_id = $2 RETURNING tm_user_id`)).WithArgs(1, 2).
		WillReturnResult(sqlmock.NewRows([]string{"tm_user_id"}).AddRow(2))

	mock.ExpectCommit()

}
