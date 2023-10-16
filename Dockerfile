FROM golang:1.21-alpine as builder
RUN apk update
RUN apk add --no-cache build-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o cosmobot

FROM alpine:3 as prod
# copy app
WORKDIR /app
COPY --from=builder /app/cosmobot .

EXPOSE 8080

CMD ["/bin/sh", "-c", "./cosmobot"]