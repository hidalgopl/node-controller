FROM golang:1.12-alpine3.9 as builder

WORKDIR /controller
COPY . .

RUN apk add --no-cache \
    gcc \
    git \
    linux-headers \
    make \
    musl-dev

RUN make build

FROM alpine:3.9

COPY --from=builder /controller/bin/ /usr/local/bin/

RUN addgroup -g 1000 runnergroup && \
    adduser -h / -D -u 1000 -G runnergroup runner
USER runner

ENTRYPOINT ["node-controller"]
