FROM golang:1.9.3 as builder
WORKDIR /go/src/github.com/bdimcheff/rainbow-deploys/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make build

FROM scratch
COPY --from=builder /go/src/github.com/bdimcheff/rainbow-deploys/rainbow-deploys .
ARG COLOR
ENV COLOR ${COLOR}
ENTRYPOINT ["/rainbow-deploys"]