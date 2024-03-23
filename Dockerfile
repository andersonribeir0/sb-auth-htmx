FROM golang:latest as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o migrate cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o dreampicai cmd/api/main.go

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/migrate .
COPY --from=builder /app/dreampicai .

COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]

