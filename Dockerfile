FROM golang:1.25

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /app/exe /app/cmd/api/main.go

CMD ["/app/exe"]
