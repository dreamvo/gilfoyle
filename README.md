# gilfoyle

<div align="left">
  <a href="https://godoc.org/github.com/dreamvo/gilfoyle">
    <img src="https://godoc.org/github.com/dreamvo/gilfoyle?status.svg" alt="GoDoc">
  </a>
  <a href="https://github.com/dreamvo/gilfoyle/actions">
    <img src="https://img.shields.io/endpoint.svg?url=https://actions-badge.atrox.dev/dreamvo/gilfoyle/badge?ref=master" alt="build status" />
  </a>
  <a href="https://goreportcard.com/report/github.com/dreamvo/gilfoyle">
    <img src="https://goreportcard.com/badge/github.com/dreamvo/gilfoyle" alt="go report" />
  </a>
  <a href="https://github.com/dreamvo/gilfoyle/releases">
    <img src="https://img.shields.io/github/release/dreamvo/gilfoyle.svg" alt="Latest version" />
  </a>
</div>

Gilfoyle is a web application from the [Dreamvo project](https://dreamvo.com) that runs a self-hosted video streaming server. This application allows you to setup a production-ready video hosting & live streaming platform in minutes. Gilfoyle handles video upload, processing and streaming.

It's written in Golang and runs as a single Linux binary with PostgreSQL and Redis.

<a href="https://www.redbubble.com/fr/people/andromeduh/shop"><img src="https://ih1.redbubble.net/image.71449494.3195/raf,750x1000,075,t,oatmeal_heather.u2.jpg" width="256" align="right" /></a>

## Table of content

- [Features](#features)
- [Current status](#current-status)
- [Design](#design)
- [Documentation](#documentation)
- [Discussion](#discussion)

## Features

- Deploy a RESTful API to manage your videos
- Upload files or import from third-party platforms such as: *YouTube, Dailymotion, Vimeo*
- Handle video compression and transcoding with [FFmpeg](https://ffmpeg.org/)
- Decentralize video storage with [IPFS](https://ipfs.io/) clustering feature
- Enjoy [IPFS](https://ipfs.io/)'s cache & CDN features
- Collect analytics such as watch time and view count

## Current status

As this project is very recent, it's under heavy development and not suitable for production yet. Please consider v0 as unstable. Want to contribute ? Check the [backlog](https://github.com/dreamvo/gilfoyle/projects/1).

## Design

Read full specifications & background in [DESIGN.md](DESIGN.md).

## Documentation

- For **developers**: see [godoc](https://godoc.org/github.com/dreamvo/gilfoyle), [design documentation](DESIGN.md)
- For **administrators**: see [user guide](https://dreamvo.github.io/gilfoyle/) and [API documentation](https://redocly.github.io/redoc/?url=https://raw.githubusercontent.com/dreamvo/gilfoyle/master/api/docs/swagger.yaml#tag/videos)
- The [GPL v3](LICENSE) license

## Discussion

- [Report a bug](https://github.com/dreamvo/gilfoyle/issues/new)
- Follow us on [Twitter](https://twitter.com/dreamvoapp)
- contact us [contact@dreamvo.com](mailto:contact@dreamvo.com)
