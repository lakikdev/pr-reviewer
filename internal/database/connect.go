package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

var (
	databaseUserName   = flag.String("db-username", "postgres", "Database User.")
	databasePassword   = flag.String("db-password", "password", "Database Password.")
	databaseSocketPath = flag.String("db-socket-path", "", "Database Socket Path.")
	databaseHost       = flag.String("db-host", "localhost", "Database Host.")
	databasePort       = flag.String("db-port", "5432", "Database Port.")
	databaseName       = flag.String("db-name", "postgres", "Database Name.")
	databaseTimeout    = flag.Int64("database-timeout-ms", 10000, "")
)

type options struct {
	WithMigration bool
}

func NewOptions() options {
	return options{
		WithMigration: true, // Set the default value for WithMigration to true
	}
}

func New() (*DB, error) {
	return NewWithOptions(NewOptions())
}

// New creates a new database
func NewWithOptions(options options) (*DB, error) {
	conn, err := connect(options)
	if err != nil {
		return nil, err
	}

	return &DB{
		Conn: conn,
	}, nil
}

// Connect creates a new database connection
func connect(options options) (*sqlx.DB, error) {
	fmt.Print("Connecting to database\n")
	host := databaseSocketPath
	if *host == "" {
		host = databaseHost
	}

	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s sslmode=disable", *host, *databaseUserName, *databasePassword, *databasePort, *databaseName)
	// conn is the pool of database connections.
	conn, err := sqlx.Open("postgres", dbURI)
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to database")
	}

	conn.SetMaxOpenConns(64)

	if err := conn.Ping(); err != nil {
		return nil, errors.Wrap(err, "Could not connect to database")
	}

	//Check if database running
	if err := waitForDb(conn.DB); err != nil {
		return nil, err
	}

	if options.WithMigration {
		//Migrate database schema
		if err := migrateDb(conn.DB); err != nil {
			return nil, errors.Wrap(err, "could not migrate database")
		}
	}

	return conn, nil
}

func waitForDb(conn *sql.DB) error {
	ready := make(chan struct{})
	go func() {
		for {
			if err := conn.Ping(); err == nil {
				close(ready)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-ready:
		return nil
	case <-time.After(time.Duration(*databaseTimeout) * time.Millisecond):
		return errors.New("database not ready")
	}
}
