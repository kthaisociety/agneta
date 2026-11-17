FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /agneta ./cmd/agneta

# Alpine with CA certificates and a non-root user, like lucid.
FROM alpine:latest
RUN apk --no-cache add ca-certificates \
    && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /agneta .
ENV HOST=0.0.0.0 \
    PORT=8000
USER app
EXPOSE 8000
CMD ["./agneta"]
