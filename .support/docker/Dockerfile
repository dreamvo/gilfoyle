FROM node:16.4.2-alpine

LABEL maintainer="raphael@crvx.fr" \
  org.label-schema.name="Gilfoyle" \
  org.label-schema.description="Cloud-native solution to embed media streaming in any application at any scale." \
  org.label-schema.url="https://github.com/dreamvo/gilfoyle" \
  org.label-schema.vcs-url="https://github.com/dreamvo/gilfoyle" \
  org.label-schema.vendor="Dreamvo" \
  org.label-schema.schema-version="0.1"

WORKDIR /usr/src/app
COPY ./dashboard/ui .
RUN yarn install
RUN yarn build

FROM golang:1.16.0-alpine

WORKDIR /app

RUN set -ex && \
    apk add --no-cache ffmpeg ffmpeg-libs && \
    apk add --no-cache gcc build-base linux-headers

COPY . .
COPY --from=0 /usr/src/app/dist ./dashboard/ui/dist

RUN go get -v -t -d ./...
RUN go generate ./...
RUN GOOS=linux go build -v -ldflags="-s -w -X 'github.com/dreamvo/gilfoyle/config.Version=$(git describe --abbrev=0 --tags)' -X 'github.com/dreamvo/gilfoyle/config.Commit=$(git rev-parse --short HEAD)'" -v -o gilfoyle ./cmd/main.go

ENTRYPOINT ["/app/gilfoyle"]
