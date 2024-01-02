FROM golang:1.18.6 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app internal/app.go

CMD ["app"]