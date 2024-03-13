FROM golang:1.22.1-alpine3.18 as builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o data-manager cmd/data-management.go

CMD ["go", "run", "cmd/data-management.go", "run"]

FROM alpine as swagger

WORKDIR /app

RUN apk add unzip wget curl

RUN wget https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.11.10.zip

RUN unzip v5.11.10.zip

RUN sed

FROM alpine

COPY --from=builder /app/data-manager /data-manager

COPY --from=swagger /app/swagger-ui-5.11.10/dist /app/

CMD ["/data-manager", "run"]