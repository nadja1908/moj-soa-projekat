FROM golang:1.25.4-alpine AS builder

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

WORKDIR /app

# Copy source code first
COPY . .

# Initialize module and download dependencies
RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8000


CMD ["./main"]