package config

import "flag"

// переменные
var (
	// Адрес сервера
	FlagRunAddr string
	// Использование СУБД
	FlagDSN string
	// Логин
	FlagLogin string
	// Пароль
	FlagPassword string
)

// config описание структур данных среды окружения
type config struct {
	ServerAddress string
	DSN           string
	Login         string
	Password      string
}

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagDSN, "d", "postgres://postgres:1303@localhost:5432/postgres", "access to DBMS")
	// Для работы в режиме CLI
	flag.StringVar(&FlagLogin, "u", "", "login access to app")
	// Для работы в режиме CLI
	flag.StringVar(&FlagPassword, "p", "", "password access to app")

	flag.Parse()
}

// LoadConfig загружаем данные среды окружения
func LoadConfig() *config {
	ParseFlags()

	var config = &config{
		ServerAddress: FlagRunAddr,
		DSN:           FlagDSN,
		Login:         FlagLogin,
		Password:      FlagPassword,
	}

	return config
}
