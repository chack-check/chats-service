FROM golang:1.26-alpine AS builder

RUN addgroup -g 1001 -S app && adduser -u 1001 -S app -G app

USER app

WORKDIR /home/app/

COPY go.mod go.sum /home/app/

RUN go mod download

COPY internal/ ./internal/

COPY configs/ ./configs/

COPY cmd/ ./cmd/

RUN go build -o api ./cmd/api/

RUN go build -o grpcserver ./cmd/grpcserver/

RUN go build -o consumer ./cmd/consumer/


FROM alpine AS api

RUN addgroup -g 1001 -S app && adduser -u 1001 -S app -G app

USER app

WORKDIR /home/app/

COPY --from=builder /home/app/api /home/app/api

ENTRYPOINT ["/home/app/api"]


FROM alpine AS grpcserver

RUN addgroup -g 1001 -S app && adduser -u 1001 -S app -G app

USER app

WORKDIR /home/app/

COPY --from=builder /home/app/grpcserver /home/app/grpcserver

ENTRYPOINT ["/home/app/grpcserver"]


FROM alpine AS consumer

RUN addgroup -g 1001 -S app && adduser -u 1001 -S app -G app

USER app

WORKDIR /home/app/

COPY --from=builder /home/app/consumer /home/app/consumer

ENTRYPOINT ["/home/app/consumer"]
