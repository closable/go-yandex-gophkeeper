proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/services/proto/gophkeeper.proto

build-win:
	GOOS=windows GOARCH=amd64 go build -o cmd/gophkeeper/client/bin/client-windows-x86-64.exe -ldflags "\
	-X 'github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version.buildVersion=v1.0.0' \
	-X 'github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version.gitCommit=$(git rev-parse HEAD)' \
	-X 'github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')' " \
	cmd/gophkeeper/client/main.go

build-m1:
	GOOS=darwin GOARCH=arm64 go build -o cmd/gophkeeper/client/bin/client-darwin-m1 -ldflags "\
	-X 'github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version.buildVersion=v1.0.0' \
	-X 'github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version.gitCommit=$(git rev-parse HEAD)' \
	-X 'github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')' " \
	cmd/gophkeeper/client/main.go

build-linux:
	GOOS=linux GOARCH=386 go build -o cmd/gophkeeper/client/bin/client-linux-386 -ldflags "\
	-X 'github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version.buildVersion=v1.0.0' \
	-X 'github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version.gitCommit=$(git rev-parse HEAD)' \
	-X 'github.com/closable/go-yandex-gophkeeper/cmd/gophkeeper/version.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')' " \
	cmd/gophkeeper/client/main.go

test:
	go test ./...

cover:
	go test -cover ./...

.PHONY: proto build-win build-m1 build-linux test cover