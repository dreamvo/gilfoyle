FROM golang:1.14-alpine as builder

LABEL maintainer="Sundowndev" \
  org.label-schema.name="Gilfoyle" \
  org.label-schema.description="Video streaming server backed by decentralized filesystem." \
  org.label-schema.url="https://github.com/dreamvo/gilfoyle" \
  org.label-schema.vcs-url="https://github.com/dreamvo/gilfoyle" \
  org.label-schema.vendor="Dreamvo" \
  org.label-schema.schema-version="0.1"

WORKDIR /app

COPY . .

RUN go get -v -t -d ./...

RUN go generate ./...

RUN go build -v -o gilfoyle .

FROM scratch

WORKDIR /app

COPY --from=builder /app/gilfoyle .

ENTRYPOINT ["/app/gilfoyle"]
