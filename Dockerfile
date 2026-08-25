FROM golang:1.26-alpine AS builder

RUN addgroup -g 1001 -S app && adduser -u 1001 -S app -G app

USER app

WORKDIR /home/app/

COPY src/go.mod src/go.sum /home/app/

RUN go mod download

COPY src/ /home/app/

RUN go build -o server .


FROM alpine

RUN addgroup -g 1001 -S app && adduser -u 1001 -S app -G app

USER app

WORKDIR /home/app/

COPY --from=builder /home/app/server /home/app/server

ENTRYPOINT ["/home/app/server"]
