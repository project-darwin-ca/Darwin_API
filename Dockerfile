FROM golang:1.22.1-alpine3.18 as builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o data-manager cmd/data-management.go

FROM alpine

COPY --from=builder data-manager /data-manager

CMD ["/data-manager", "run"]