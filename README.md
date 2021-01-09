# gilfoyle [![godoc](https://godoc.org/github.com/dreamvo/gilfoyle?status.svg)](https://godoc.org/github.com/dreamvo/gilfoyle) [![build](https://img.shields.io/endpoint.svg?url=https://actions-badge.atrox.dev/dreamvo/gilfoyle/badge?ref=master)](https://github.com/dreamvo/gilfoyle/actions) [![goreport](https://goreportcard.com/badge/github.com/dreamvo/gilfoyle)](https://goreportcard.com/report/github.com/dreamvo/gilfoyle) [![Coverage Status](https://coveralls.io/repos/github/dreamvo/gilfoyle/badge.svg?branch=master)](https://coveralls.io/github/dreamvo/gilfoyle?branch=master) [![release](https://img.shields.io/github/release/dreamvo/gilfoyle.svg)](https://github.com/dreamvo/gilfoyle/releases)

Gilfoyle is a web application from the [Dreamvo project](https://dreamvo.com) that runs a self-hosted video streaming server. This application allows you to setup an enterprise-grade media encoding & streaming platform in minutes. Gilfoyle handles media upload, processing and streaming.

It's written in Golang, designed for [Kubernetes](http://kubernetes.io/) and runs as a single Linux binary with [PostgreSQL](https://www.postgresql.org/) and [RabbitMQ](https://www.rabbitmq.com/).

## Table of content

- [Features](#features)
- [Current status](#current-status)
- [Design](#design)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [Discussion](#discussion)

## Features

- Deploy a RESTful API and HLS server to manage and stream audio & video
- Handle video compression and encoding with [FFmpeg](https://ffmpeg.org/)
- Customize media transcoding
- Highly scalable architecture
- Monitoring: Prometheus exported metrics, embedded Web UI
- Media attachments: attach files such as subtitles or preview images to medias

### What's next ?

- Media asset generation (thumbnail, video preview...)
- More supported formats (e.g: *360° videos, 60fps*...)
- Authentication and delegated upload
- Live streaming
- [IPFS](https://ipfs.io/) support
- Encryption support

## Current status

It's a **Work In Progress**. As this project is very recent, it's under heavy development and not suitable for production yet. Please consider v0 as unstable. Want to contribute ? Check the [backlog](https://github.com/dreamvo/gilfoyle/projects/1).

## Design

See [this document](DESIGN.md) for a high level design and goals.

## Documentation

- For **contributors**: see [godoc](https://godoc.org/github.com/dreamvo/gilfoyle), [high-level design documentation](DESIGN.md)
- For **users**: see [user guide](https://dreamvo.github.io/gilfoyle/) (WIP) and [API documentation](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/dreamvo/gilfoyle/master/api/docs/swagger.json)
- The [GPL v3](LICENSE) license

## Contributing

This project is in an early stage. We appreciate feedbacks and [discussions](#discussion) about the design and [features](#features).

## Discussion

- [Discuss on GitHub](https://github.com/dreamvo/gilfoyle/discussions)
- [Report a bug](https://github.com/dreamvo/gilfoyle/issues/new)
- Follow us on [Twitter](https://twitter.com/dreamvoapp)
- Contact us at [contact@dreamvo.com](mailto:contact@dreamvo.com)
