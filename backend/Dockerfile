# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /daily-seed-api .

# Runtime stage
FROM gcr.io/distroless/static-debian12

COPY --from=builder /daily-seed-api /daily-seed-api

EXPOSE 8080

ENTRYPOINT ["/daily-seed-api"]
