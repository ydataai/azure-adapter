ARG GOLANG_VERSION=1.20
FROM golang:${GOLANG_VERSION} as builder
ARG COMPILE_CMD

WORKDIR /workspace

COPY . .

RUN go mod download

# Build
RUN CGO_ENABLED=0 go build -a -o main "./cmd/${COMPILE_CMD}"

# Use distroless as minimal base image to package the manager binary
FROM gcr.io/distroless/base:latest-amd64
WORKDIR /

LABEL org.opencontainers.image.source https://github.com/ydataai/azure-adapter

COPY --from=builder /workspace/main .

ENTRYPOINT ["/main"]
