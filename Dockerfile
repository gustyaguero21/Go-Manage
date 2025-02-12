FROM golang:1.23

WORKDIR /app

COPY . .

RUN ls -R /app

RUN go mod tidy

RUN go build -v -o mi-app ./cmd/api

EXPOSE 8080

CMD ["/app/mi-app"]
