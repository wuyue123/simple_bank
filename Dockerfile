# build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=off
ENV GO111MODULE=on

RUN go build -o main main.go

# run stage
FROM alpine:3.18 AS final
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
EXPOSE 8080
CMD ["./main"]