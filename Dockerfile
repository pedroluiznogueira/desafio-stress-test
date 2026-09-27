FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o /out/stress-test .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /out/stress-test /stress-test
ENTRYPOINT ["/stress-test"]
