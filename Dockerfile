FROM golang:1.20.4-bullseye

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN go get collector

RUN go build -o /collector

CMD ["/collector"]