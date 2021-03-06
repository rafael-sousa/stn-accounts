// Package mysql_test includes assets required to run tests using dockertest and mysql implementation
package mysql_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/ory/dockertest/v3"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
)

var db *sql.DB
var txr repository.Transactioner

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	logFatal(err, "unable to create dockertest pool")

	dbCtnr, err := pool.Run("mysql", "5.6", []string{"MYSQL_ROOT_PASSWORD=secret"})
	logFatal(err, "unable to start db container")

	err = pool.Retry(func() error {
		var err error
		if db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?multiStatements=true&parseTime=true", dbCtnr.GetPort("3306/tcp"))); err == nil {
			return db.Ping()
		}
		return err
	})
	logFatal(err, "unable to connect to db container")

	runMigrations(pool)

	txr = repository.NewTxr(db)

	code := m.Run()

	err = pool.Purge(dbCtnr)
	logFatal(err, "unable to purge db container")

	os.Exit(code)
}

func runMigrations(pool *dockertest.Pool) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	logFatal(err, "unable to get driver instance")
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	logFatal(err, "unable to prepare migrations for execution")

	err = m.Up()
	logFatal(err, "unable to exec db migration")
}

func dbWipe() {
	_, err := db.Exec("DELETE FROM transfer")
	logFatal(err, "unable to clean the transfer table")

	_, err = db.Exec("DELETE FROM account")
	logFatal(err, "unable to clean the account table")
}

func logFatal(err error, msg string) {
	if err != nil {
		log.Fatalf("%s, %v", msg, err)
	}
}
