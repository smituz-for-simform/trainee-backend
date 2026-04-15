# # 🔹 Stage 1: Build
# FROM golang:1.25-alpine AS builder

# WORKDIR /app

# # Copy go mod first (cache optimization)
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy rest of code
# COPY . .

# # Build binary
# RUN go build -o main .

# # 🔹 Stage 2: Run
# FROM alpine:latest

# WORKDIR /app

# # Copy binary from builder
# COPY --from=builder /app/main .

# # Create uploads directory
# RUN mkdir -p uploads

# # Expose port
# EXPOSE 8080

# # Run app
# CMD ["./main"]

# 🔹 Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go clean -cache && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 🔹 Stage 2: Run
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p uploads

#for healthcheck and ready db check
RUN apk add --no-cache wget

RUN adduser -D appuser \
    && mkdir -p uploads \
    && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./main"]