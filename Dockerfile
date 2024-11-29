FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.sum ./
COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o goapp

FROM scratch

COPY --from=builder /app/goapp /bin/goapp

ENTRYPOINT [ "/bin/goapp" ]