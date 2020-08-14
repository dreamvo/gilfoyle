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

Gilfoyle is a web application from the Dreamvo project that runs a self-hosted video streaming server. This application allows you to setup a production-ready video hosting & live streaming platform in minutes. Gilfoyle handles video upload, processing and streaming.

It's written in Golang and runs as a single Linux binary with PostgreSQL and Redis.

<a href="https://www.redbubble.com/fr/people/andromeduh/shop"><img src="https://ih1.redbubble.net/image.71449494.3195/raf,750x1000,075,t,oatmeal_heather.u2.jpg" width="256" align="right" /></a>

## Table of content

- [Features](#features)
- [Current status](#current-status)
- [Design](#design)
  - [Goals](#goals)
- [Documentation](#documentation)
- [License](#license)

## Features

- Deploy a RESTful API to manage your videos
- Upload videos from file or import from platforms such as: *YouTube, Dailymotion, Vimeo*
- Handle video compression and conversion asynchronously with [FFmpeg](https://ffmpeg.org/)
- Decentralize video storage with [IPFS](https://ipfs.io/) clustering feature
- Enjoy [IPFS](https://ipfs.io/)'s cache & CDN features
- Collect analytics such as watch time and view count

## Current status

As this project is very recent, it's under heavy development and not suitable for production yet. Please consider v0 as instable. Want to contribute ? Check the [backlog](https://github.com/dreamvo/gilfoyle/projects/1).

## Design

Read full specifications on the [dedicated website](https://dreamvo.github.io/specs/).

### Goals

1. Privacy & security by default

Gilfoyle was created in an effort to bring a new privacy and watch experience to end-users. This application collects very few things about end-user and tend to keep it that way. Althrough we take security very seriously, the API was not designed to be exposed to the public network. Usually administrator would isolate the interface in a security group with access limited to other services. Gilfoyle shouldn't be the backend of your application but a private storage service used by your own API. Still, you can deploy and expose this service in production.

2. Customization

The service may have some configuration settings to be controlled by administrator such as max file size, target transcoding format, compression rate... To achieve that, administrator would use a simple Yaml file that centralize these settings. If any config file is provided, default settings are used. Some open source projects can be difficult to use because of too many configuration settings. Gilfoyle is easy to use : simply download a binary, run it and access the web service. Want to deploy to production? Use the production-ready Docker image or see container orchestration examples. Of course, this application was designed to follow your application's scale. You can even [scale your own IPFS cluster](https://cluster.ipfs.io/).

3. Documentation as code

We want documentation to be part of the code, so its always up-to-date and developers can understand snippets very quickly. We defines API's specifications using OpenAPI with swagger.

## Documentation

- For **developers**: see [godoc](https://godoc.org/github.com/dreamvo/gilfoyle)
- For **administrators**: see [API documentation](https://redocly.github.io/redoc/?url=https://raw.githubusercontent.com/dreamvo/gilfoyle/master/api/docs/swagger.yaml#tag/videos) and [user guide](#)

## License

This project is licensed under the [GPL v3](LICENSE) license.
