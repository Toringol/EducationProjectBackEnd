FROM golang:1.15

ENV GO111MODULE=on

WORKDIR /app

COPY . .

EXPOSE 8080
EXPOSE 8081

CMD ["go", "run", "./cmd/authService/authService.go"]