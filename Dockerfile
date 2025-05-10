FROM golang:1.24.2-bookworm
WORKDIR /app

COPY cmd/ ./cmd/
COPY go.mod .
COPY go.sum .
COPY public/ ./public/
COPY templates/ ./templates/
RUN mkdir -p ./bin
RUN wget https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -O ./bin/tailwind
RUN chmod +x ./bin/tailwind
RUN /app/bin/tailwind -i ./templates/main.css -o ./dist/output.css --minify

ENV PORT=8080
EXPOSE ${PORT}

RUN go build -o /app/bin/server ./cmd/server/main.go
CMD [ "/app/bin/server" ]