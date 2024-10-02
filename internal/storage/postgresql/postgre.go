package postgresql

import (
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/nabishec/restapi/internal/config"
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

func NewDatabase(cfg config.DBDataSourceName) (*Database, error) {
	var database Database
	dbConnfig := &dataSourceName{
		protocol: cfg.Protocol,
		userName: cfg.UserName,
		password: cfg.Password,
		host:     cfg.Host,
		port:     cfg.Port,
		dbName:   cfg.DBName,
		options:  cfg.Options,
	}
	err := database.connectDatabase(dbConnfig)
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
