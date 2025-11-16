FROM alpine:3.22

WORKDIR /migrations

RUN apk update && \
    apk add --no-cache \
    git \
    go \
    make \
    postgresql-client

ENV PATH="/root/go/bin:${PATH}"
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

ENV GOOSE_DRIVER=postgres

CMD ["goose", "up"]
