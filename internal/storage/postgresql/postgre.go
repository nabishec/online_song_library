package postgresql

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/nabishec/restapi/internal/storage/postgresql/migration"
)

type Database struct {
	dataSourceName string
	DB             *sqlx.DB
}

type dataSourceName struct {
	protocol string
	userName string
	password string
	host     string
	port     string
	dbName   string
	options  string
}

func NewDatabase() (*Database, error) {
	var database Database
	dbConnfig, err := NewDSN()
	if err != nil {
		return nil, err
	}

	err = database.connectDatabase(dbConnfig)
	if err != nil {
		return nil, err
	}
	err = migration.MigrationsUp(database.DB, database.dataSourceName)
	if err != nil {
		return nil, err
	}

	return &database, err
}

func (db *Database) connectDatabase(config *dataSourceName) error {
	const op = "internal.storage.postgresql.ConnectDatabase()"

	db.dataSourceName = config.protocol + "://" + config.userName + ":" + config.password + "@" +
		config.host + ":" + config.port + "/" + config.dbName + "?" + config.options

	var connectError error
	db.DB, connectError = sqlx.Connect("pgx", db.dataSourceName)
	if connectError != nil {
		return fmt.Errorf("%s,%w", op, connectError)
	}
	return connectError
}

func (db *Database) PingDatabase() error {
	const op = "internal.storage.postgresql.PingDatabase()"

	if db.DB == nil {
		return fmt.Errorf("%s:%s", op, "database isn`t established")
	}

	var pingError = db.DB.Ping()
	if pingError != nil {
		return fmt.Errorf("%s,%w", op, pingError)
	}
	return nil
}

func (db *Database) CloseDatabase() error {
	const op = "internal.storage.postgresql.CloseDatabase()"

	var closingError = db.DB.Close()
	if closingError != nil {
		return fmt.Errorf("%s:%w", op, closingError)
	}
	return nil
}

func NewDSN() (*dataSourceName, error) {
	const op = "internal.storage.postgresql.NewDSN()"
	dsn := &dataSourceName{}

	dsn.protocol = os.Getenv("DB_PROTOCOL")
	if dsn.protocol == "" {
		return nil, fmt.Errorf("%s:%s", op, "DB_PROTOCOL isn't set")
	}

	dsn.userName = os.Getenv("DB_USER")
	if dsn.userName == "" {
		return nil, fmt.Errorf("%s:%s", op, "DB_USER isn't set")
	}

	dsn.password = os.Getenv("DB_PASSWORD")
	if dsn.password == "" {
		return nil, fmt.Errorf("%s:%s", op, "DB_PASSWORD isn't set")
	}

	dsn.host = os.Getenv("DB_HOST")
	if dsn.host == "" {
		return nil, fmt.Errorf("%s:%s", op, "DB_HOST isn't set")
	}

	dsn.port = os.Getenv("DB_PORT")
	if dsn.port == "" {
		return nil, fmt.Errorf("%s:%s", op, "DB_PORT isn't set")
	}

	dsn.dbName = os.Getenv("DB_NAME")
	if dsn.dbName == "" {
		return nil, fmt.Errorf("%s:%s", op, "DB_NAME isn't set")
	}

	dsn.options = os.Getenv("DB_OPTIONS")
	if dsn.options == "" {
		return nil, fmt.Errorf("%s:%s", op, "DB_OPTIONS isn't set")
	}

	return dsn, nil
}
