package database

import "fmt"

type Config struct {
	Hostname string `env:"DB_HOSTNAME"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Port     uint32 `env:"DB_PORT" envDefault:"3306"`
	Name     string `env:"DB_NAME"`
}

func (databaseConfig Config) DSN() string {
	return databaseConfig.Username +
		":" +
		databaseConfig.Password +
		"@tcp(" +
		databaseConfig.Hostname +
		":" +
		fmt.Sprintf("%d", databaseConfig.Port) +
		")/" +
		databaseConfig.Name + "?parseTime=true"
}
