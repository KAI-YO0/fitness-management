FROM golang:1.21-alpine

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app

# Add dumb-init for support prefrok mode
RUN apk add dumb-init

# Opening ports
EXPOSE 8000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["app"]
