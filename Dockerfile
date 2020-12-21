FROM golang:1.15.6

WORKDIR /oengus-patreon
#COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main .

CMD ["./main"]
