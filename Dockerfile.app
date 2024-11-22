FROM golang:latest

RUN mkdir /app
WORKDIR /app
ENV CONFIG_PATH=./config/remote.yaml
ENV POSTGRES_PASS=12345pass
ENV REDIS_PASS=admin

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cmd/sso ./cmd/sso/main.go

CMD ["/app/cmd/sso/main"]