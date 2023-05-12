*Author: raphael@crvx.fr*

Gilfoyle was created with scale in mind. It's not only a storage server for your videos, it's a complete media streaming server. It handles file upload, processing, and streaming at your own scale.

This is intended to be a **high level** design document. Some of the implementation details are going to be decided after the prototyping phase.

# Overview

Over time, media streaming evolved in a good way. Users can host videos for free, they can enjoy fast and adaptative streaming. But these pros often come with some conditions or limitations. For example, you are required to agree to the service's business model (YouTube, Vimeo, ...), you have to agree with the upload rate limit, bad transcoding settings making your content looking ugly. Also, at an enterprise grade, your needs may be huge but you don't want to relay on another SaaS/PaaS to host your content, as it can become very expansive. So Gilfoyle is a simple start to setup a media streaming server at your own scale. It's simple to **setup**, to **scale**, it's **customizable**, also it's free and open source.

As you may ask, why would we create another open-source video streaming server? [PeerTube](https://github.com/Chocobozzz/PeerTube), [D.Tube](https://d.tube/), already does the job, right? Yes. Gilfoyle takes those projects as an acknowledgment to create an alternate solution to the common problem. Gilfoyle is also a bit different: it's not a social network or a federated video streaming platform. It's a self-hosted service that only handle video/audio hosting, transcoding and streaming.

To resume, Gilfoyle is a **self-hosted** and **open source** version of existing SaaS such as [api.video](https://api.video/), [Dailymotion Cloud](https://dmcloud.net/) or [mux](https://mux.com/).

## Goals and Non-Goals

### Goals

#### G1: Cloud-native

> Cloud native applications are built from the ground up—optimized for cloud scale and performance. They’re based on microservices architectures, use managed services, and take advantage of continuous delivery to achieve reliability and faster time to market. [Read more](https://azure.microsoft.com/en-us/overview/cloudnative/)

We want to provide a cloud-native application for high scale businesses. That's why it should be easy for any user to scale and mutate web services and databases on demand. Cloud-native also means extensible configuration for distributed environments, it defines how "cloud-friendly" your application is.

#### G2: Customizable

The service may have some configuration settings to be controlled by the user such as max file size, target transcoding format, compression rate... To achieve that, user would use a simple Yaml file that centralize these settings. If any config file is provided, default settings are used. Some open source projects can be difficult to use because of too many configuration settings. Gilfoyle is easy to use : simply download a binary, run it and access the web service. Want to deploy to production? Use the production-ready Docker image or see container orchestration examples.

#### G3: Stateless architecture, flexible storage

You can choose the appropriate storage system between: **filesystem** or **object storage**.

##### Filesystem

> Filesystem refers to the actual physical disk of the running machine.

This option is for stateful architectures.

##### Object storage

> Cloud storage is an external object storage system, such as [AWS S3](https://aws.amazon.com/s3/) (or any S3-compatible alternative), [OpenStack Object Storage (swift)](https://www.ovhcloud.com/en-gb/public-cloud/object-storage/) or [Google Cloud Storage](https://cloud.google.com/storage/).

This option is for businesses with high trafic, large and private files.

#### G5: Multimedia

The service handles both video and audio. It means you can use it to create your own podcast streaming platform, or your own video streaming platform.

### Non-goals

#### Security by Default

Althrough we take security very seriously, the API wasn't designed to be exposed to the public network. Usually user would isolate the interface in a security group with access restricted to other services. Gilfoyle shouldn't be the primary backend of your application, but a private storage service used by your own API. Still, you can deploy and expose publicly this service in production.

#### Social features

Gilfoyle is not another YouTube alternative. It doesn't provide social features such as likes, comments, channels or subscriptions.

#### Federation

Federation is for P2P-based platforms for which reliability is a top priority. As we want to prioritize business usage, we cannot support federation.

## Design

### Definitions

- **Media** : A multimedia entity representing either a video or audio only that users can download.
- **Media rendition** : A version of media with specific encoding parameters. 
- **Media format** : hls, dash, mp4
- **Media event** : Log events about a single media. Providing useful information about the status of a media and how's the encoding going.
- **Media probe** : Information about the source file of a media.
- **Media source file** : The original file provided by the end-user.
- **Storage driver** : The interface used to interact with a storage backend (e.g.: S3, local filesystem, GCS...).
- **Workload** : all the resources necessary to run a Gilfoyle instance

### Workload architecture

Gilfoyle is composed of the following parts : 

- A PostgreSQL database
- A RabbitMQ server
- A REST API
- A set of worker nodes

Worker nodes are processing messages from the RabbitMQ server directly. A Gilfoyle instance only needs 1 worker node to run properly but multiple nodes are recommanded in production.

### Lifecycle of a media

```
1. User creates a media --> Media A gets created with status "AwaitingUpload"
2. User uploads a file for media A --> File is stored, the media status change to "ScheduledProcessing", a job is scheduled to start the encoding process
3. In background, a worker process the jobs and start encoding all the necessary renditions. Media status change to "Processing"
4. If one rendition fails, the worker will reschedule the job X times (defined by user)
5. Once encoding jobs are finished, another job finalize the process and set the media status either to "Ready" or "Errored"
6. If the encoding failed, the administrator can look at media events to see what happened at each step
7. When media becomes available for streaming, user can retrieve the stream either using a CDN or the Gilfoyle API directly
```

### How we handle fail over

#### External services

We can't control the availability of external services so we'll simply delegate that to the user. User is responsible for the availability of the PSGL database and the RabbitMQ server.

#### Workers

We may experience issues in the following situations :

- The RabbitMQ client disconnect unexpectedly
- The worker panic while processing a job and terminate

...

### Dependencies

Of course, this list can evolve over time with no warrancy.

- Go
  - [swaggo/swag](https://github.com/swaggo/swag)
  - [gin-gonic/gin](https://github.com/gin-gonic/gin)
  - [facebook/ent](https://github.com/facebook/ent)
  - [spf13/cobra](https://github.com/spf13/cobra)
  - [lib/pq](https://github.com/lib/pq)
  - [google/uuid](https://github.com/google/uuid)
  - [jinzhu/configor](https://github.com/jinzhu/configor)
  - [yaml.v2](https://github.com/go-yaml)
  - [streadway/amqp](https://github.com/streadway/amqp)
  - [go.uber.org/zap](https://go.uber.org/zap)
  - [natefinch/lumberjack.v2](https://gopkg.in/natefinch/lumberjack.v2)
  - [vansante/go-ffprobe.v2](https://gopkg.in/vansante/go-ffprobe.v2)
  - [minio/minio-go/v7](https://github.com/minio/minio-go/v7)
  - [go-playground/validator/v10](https://github.com/go-playground/validator/v10)
- PostgreSQL
- RabbitMQ
- FFmpeg
- FFprobe

### Technical architecture

#### High level architecture

![high level architecture](https://i.imgur.com/iyhen9k.png)

#### External interfaces

We want documentation to be part of the code, so its always up-to-date and developers can understand snippets very quickly. We defines API's specifications using OpenAPI with swagger. The specification is defined in the code and used to generate a JSON file (see [swagger.yaml](https://github.com/dreamvo/gilfoyle/blob/master/api/docs/swagger.yaml)).
