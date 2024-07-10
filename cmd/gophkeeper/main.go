package main

import "fmt"

var buildVersion, buildDate, buildCommit = "N/A", "N/A", "N/A"

func main() {
	// go run -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" cmd/gophkeeper/main.go
	// go build -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" cmd/gophkeeper/main.go
	// start bin file -> ./main
	fmt.Printf("Build version:%s\nBuild date:%s\nBuild commit:%s\n", buildVersion, buildDate, buildCommit)

}
