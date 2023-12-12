css:
	@npx tailwindcss -i ./templates/main.css -o ./dist/output.css

build: css
	@go build -o ./bin/server ./cmd/server/main.go

run: build
	@./bin/server