##########################################
# Builder image to build the application #
##########################################
FROM golang:1.23-alpine as builder

ARG SRC_PATH

RUN apk add --no-cache upx=4.2.4-r0

WORKDIR /go/src/github.com/thoughtgears/pam-manager

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/app-linux-amd64 . && \
    upx --best --lzma builds/app-linux-amd64


##########################################
# Final image to run the application     #
##########################################
FROM alpine:3.20

ARG SRC_PATH

RUN apk --no-cache add ca-certificates=20240705-r0

WORKDIR /app/
COPY --from=builder /go/src/github.com/thoughtgears/pam-manager/builds/app-linux-amd64 ./app

EXPOSE 8080

CMD ["./app"]
