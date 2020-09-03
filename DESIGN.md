*Author: raphael@crvx.fr*

Gilfoyle was created with privacy, scale and security in mind. It's only a storage server for your videos, it's a complete video streaming platform backend. It handles video file upload, processing, and streaming.

This is intended to be a **high level** design document. Some of the implementation details are going to be decided after the prototyping phase.

# Overview

Over time, media streaming evolved in a good way. Users can host videos for free, they can enjoy fast and adaptative streaming. But these pros often come with some conditions or limitations. For example, you are required to agree to the service's business model (YouTube, Vimeo, ...), you have to agree with the upload rate limit, bad transcoding settings making your content looking ugly. Also, at an enterprise grade, your needs may be huge but you don't want to relay on another SaaS/PaaS to host your content. So Gilfoyle is a simple start to setup a media streaming server at your own scale. It's simple to **setup**, to **scale**, it's **customizable**, it's **decentralized** over the IPFS network, also it's free and open source.

As you may ask, why would we create another open-source video streaming server? [PeerTube](https://github.com/Chocobozzz/PeerTube), [D.Tube](https://d.tube/), already does the job, right? Yes. Gilfoyle takes those projects as an aknownledgment to a better solution to the common problem. Gilfoyle is also a bit different: it's not a social network or a federated video streaming platform. It's a self-hosted service that only handle video/audio hosting, transcoding and streaming.

## Goals and Non-Goals

### Goals

#### G1: Performances & Scale

We want to provide a efficient product for high scale businesses. It should be easy for any administrator to scale web service and databases on demand, because of distributed services. For example, the technical choice of etcd over Redis or CockroachDB over PostgreSQL can make the difference. You can choose to scale the server or the storage (IPFS Swarm) as you wish, independently.

#### G2: Privacy & Security by Default

Gilfoyle was created in an effort to bring a new privacy and watch experience to end-users. This application collects very few things about end-user and tend to keep it that way. Althrough we take security very seriously, the API wasn't designed to be exposed to the public network. Usually administrator would isolate the interface in a security group with access restricted to other services. Gilfoyle shouldn't be the primary backend of your application, but a private storage service used by your own API. Still, you can deploy and expose publicly this service in production.

#### G3: Customizable

The service may have some configuration settings to be controlled by administrator such as max file size, target transcoding format, compression rate... To achieve that, administrator would use a simple Yaml file that centralize these settings. If any config file is provided, default settings are used. Some open source projects can be difficult to use because of too many configuration settings. Gilfoyle is easy to use : simply download a binary, run it and access the web service. Want to deploy to production? Use the production-ready Docker image or see container orchestration examples. Of course, this application was designed to follow your application's scale. You can even [scale your own IPFS cluster](https://cluster.ipfs.io/).

#### G4: Decentralized

...

#### G5: Cloud-native

...

### Non-goals

#### Extensible

...

#### Social features

...

#### User tracking

...

## Design

TODO

### Dependencies

...

### Technical architecture

#### Application file structure

...

#### External interfaces

We want documentation to be part of the code, so its always up-to-date and developers can understand snippets very quickly. We defines API's specifications using OpenAPI with swagger. The specification is defined in the code and used to generate a JSON file (`api/docs/swagger.json`).

### Security

TODO

### Features

#### Video upload

...

#### Video processing

...
