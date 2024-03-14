FROM golang:1.22.1-alpine3.18 as builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o data-manager cmd/data-management.go

CMD ["go", "run", "cmd/data-management.go", "run"]

FROM swaggerapi/swagger-ui as swagger-builder

COPY dev/swagger.json /usr/share/nginx/html/

FROM alpine

COPY --from=builder /app/data-manager /data-manager

COPY --from=swagger-builder /usr/share/nginx/html /opt/swagger

COPY entrypoint.sh .

ENTRYPOINT ["/entrypoint.sh"]
