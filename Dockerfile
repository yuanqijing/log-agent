FROM golang:1.18-alpine as builder

WORKDIR /go/src/github.com/yuanqijing/log-agent
RUN apk add --update make git bash rsync gcc musl-dev

# Copy the go source
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

# Copy the go source
COPY apis/ apis/
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build
ENV GOOS linux
ENV GOARCH amd64
RUN go build -a -o log-agent /go/src/github.com/yuanqijing/log-agent/cmd/log-agent/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to  https://github.com/GoogleContainerTools/distroless for more details
FROM alpine:3.12
WORKDIR /
COPY --from=builder /go/src/github.com/yuanqijing/log-agent/log-agent .

USER root:root