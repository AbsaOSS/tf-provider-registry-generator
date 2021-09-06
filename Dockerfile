# build stage
FROM golang:1.16 as builder
WORKDIR /workspace
ARG GOARCH=amd64

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY internal/ internal/ 
COPY main.go main.go
RUN go build -o tfreg-golang

FROM alpine

RUN apk --no-cache add git bash coreutils git-lfs gnupg
RUN mkdir /data
COPY data /data
COPY --from=builder /workspace/tfreg-golang /usr/bin/tfreg-golang
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /usr/bin/tfreg-golang

ENTRYPOINT ["/entrypoint.sh"]
