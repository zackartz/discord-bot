FROM andersfylling/disgord:latest as builder
MAINTAINER https://github.com/zackartz
WORKDIR /build
COPY . /build
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags \"-static\"' -o discordbot .

FROM debian
RUN apt-get update && apt-get install ffmpeg
WORKDIR /bot
COPY --from=builder /build/discordbot .
CMD ["/bot/discordbot"]
