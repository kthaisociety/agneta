FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /agneta ./cmd/agneta

FROM scratch
COPY --from=build /agneta /agneta
USER 10001
EXPOSE 8000
ENTRYPOINT ["/agneta"]
