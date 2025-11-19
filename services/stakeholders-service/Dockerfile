# Stage 1: Build faza
FROM golang:1.23-alpine AS builder

# Radni direktorijum unutar kontejnera
WORKDIR /app

# Copy source code first
COPY . .

# Initialize module and download dependencies
RUN go mod tidy
RUN go mod download

# Build-uj stakeholders-service
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main

# Stage 2: Production faza
FROM alpine:latest

# Instaliraj wget za health check
RUN apk --no-cache add ca-certificates wget

WORKDIR /root/
COPY --from=builder /app/main .

COPY assets/images /root/images

EXPOSE 8001
CMD ["./main"]