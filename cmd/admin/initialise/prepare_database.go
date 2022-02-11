package initialise

import (
	"database/sql"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/database"
)

func prepareDB(config database.Config, user, password, sslCert, sslKey string) error {
	adminConfig := config
	adminConfig.User = user
	adminConfig.Password = password
	adminConfig.SSL.Cert = sslCert
	adminConfig.SSL.Key = sslKey

	db, err := database.Connect(adminConfig)
	if err != nil {
		return err
	}

	logging.Info("verify user")
	if err = verifyUser(db, config); err != nil {
		return err
	}
	logging.Info("verify database")
	if err = verifyDB(db, config); err != nil {
		return err
	}
	logging.Info("verify grant")
	if err = verifyGrant(db, config); err != nil {
		return err
	}

	return db.Close()
}

func verifyUser(db *sql.DB, config database.Config) error {
	exists, err := existsUser(db, config)
	if exists || err != nil {
		return err
	}
	return createUser(db, config)
}

func existsUser(db *sql.DB, config database.Config) (exists bool, err error) {
	row := db.QueryRow("SELECT EXISTS(SELECT username FROM [show roles] WHERE username = $1)", config.User)
	err = row.Scan(&exists)
	return exists, err
}

func createUser(db *sql.DB, config database.Config) error {
	_, err := db.Exec("CREATE USER $1 WITH PASSWORD $2", config.User, &sql.NullString{String: config.Password, Valid: config.Password != ""})
	return err
}

func verifyDB(db *sql.DB, config database.Config) error {
	exists, err := existsDatabase(db, config)
	if exists || err != nil {
		return err
	}
	return createDatabase(db, config)
}

func existsDatabase(db *sql.DB, config database.Config) (exists bool, err error) {
	row := db.QueryRow("SELECT EXISTS(SELECT database_name FROM [show databases] WHERE database_name = $1)", config.Database)
	err = row.Scan(&exists)
	return exists, err
}

func createDatabase(db *sql.DB, config database.Config) error {
	_, err := db.Exec("CREATE DATABASE " + config.Database)
	return err
}

func verifyGrant(db *sql.DB, config database.Config) error {
	exists, err := hasGrant(db, config)
	if exists || err != nil {
		return err
	}
	return grant(db, config)
}

func hasGrant(db *sql.DB, config database.Config) (has bool, err error) {
	row := db.QueryRow("SELECT EXISTS(SELECT * FROM [SHOW GRANTS ON DATABASE "+config.Database+"] where grantee = $1 AND privilege_type = 'ALL')", config.User)
	err = row.Scan(&has)
	return has, err
}

func grant(db *sql.DB, config database.Config) error {
	_, err := db.Exec("GRANT ALL ON DATABASE " + config.Database + " TO " + config.User)
	return err
}
