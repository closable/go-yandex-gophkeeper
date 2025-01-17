package config

import "flag"

// переменные
var (
	// Адрес сервера
	FlagRunAddr string
	// Адрес сервера
	FlagFileRunAddr string
	// Адрес minij server
	FlagMinioAddr string
	// Использование СУБД
	FlagDSN string
	// Логин
	FlagLogin string
	// Пароль
	FlagPassword string
	// Inner CLI interface
	FlagCLI bool
)

// config описание структур данных среды окружения
type config struct {
	ServerAddress      string
	FileServerAddress  string
	MinioServerAddress string
	DSN                string
	Login              string
	Password           string
	CLI                bool
}

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную flagRunAddr
	flag.StringVar(&FlagRunAddr, "a", ":3000", "address and port to run server")
	flag.StringVar(&FlagFileRunAddr, "f", ":3100", "address and port to run file server")
	flag.StringVar(&FlagMinioAddr, "m", "localhost:9000", "address and port minio server")
	//flag.StringVar(&FlagDSN, "d", "postgres://postgres:1303@localhost:5432/postgres", "access to DBMS")
	flag.StringVar(&FlagDSN, "d", "postgres://postgres:postgres@host.docker.internal:25432/postgres?sslmode=disable", "access to DBMS")
	// Для работы в режиме CLI
	flag.StringVar(&FlagLogin, "u", "", "login access to app")
	// Для работы в режиме CLI
	flag.StringVar(&FlagPassword, "p", "", "password access to app")
	flag.BoolVar(&FlagCLI, "c", false, "use innser CLI interfase (default TUI)")

	flag.Parse()
}

// LoadConfig загружаем данные среды окружения
func LoadConfig() *config {
	ParseFlags()

	var config = &config{
		ServerAddress:      FlagRunAddr,
		FileServerAddress:  FlagFileRunAddr,
		MinioServerAddress: FlagMinioAddr,
		DSN:                FlagDSN,
		Login:              FlagLogin,
		Password:           FlagPassword,
		CLI:                FlagCLI,
	}

	return config
}
