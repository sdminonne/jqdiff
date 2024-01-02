

jqdiff:
	go build .

PONY: test
test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./...

view-cover:
	go tool cover -html=coverage.out
