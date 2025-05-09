css:
	@npx @tailwindcss/cli -i ./templates/main.css -o ./dist/output.css --minify

build: css
	@go build -o ./bin/server ./cmd/server/main.go

run: build
	@./bin/server