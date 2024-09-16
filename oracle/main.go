package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/godror/godror"
)

const createTableStatement = "CREATE TABLE TEMP_TABLE ( NAME VARCHAR2(100), CREATION_TIME TIMESTAMP DEFAULT SYSTIMESTAMP, VALUE  NUMBER(5))"
const dropTableStatement = "DROP TABLE TEMP_TABLE PURGE"
const insertStatement = "INSERT INTO TEMP_TABLE ( NAME , VALUE) VALUES (:name, :value)"
const queryStatement = "SELECT name, creation_time, value FROM TEMP_TABLE"

// func sqlOperations(db *sql.DB) {
// 	_, err := db.Exec(createTableStatement)
// 	handleError("create table", err)
// 	defer db.Exec(dropTableStatement) // make sure the table is removed when all is said and done
// 	stmt, err := db.Prepare(insertStatement)
// 	handleError("prepare insert statement", err)
// 	sqlresult, err := stmt.Exec("John", 42)
// 	handleError("execute insert statement", err)
// 	rowCount, _ := sqlresult.RowsAffected()
// 	fmt.Println("Inserted number of rows = ", rowCount)

// 	var queryResultName string
// 	var queryResultTimestamp time.Time
// 	var queryResultValue int32
// 	row := db.QueryRow(queryStatement)
// 	err = row.Scan(&queryResultName, &queryResultTimestamp, &queryResultValue)
// 	handleError("query single row", err)
// 	if err != nil {
// 		panic(fmt.Errorf("error scanning db: %w", err))
// 	}
// 	fmt.Println(fmt.Sprintf("The name: %s, time: %s, value:%d ", queryResultName, queryResultTimestamp, queryResultValue))
// 	_, err = stmt.Exec("Jane", 69)
// 	handleError("execute insert statement", err)
// 	_, err = stmt.Exec("Malcolm", 13)
// 	handleError("execute insert statement", err)

// 	// fetching multiple rows
// 	theRows, err := db.Query(queryStatement)
// 	handleError("Query for multiple rows", err)
// 	defer theRows.Close()
// 	var (
// 		name  string
// 		value int32
// 		ts    time.Time
// 	)
// 	for theRows.Next() {
// 		err := theRows.Scan(&name, &ts, &value)
// 		handleError("next row in multiple rows", err)
// 		fmt.Println(fmt.Sprintf("The name: %s and value:%d created at time: %s ", name, value, ts))
// 	}
// 	err = theRows.Err()
// 	handleError("next row in multiple rows", err)
// }

func handleError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

type Employee struct {
	ID   uint64
	Name string
	City string
}

func sqlOperations(db *sql.DB) {
	dbQuery := "SELECT * FROM EMPLOYEES"
	rows, err := db.Query(dbQuery)
	if err != nil {
		fmt.Println(".....Error processing query")
		fmt.Println(err)
		return
	}
	defer rows.Close()

	fmt.Println("... Parsing query results")
	for rows.Next() {
		var ID uint64
		var Name string
		var City string
		err := rows.Scan(&ID, &Name, &City)
		if err != nil {
			fmt.Printf("error scanning query result from database into target variable: %s \n ", err.Error())
			continue
		}

		employee := Employee{
			ID:   ID,
			Name: Name,
			City: City,
		}
		fmt.Printf("employee: %+v \n", employee)

	}
}

func GetSqlDBWithGoDrOrDriver(dbParams map[string]string) *sql.DB {
	var err error

	var P godror.ConnectionParams
	P.Username = dbParams["username"]
	P.Password = godror.NewPassword(dbParams["password"])
	P.ConnectString = fmt.Sprintf("%s:%s/%s?connect_timeout=2", dbParams["server"], dbParams["port"], dbParams["service"])
	P.SessionTimeout = 42 * time.Second
	P.SetSessionParamOnInit("NLS_NUMERIC_CHARACTERS", ",.")
	P.SetSessionParamOnInit("NLS_LANGUAGE", "FRENCH")
	P.Timezone = time.Local
	fmt.Println(P.StringWithPassword())
	db := sql.OpenDB(godror.NewConnector(P))
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	return db
}

var localDB = map[string]string{
	"service":  "ORCLPDB1",
	"server":   "103.161.38.151",
	"port":     "11521",
	"username": "testuser1",
	"password": "testuser1",
}

func main() {
	db := GetSqlDBWithGoDrOrDriver(localDB)
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("Can't close connection: ", err)
		}
	}()
	sqlOperations(db)
}
