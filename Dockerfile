# build stage
FROM golang:1.24 AS build

# set working directory
WORKDIR /app

# copy source code
COPY . .

# install dependencies
RUN go mod download

# build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./gorest ./cmd/http/main.go

EXPOSE 8080

ENTRYPOINT [ "./gorest" ]