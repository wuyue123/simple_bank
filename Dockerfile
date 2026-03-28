# build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=off
ENV GO111MODULE=on
#RUN apk add --no-cache curl
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-arm64.tar.gz
RUN go build -o main main.go

# run stage
FROM alpine:3.18 AS final
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-arm64.tar.gz ./migrate
COPY app.env .
COPY db/migration ./migration
COPY sqlc.yaml .
COPY start.sh .

EXPOSE 8080
CMD ["./main"]
ENTRYPOINT [ "/app/start.sh" ]