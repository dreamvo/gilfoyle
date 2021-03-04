---
title: "Introduction"
description: "Gilfoyle is a web application from the Dreamvo project that runs a self-hosted media processing server."
lead: "Welcome to the Gilfoyle documentation! Before you learn how to install and deploy it, you may want to understand the design and the motivation that drives the project."
date: 2021-02-27T08:48:57+00:00
lastmod: 2021-02-27T08:48:57+00:00
draft: false
images: []
menu:
  docs:
    parent: "prologue"
weight: 100
toc: true
---

{{< alert icon="❗" text="As this project is very recent, it's under heavy development and not suitable for production yet. Please consider version 0 as unstable." >}}

Over time, media streaming evolved in a good way. Users can host videos for free, they can enjoy fast and adaptative streaming. But these pros often come with some conditions or limitations. For example, you are required to agree to the service's business model (YouTube, Vimeo, ...), you have to agree with the upload rate limit, bad transcoding settings making your content looking ugly. Also, at an enterprise grade, your needs may be huge but you don't want to relay on another SaaS/PaaS to host your content, as it can become very expansive. So Gilfoyle is a simple start to setup a media streaming server at your own scale. It's simple to **setup**, to **scale**, it's **customizable**, also it's free and open source.

As you may ask, why would we create another open-source video streaming server? [PeerTube](https://github.com/Chocobozzz/PeerTube), [D.Tube](https://d.tube/), already does the job, right? Yes. Gilfoyle takes those projects as an acknowledgment to create another solution to the common problem. Gilfoyle is also a bit different: it's not a social network or a federated video streaming platform. It's a self-hosted service that only handle video/audio hosting, transcoding and streaming.

To resume, Gilfoyle is a **self-hosted** and **open source** alternative to existing SaaS (Software-As-A-Service) such as [api.video](https://api.video/), [Dailymotion Cloud](https://dmcloud.net/), [Red5pro](https://www.red5pro.com/) or [mux](https://mux.com/).

## Goals and Non-Goals

### Goals

#### G1: Cloud-native

> Cloud native applications are built from the ground up—optimized for cloud scale and performance. They’re based on microservices architectures, use managed services, and take advantage of continuous delivery to achieve reliability and faster time to market. [Read more](https://azure.microsoft.com/en-us/overview/cloudnative/)

We want to provide a cloud-native application for high scale businesses. That's why it should be easy for any user to scale and mutate web services and databases on demand. Cloud-native also means extensible configuration for distributed environments, it defines how "cloud-friendly" your application is.

#### G2: Customizable

The service may have some configuration settings to be controlled by the user such as max file size, target transcoding format, compression rate... To achieve that, user would use a simple Yaml file that centralize these settings. If any config file is provided, default settings are used.

#### G3: Stateless architecture, flexible storage

You can choose the appropriate storage system between: **filesystem** or **object storage**.

##### Filesystem

> Filesystem refers to the actual physical disk of the running machine.

This option is for stateful architectures.

##### Object storage

> Cloud storage is an external object storage system, such as [AWS S3](https://aws.amazon.com/s3/) (or any S3-compatible alternative), [OpenStack Object Storage (swift)](https://www.ovhcloud.com/en-gb/public-cloud/object-storage/) or [Google Cloud Storage](https://cloud.google.com/storage/).

This option is for businesses with high trafic, large and private files.

#### G5: Multimedia

The service handles both video and audio. It means you can rely on it to create your own podcast streaming platform, or your own video streaming platform.

#### G6: Agnostic

The media streaming industry include various use-cases with specific needs. For example, YouTube and Netflix both needs VOD streaming but does not handle encoding and streaming the same way. They both handle different features because of their use-case. We try our best to make Gilfoyle fit all of those requirements, without favouring a specific use-case.

### Non-goals

#### NG1: Security by Default

Althrough we take security very seriously, the API wasn't designed to be exposed to the public network. Usually user would isolate the interface in a security group with access restricted to other services. Gilfoyle shouldn't be the primary backend of your application, but a private storage service used by your own API. Still, you can deploy and expose publicly this service in production.

#### NG2: Social features

Gilfoyle is not another YouTube alternative. It doesn't provide social features such as likes, comments, channels or subscriptions.

#### NG3: Federation

Federation is for P2P-based platforms for which reliability is a top priority. As we want to prioritize business usage, we cannot support federation.
