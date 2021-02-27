---
title: "How it works"
description: "Learn how Gilfoyle works."
date: 2021-02-27T08:48:57+00:00
lastmod: 2021-02-27T08:48:57+00:00
draft: false
images: []
menu:
  docs:
    parent: "prologue"
weight: 110
toc: true
---

Gilfoyle was designed to be hidden from the public network, so only your internal services can communicate with it. This way, you can create a proxy gateway for your users, or simply add your own logic before sending user data to Gilfoyle.

## Architecture

Here's an overview of how your application interact with Gilfoyle.

<img width="100%" src="/images/architecture_draw_1.png" />

When a client uploads a media file, it is processed by a background job that perform analysis on the file then start encoding the media in several renditions, according to the provided settings.
