FROM golang:latest as builder
ARG MOD
ENV MOD ${MOD:-readonly}
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN echo "go mod flag: $MOD"
RUN CGO_ENABLED=0 GOOS=linux go build -mod=$MOD -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM alpine
WORKDIR /app
COPY --from=builder /build/main /app
COPY --from=builder /build/static static
COPY --from=builder /build/templates templates
EXPOSE 8080
CMD ["./main"]
