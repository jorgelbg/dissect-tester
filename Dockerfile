FROM golang:latest as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM alpine
WORKDIR /app
COPY --from=builder /build/main /app
COPY --from=builder /build/static static
COPY --from=builder /build/templates templates
EXPOSE 8080
CMD ["./main"]
