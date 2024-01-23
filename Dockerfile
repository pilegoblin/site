FROM node:latest as css
WORKDIR /app
COPY templates/ ./templates/
COPY tailwind.config.js .
RUN npx tailwindcss -i ./templates/main.css -o ./dist/output.css --minify

FROM golang:latest
WORKDIR /app

COPY cmd/ ./cmd/
COPY go.mod .
COPY go.sum .
COPY public/ ./public/
COPY --from=css /app/dist/output.css /app/dist/output.css 
COPY --from=css /app/templates/ /app/templates/

ENV PORT=8080
EXPOSE ${PORT}

RUN go build -o /app/bin/server ./cmd/server/main.go
CMD [ "/app/bin/server" ]