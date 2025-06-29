FROM golang:1.23.10

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod tidy

COPY . .

RUN go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 2345

# 注意這裡不要 build binary，留給 delve debug
CMD ["dlv", "debug", "--headless", "--listen=:2345", "--api-version=2", "--accept-multiclient"]
