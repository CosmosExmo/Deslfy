FROM golang:1.20.5-alpine3.18 AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o main main.go

FROM alpine:3.18.0
WORKDIR /app
COPY --from=builder /go/src/app/main /app
COPY start.sh .
COPY app.env .
COPY wait-for.sh .
COPY db/migration ./db/migration
EXPOSE 8080
EXPOSE 9090
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]