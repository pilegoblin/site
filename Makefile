css:
	@./bin/tailwind -i ./templates/main.css -o ./dist/output.css --minify

build: css
	@go build -o ./bin/server ./cmd/server/main.go

run: build
	@./bin/server

update-tailwind:
	curl -sL0 https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 > ./bin/tailwind
	chmod +x ./bin/tailwind

docker-build:
	docker build -t pilegoblin/website .

docker-run: docker-build
	docker run -p 8080:8080 pilegoblin/website