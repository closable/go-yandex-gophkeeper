# go-yandex-gophkeeper
## _Менеджер приватной информации_
[выпускная работа practicum.yandex.ru](https://practicum.yandex.ru/profile/go-advanced/)

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://github.com/closable/go-yandex-gophkeeper)

go-yandex-gophkeeper это современное и производительное клиент-серверное решение для работы с конфидициальной иноформацией написанной на крутячем языке Go lang ✨Magic ✨

## Features

- Сохранение(шифрование) любой приватной текстовой информации
- Просмотр содержимого в незашифрованном виде
- Возможность удаления(изменения информации)
- Сохранение файлов с последующим восстановлением
- Возможность сохранения(шифрования) указанных директорий
- Все файлы(директории) размером более 100 МБ будут сохраняться в Minio ObjectStore
  (Warning! bucket и KeyAccess необходимо задать самостоятельной)  
- Наличие готовых клиентов под разные платформы
- Удобный полнофункциональный CLI интерфейс (для сервера и клиента)
- Конфигурирование
- Офлайн режим (возможность только просмотра приватной информации)

## Tech

- [Go lang](https://go.dev/) An open-source programming language supported by Google
- [gRPC](https://grpc.io/) A high performance, open source universal RPC framework
- [bubbletea](https://github.com/charmbracelet/bubbletea) Bubble Tea is well-suited for simple and complex terminal applications, either inline, full-window, or a mix of both.
- [postgresql](https://www.postgresql.org/) PostgreSQL: The World's Most Advanced Open Source Relational Database
- [minio](https://min.io/) MinIO is a high-performance, S3 compatible object store.

### Парметры командной строки

- a=":3000" Сервер
- f=":3100" Сервер для хранения файлов < 100 MB
- m=":9000" Minio cервер для хранения файлов > 100 MB
- u="user" Логин для входа на серевере
- p="password" Парол для входа на серевере
- c  Внутренний интерфейс для работы в режиме CLI (по умолчанию использован TUI)

For server...

```sh
docker build -t go-yandex-gophkeeper:local .
```

```sh
docker-compose up    
```

For client
```sh
.client-darwin-m1 -a=:3000 -f=:3100 m=:9000
```

### Внимание! для работы с minio необходимо параметры доступа задать самостоятельно 
```
example .env file in the root folder app
  S3_BUCKET=gophkeeper
  S3_ACCESS_KEY=DNwRXfu3SAqGxZBtqLTi
  S3_SECRET_KEY=saC76TvvR5P7calPgkhvMfxO3HE68OtfaFYt1HYb
```

## License

MIT

**Free Software, Hell Yeah!**
