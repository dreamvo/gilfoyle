---
title: "Deployment"
description: "Deploy to production."
lead: ""
date: 2020-11-12T15:22:20+01:00
lastmod: 2020-11-12T15:22:20+01:00
draft: false
images: []
menu:
  docs:
    parent: "getstarted"
weight: 140
toc: true
---

## Docker

You can use a single docker-compose file to run the tool without downloading the source code.

```
version: '3.7'

services:
    api:
      container_name: gilfoyle-api
      restart: on-failure
      image: dreamvo/gilfoyle:latest
      command: serve -p 5000
      ports:
        - "80:5000"

    worker:
      container_name: gilfoyle-worker
      restart: on-failure
      image: dreamvo/gilfoyle:latest
      command: worker

    dashboard:
      container_name: gilfoyle-dashboard
      restart: on-failure
      image: dreamvo/gilfoyle:latest
      command: dashboard -p 5000 --endpoint http://api:5000
      ports:
        - "81:5000"
```

### Troubleshooting

All output is sent to stdout so it can be inspected by running:

```shell
docker logs -f <container-id|container-name>
```

