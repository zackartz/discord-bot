FROM golang:1.15 AS build

WORKDIR /go/src/db

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

FROM alpine:latest

WORKDIR /app

RUN apk add --update ffmpeg

COPY --from=build /go/src/db/app .

EXPOSE 1337

CMD ["./app", "-prod"]
