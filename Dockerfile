FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o ./server ./cmd/server/main.go

EXPOSE 8080
CMD [ "./server" ]