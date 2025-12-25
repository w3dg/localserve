.PHONY: run

lcsv:
	mkdir -p dist/
	go build -o dist/lcsv ./cmd/main.go

run: lcsv
	./dist/lcsv
