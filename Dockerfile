FROM golang:1.21-alpine as builder
RUN apk update
RUN apk add --no-cache build-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o cosmoBot

FROM gcr.io/distroless/base-debian12
WORKDIR /app
# copy app
COPY --from=builder /app/cosmoBot /app/cosmoBot

EXPOSE 8080
ENTRYPOINT [ "/app/cosmoBot" ]