FROM golang:1.21-alpine AS builder

RUN apk --no-cache add bash git make gcc gettext musl-dev

WORKDIR /usr/local/src

COPY ["./go.mod", "./go.sum", "./"]
RUN go mod download 

COPY ["cmd", "./cmd"]
COPY ["internal", "./internal"]

RUN go build -o ./bin/gophkeeper ./cmd/gophkeeper/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/gophkeeper /

CMD ["/gophkeeper"]
