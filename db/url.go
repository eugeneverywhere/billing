package db

import (
	"fmt"
	"github.com/eugeneverywhere/billing/config"
)

// MySQLDBDriver - driver name for mysql db
const MySQLDBDriver = "mysql"

// GenerateMySQLDatabaseURL - generate path to mysql database
func GenerateMySQLDatabaseURL(conf config.DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?multiStatements=true&parseTime=true",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name,
	)
}
