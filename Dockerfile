# Build stage
FROM golang:bullseye AS base
RUN apt-get update && apt-get install -y ffmpeg
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-yt

# Final stage
FROM alpine:latest 
RUN apk --no-cache add ffmpeg

COPY --from=base /docker-yt /docker-yt
EXPOSE 5000

CMD ["/docker-yt"]
