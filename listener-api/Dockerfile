# Build phase
FROM golang:alpine AS builder

LABEL maintainer="Pedro Prado"

WORKDIR /app

COPY . .

RUN go build -o listener .

# Copy phase
FROM alpine

WORKDIR /app

COPY --from=builder /app/listener .

EXPOSE 4000

CMD ["./listener"]