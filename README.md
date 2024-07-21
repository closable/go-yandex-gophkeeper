# go-yandex-gophkeeper
## _Менеджер приватной информации (выпускная работа)_

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://github.com/closable/go-yandex-gophkeeper)

go-yandex-gophkeeper это современное и производительное клиент-серверное решение для работы с конфидициальной иноформацией написанной на крутячем языке Go lang ✨Magic ✨

## Features

- Сохранение(шифрование) любой приватной текстовой информации
- Просмотр содержимого в незашифрованном виде
- Возможность удаления(изменения информации)
- Сохранение файлов с последующим восстановлением
- Возможность сохранения(шифрования) указанных директорий
- Наличие готовых клиентов под разные платформы
- Удобный полнофункциональный CLI интерфейс (для сервера и клиента)
- Конфигурирование

## Tech

- [Go lang](https://go.dev/) An open-source programming language supported by Google
- [gRPC](https://grpc.io/) A high performance, open source universal RPC framework
- [postgresql](https://www.postgresql.org/) PostgreSQL: The World's Most Advanced Open Source Relational Database

For server...

```sh
./server -a="192.168.0.116:8080" -f="192.168.0.116:8090"
```

For client
```sh
./client-darwin-m1 -a="192.168.0.116:8080" -f="192.168.0.116:8090"
```


## License

MIT

**Free Software, Hell Yeah!**
