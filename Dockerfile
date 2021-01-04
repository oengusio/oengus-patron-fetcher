FROM golang:1.15.6-alpine AS builder

WORKDIR /oengus-patreon
#COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /oengus-patreon
COPY --from=builder /oengus-patreon/main ./main
RUN chmod +x main

CMD ["./main"]
