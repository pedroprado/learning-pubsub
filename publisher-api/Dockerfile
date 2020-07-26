FROM golang:alpine AS builder
# Build phase
LABEL maintainer="Pedro Prado"

WORKDIR /app

COPY . .

RUN go build -o publisher .


# Copy phase
FROM alpine

EXPOSE 5000

WORKDIR /app

COPY --from=builder /app/publisher .

CMD ["./publisher"]