FROM golang:1.21.12-alpine3.20 AS builder
RUN go env -w GO111MODULE=on
WORKDIR /sast-integrator
COPY ./    ./
RUN CGO_ENABLED=0 GOOS=linux go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.20 as golang-app
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash
RUN apk add --no-cache git
RUN apk add --no-cache curl
RUN apk add --no-cache python3 py3-pip && \
    pip3 install semgrep
WORKDIR /root/
COPY --from=builder /sast-integrator ./
#COPY --from=builder /sca-integrator/_public_key.pem ./
RUN mkdir "_scanned-project-files"
RUN mkdir "_project-repository"
CMD ["./main"]
