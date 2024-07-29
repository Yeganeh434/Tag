FROM golang:1.22.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /tag_app ./cmd/main.go

RUN apt-get update && apt-get install -y curl && \
    curl -o /wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && \
    chmod +x /wait-for-it.sh

EXPOSE 8080

CMD ["/wait-for-it.sh", "db:3306", "--", "/tag_app"]
