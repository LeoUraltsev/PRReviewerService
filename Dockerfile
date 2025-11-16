FROM golang:1.25.4-alpine As builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/pr-reviewer-service ./cmd/main.go

FROM alpine:3.22

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/pr-reviewer-service .

EXPOSE 8080

ENV POSTGRES_PR_CONNECTION_STRING=postgres://postgres:password@postgres:5432/database

CMD ["./pr-reviewer-service"]
