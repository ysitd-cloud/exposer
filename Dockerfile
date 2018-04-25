FROM ysitd/dep AS builder

WORKDIR /go/src/code.ysitd.cloud/component/exposer

COPY . /go/src/code.ysitd.cloud/component/exposer

RUN dep ensure -vendor-only && \
    go build -v

FROM alpine:3.6

COPY --from=builder /go/src/code.ysitd.cloud/component/exposer/exposer /

CMD ["exposer"]