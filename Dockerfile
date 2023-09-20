FROM golang:1.21-alpine

RUN apk add --no-cache tzdata
ENV TZ=America/Denver

WORKDIR /usr/src/jime

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY jime.go config_jime.json log_jime.txt README.md .
RUN go build -v -o /usr/local/bin/jime ./...

CMD ["jime"]
