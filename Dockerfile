FROM golang:bullseye

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-yt

EXPOSE 3000

CMD ["/docker-yt"]
