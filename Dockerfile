ARG GOLANG_VERSION=1.15
FROM golang:${GOLANG_VERSION} as builder

WORKDIR /azure-quota-provider

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN cd /azure-quota-provider && go mod download

# Build
RUN go build -a -o server ./cmd/server

# Use distroless as minimal base image to package the manager binary
FROM gcr.io/distroless/base:latest-amd64
WORKDIR /

LABEL org.opencontainers.image.source https://github.com/ydataai/azure-quota-provider

COPY --from=builder /azure-quota-provider/server .

ENTRYPOINT ["/server"]
