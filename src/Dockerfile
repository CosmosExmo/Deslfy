FROM golang:1.20.5-alpine3.18 AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o main main.go
RUN apk add --update curl && rm -rf /var/cache/apk/*
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.1/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.18.0
WORKDIR /app
COPY --from=builder /go/src/app/main /app
COPY --from=builder /go/src/app/migrate ./migrate
COPY start.sh .
COPY app.env .
COPY wait-for.sh .
COPY db/migration ./migration
EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]