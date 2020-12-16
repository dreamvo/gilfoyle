LABEL maintainer="Sundowndev" \
  org.label-schema.name="Gilfoyle" \
  org.label-schema.description="Cloud-native solution to embed media streaming in any application at any scale." \
  org.label-schema.url="https://github.com/dreamvo/gilfoyle" \
  org.label-schema.vcs-url="https://github.com/dreamvo/gilfoyle" \
  org.label-schema.vendor="Dreamvo" \
  org.label-schema.schema-version="0.1"

FROM node:12.20-alpine

WORKDIR /usr/src/app
COPY ./dashboard/ui .
RUN yarn install
RUN yarn build

FROM golang:1.15.6-alpine

WORKDIR /app
RUN apk add build-base
COPY . .
COPY --from=0 /usr/src/app/dist ./dashboard/ui/dist
RUN apk add git ffmpeg
RUN go get -v -t -d ./...
RUN go generate ./...
RUN go build -v -ldflags="-s -w -X 'github.com/dreamvo/gilfoyle/config.Version=$(git describe --abbrev=0 --tags)' -X 'github.com/dreamvo/gilfoyle/config.Commit=$(git rev-parse --short HEAD)'" -v -o gilfoyle .

FROM scratch
COPY --from=1 /app/gilfoyle /bin/gilfoyle
ENTRYPOINT ["gilfoyle"]
