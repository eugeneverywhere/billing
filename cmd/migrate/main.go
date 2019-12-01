package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/eugeneverywhere/billing/config"
	"github.com/eugeneverywhere/billing/db"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/lillilli/vconf"
)

const (
	errorNoMigrationsChange = "no change"
)

var (
	configFile  = flag.String("config", "", "set service config file")
	down        = flag.Bool("drop", false, "set drop all migrations")
	steps       = flag.Int("steps", 0, "steps number for migrations")
	migratePath = flag.String("migrate-path", "", "path to migrations")
)

func main() {
	flag.Parse()
	cfg := &config.Config{}

	if err := vconf.InitFromFile(*configFile, cfg); err != nil {
		fmt.Printf("unable to load config: %s\n", err)
		os.Exit(1)
	}

	m, err := createMigrations(cfg.DB)
	if err != nil {
		fmt.Printf("creating migrations failed: %s\n", err)
		os.Exit(1)
	}

	if err := runMigrations(m, *steps); err != nil && err.Error() != errorNoMigrationsChange {
		fmt.Printf("migrations failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Migrations complete")
}

func createMigrations(conf config.DBConfig) (*migrate.Migrate, error) {
	database, err := sql.Open(db.MySQLDBDriver, db.GenerateMySQLDatabaseURL(conf))
	if err != nil {
		return nil, err
	}

	driver, err := mysql.WithInstance(database, &mysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", *migratePath),
		db.MySQLDBDriver,
		driver,
	)

	return m, err
}

func runMigrations(m *migrate.Migrate, steps int) (err error) {
	if steps != 0 {
		if err = m.Steps(steps); err != nil {
			return err
		}
	}

	if *down {
		err = m.Down()
	} else {
		err = m.Up()
	}

	if err != nil {
		return err
	}

	sourceErr, err := m.Close()
	if sourceErr != nil {
		return sourceErr
	}

	return err
}
