FROM golang:1.21

WORKDIR /src/

COPY src/go.mod src/go.sum /src/

RUN go mod download

COPY src/ /src/

RUN go build -o server .

ENTRYPOINT [ "./server" ]