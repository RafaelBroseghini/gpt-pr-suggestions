FROM golang:1.20

WORKDIR /
COPY . /
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o gpt-pr-suggestions main.go

ENTRYPOINT ["/entrypoint.sh"]