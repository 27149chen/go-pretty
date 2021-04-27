FROM golang:1.15 as build

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ENV GOPROXY=https://goproxy.cn,direct

COPY third_party third_party
COPY go.mod go.sum ./

RUN go mod download

# === Code Reserved for enterprise only. START ===
COPY cmd cmd
COPY pkg pkg
# === Code Reserved for enterprise only. END ===

RUN go build -v -o /pretty ./cmd/pretty/main.go

FROM ubuntu:16.04

WORKDIR /app

COPY --from=build /pretty .

ENTRYPOINT ["/app/pretty"]
