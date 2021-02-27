---
title: "Installation"
description: "Install the CLI application."
lead: ""
date: 2020-10-06T08:49:31+00:00
lastmod: 2020-10-06T08:49:31+00:00
draft: false
images: []
menu:
  docs:
    parent: "getstarted"
weight: 120
toc: true
---

To install Gilfoyle, you'll need to download the latest binary available on GitHub or build the software from source.

{{< alert icon="❗" text="Linux and MacOS are the only supported operating systems, with the following architectures : <strong>x86_64, arm64, armv6, i386</strong>." >}}

## Binary installation

Follow the instructions :

- Go to the [release page on GitHub](https://github.com/dreamvo/gilfoyle/releases)
- Choose a binary according to your OS and architecture
- Download the archive, extract the binary then run it in a terminal

You can also do all the steps above from the terminal :

```
os="$(uname -s)_$(uname -m)"
gilfoyle_version=$(curl -s https://api.github.com/repos/dreamvo/gilfoyle/releases/latest | grep tag_name | cut -d '"' -f 4)

# Downloading the latest version
wget "https://github.com/dreamvo/gilfoyle/releases/download/$gilfoyle_version/gilfoyle_$os.tar.gz"

# Verify checksum
curl -sSL https://github.com/dreamvo/gilfoyle/releases/download/$gilfoyle_version/gilfoyle_checksums.txt | sha256sum -c --strict --ignore-missing

# Clear directory
tar xfv "gilfoyle_$os.tar.gz"
rm gilfoyle_$os.tar.gz

# Use the binary
./gilfoyle version
```

To ensure your system is supported, please check the output of echo `"$(uname -s)_$(uname -m)"` in your terminal and see if it's available on the [GitHub release page](https://github.com/dreamvo/gilfoyle/releases).

## Docker

### From docker hub

You can pull the production-ready Docker image directly from Docker hub.

```shell
docker pull dreamvo/gilfoyle:latest
```

Then run the program :

```shell
docker run --rm -it dreamvo/gilfoyle version
```

Learn [how to configure a Gilfoyle instance →]({{< ref "config" >}})
