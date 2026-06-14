FROM golang:1.26.3-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/freelib/exe /app/cmd/freelib

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/cmd/freelib/exe /app
CMD ["/app/exe"]