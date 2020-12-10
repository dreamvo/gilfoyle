# gilfoyle [![godoc](https://godoc.org/github.com/dreamvo/gilfoyle?status.svg)](https://godoc.org/github.com/dreamvo/gilfoyle) [![build](https://img.shields.io/endpoint.svg?url=https://actions-badge.atrox.dev/dreamvo/gilfoyle/badge?ref=master)](https://github.com/dreamvo/gilfoyle/actions) [![goreport](https://goreportcard.com/badge/github.com/dreamvo/gilfoyle)](https://goreportcard.com/report/github.com/dreamvo/gilfoyle) [![Coverage Status](https://coveralls.io/repos/github/dreamvo/gilfoyle/badge.svg)](https://coveralls.io/github/dreamvo/gilfoyle) [![release](https://img.shields.io/github/release/dreamvo/gilfoyle.svg)](https://github.com/dreamvo/gilfoyle/releases)

Gilfoyle is a web application from the [Dreamvo project](https://dreamvo.com) that runs a self-hosted video streaming server. This application allows you to setup a production-ready video hosting & live streaming platform in minutes. Gilfoyle handles video upload, processing and streaming.

It's written in Golang and runs as a single Linux binary with PostgreSQL and Redis.

## Table of content

- [Features](#features)
- [Current status](#current-status)
- [Design](#design)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [Discussion](#discussion)

## Features

- **Deploy** a RESTful API and HLS server to manage and stream audio & video
- **Upload** files or import from third-party platforms such as: *YouTube, Dailymotion, Vimeo*
- Handle video compression and transcoding with [FFmpeg](https://ffmpeg.org/)
- **Decentralize** video storage with [IPFS](https://ipfs.io/)
- Enjoy [IPFS](https://ipfs.io/)'s cache & CDN features
- **Customize**: video transcoding (format, resolution, frame rate), compression rate, authentication

## Current status

It's a **Work In Progress**. As this project is very recent, it's under heavy development and not suitable for production yet. Please consider v0 as unstable. Want to contribute ? Check the [backlog](https://github.com/dreamvo/gilfoyle/projects/1).

## Design

See [this document](DESIGN.md) for a high level design and goals.

## Documentation

- For **developers**: see [godoc](https://godoc.org/github.com/dreamvo/gilfoyle), [design documentation](DESIGN.md)
- For **administrators**: see [user guide](https://dreamvo.github.io/gilfoyle/) (WIP) and [API documentation](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/dreamvo/gilfoyle/master/api/docs/swagger.json)
- The [GPL v3](LICENSE) license

## Contributing

This project is in an early stage. We appreciate feedbacks and [discussions](#discussion) about the design and [features](#features).

## Discussion

- [Report a bug](https://github.com/dreamvo/gilfoyle/issues/new)
- Follow us on [Twitter](https://twitter.com/dreamvoapp)
- Contact us at [contact@dreamvo.com](mailto:contact@dreamvo.com)
