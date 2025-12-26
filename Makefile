dist/lcsv: main.go
	mkdir -p dist/
	go build -o $@ $<

.PHONY: run
run: dist/lcsv
	@$<

.PHONY: all
all: dist/lcsv

clean:
	rm -rf ./dist/ 2>/dev/null

help:
	@echo "Build the executable from source"
	@echo "  make "
	@echo "  make lcsv"
	@echo "  make build"

	@echo "Build and run the project"
	@echo "  make run"

	@echo "Clean artifacts"
	@echo "  make clean"

	@echo "Display this help message"
	@echo "  make help"
