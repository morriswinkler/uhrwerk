package main

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "log"
  "fmt"
  "time"
  "errors"
)

const TestTable string = "test_table"
type TestRecord struct {
  id int
  name, datetime string
}

type Database struct {
  DBHandle *sql.DB
  host, user, password, dbname string
}

// Initializes database handle, stores connection data
func (d *Database) Init(host, user, password, dbname string) *sql.DB {

  // Close old DBHandle if it exists
  if d.DBHandle != nil {
    d.DBHandle.Close()
  }

  // Store data in our Database struct
  d.host = host
  d.user = user
  d.password = password
  d.dbname = dbname

  // Get new database handle
  connData := fmt.Sprintf("%s:%s@%s/%s", 
    d.user, 
    d.password, 
    d.host, 
    d.dbname)
  
  db, err := sql.Open("mysql", connData)
  
  if err != nil {
    ERROR.Printf("Could not initialize db handle: %s", err)
    return nil
  }

  // Open doesn't open a connection. Validate DSN data:
  err = db.Ping()
  if err != nil {
    ERROR.Printf("Could not validate db connection: %s", err)
    return nil
  }

  d.DBHandle = db
  return d.DBHandle  
}

// returns the DB handle
func (d *Database) GetHandle() *sql.DB {
  return d.DBHandle
}

// Tests if the database can connect, create table, 
// insert records and get records
func (d *Database) Test() (bool, error){
  var err error

  if d.DBHandle == nil {
    err = errors.New("Missing db handle")
    ERROR.Println("Missing db handle")
    return false, err
  }

  db := d.DBHandle

  // Create test table
  query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s, %s, %s, %s)",
    TestTable,
    "id INT NOT NULL AUTO_INCREMENT",
    "name CHAR(30)",
    "datetime DATETIME",
    "PRIMARY KEY (id)")

  _, err = db.Exec(query)
  
  if err != nil {
    ERROR.Printf("Failed to create table: %s", err)
    return false, err
  }

  // Add current time as test record to the database
  timeNow := time.Now()
  timeLayout := "2006-01-02 15:04:05"
  timeString := timeNow.Format(timeLayout)
  //log.Printf("timeString: %s", timeString)
  query = fmt.Sprintf("INSERT INTO %s VALUES (%s, '%s', '%s')",
    TestTable,
    "NULL",
    "Tsar is mighty",
    timeString)

  _, err = db.Exec(query)
  
  if (err != nil) {
    ERROR.Printf("Failed to insert test data: %s", err)
    return false, err
  }

  // Get row to test if all works both ways
  var tr TestRecord
  query = fmt.Sprintf("SELECT * FROM %s", TestTable)
  err = db.QueryRow(query).Scan(&tr.id, &tr.name, &tr.datetime)
  
  if err != nil {
    log.Printf("Failed to get test record: %s", err)
    return false, err
  }

  //fmt.Printf("id: %d, name: %s, datetime: %s\n", tr.id, tr.name, tr.datetime)

  // Drop table
  query = fmt.Sprintf("DROP TABLE IF EXISTS %s", TestTable)
  _, err = db.Exec(query)
  if err != nil {
    log.Printf("Failed to drop table: %s", err)
    return false, err
  }

  return true, nil
}

// Close the connection to the database
func (d *Database) Close() {
  if d.DBHandle != nil {
    TRACE.Println("Closing db")
    d.DBHandle.Close()
  }
}