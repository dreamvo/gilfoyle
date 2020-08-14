# Design

Gilfoyle was created with privacy, scale and security in mind. It's only a storage server for your videos, it's a complete video streaming platform backend. It handles video file upload, processing, and streaming.

## Why?

As you may ask, why would we create another open-source video streaming server? [PeerTube](), D.Tube, already does the job, right? But these projects has common limitations. ...

## Goals

1. Privacy & security by default

Gilfoyle was created in an effort to bring a new privacy and watch experience to end-users. This application collects very few things about end-user and tend to keep it that way. Althrough we take security very seriously, the API wasn't designed to be exposed to the public network. Usually administrator would isolate the interface in a security group with access limited to other services. Gilfoyle shouldn't be the backend of your application, but a private storage service used by your own API. Still, you can deploy and expose this service in production.

2. Customization

The service may have some configuration settings to be controlled by administrator such as max file size, target transcoding format, compression rate... To achieve that, administrator would use a simple Yaml file that centralize these settings. If any config file is provided, default settings are used. Some open source projects can be difficult to use because of too many configuration settings. Gilfoyle is easy to use : simply download a binary, run it and access the web service. Want to deploy to production? Use the production-ready Docker image or see container orchestration examples. Of course, this application was designed to follow your application's scale. You can even [scale your own IPFS cluster](https://cluster.ipfs.io/).

3. Documentation as code

We want documentation to be part of the code, so its always up-to-date and developers can understand snippets very quickly. We defines API's specifications using OpenAPI with swagger.
