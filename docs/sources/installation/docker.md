+++
title = "Installing using Docker"
description = "Installing Grafana using Docker guide"
keywords = ["grafana", "configuration", "documentation", "docker"]
type = "docs"
[menu.docs]
name = "Installing using Docker"
identifier = "docker"
parent = "installation"
weight = 4
+++

# Installing using Docker

Grafana is very easy to install and run using the offical docker container.

```bash
$ docker run -d -p 3000:3000 grafana/grafana
```

## Configuration

All options defined in conf/grafana.ini can be overridden using environment
variables by using the syntax `GF_<SectionName>_<KeyName>`.
For example:

```bash
$ docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  -e "GF_SERVER_ROOT_URL=http://grafana.server.name" \
  -e "GF_SECURITY_ADMIN_PASSWORD=secret" \
  grafana/grafana
```

The back-end web server has a number of configuration options. Go to the
[Configuration]({{< relref "configuration.md" >}}) page for details on all
those options.

## Running a Specific Version of Grafana

```bash
# specify right tag, e.g. 5.0.3 - see Docker Hub for available tags
$ docker run \
  -d \
  -p 3000:3000 \
  --name grafana \
  grafana/grafana:5.0.3
```

## Installing Plugins for Grafana

Pass the plugins you want installed to docker with the `GF_INSTALL_PLUGINS` environment variable as a comma separated list. This will pass each plugin name to `grafana-cli plugins install ${plugin}`.

```bash
docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  -e "GF_INSTALL_PLUGINS=grafana-clock-panel,grafana-simple-json-datasource" \
  grafana/grafana
```

## Building a custom Grafana image with pre-installed plugins

#### Dockerfile
```Dockerfile
FROM grafana/grafana:5.0.0
ENV GF_PATHS_PLUGINS=/opt/grafana-plugins
RUN mkdir -p $GF_PATHS_PLUGINS
RUN grafana-cli --pluginsDir $GF_PATHS_PLUGINS plugins install grafana-clock-panel
```

Add lines with `RUN grafana-cli ...` for each plugin you wish to install in your custom image. Don't forget to specify what version of Grafana you wish to build from (replace 5.0.0 in the example).

Example of how to build and run:
```bash
docker build -t grafana:5.0.0-custom . 
docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  grafana:5.0.0-custom
```

## Configuring AWS Credentials for CloudWatch Support

```bash
$ docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  -e "GF_AWS_PROFILES=default" \
  -e "GF_AWS_default_ACCESS_KEY_ID=YOUR_ACCESS_KEY" \
  -e "GF_AWS_default_SECRET_ACCESS_KEY=YOUR_SECRET_KEY" \
  -e "GF_AWS_default_REGION=us-east-1" \
  grafana/grafana
```

You may also specify multiple profiles to `GF_AWS_PROFILES` (e.g.
`GF_AWS_PROFILES=default another`).

Supported variables:

- `GF_AWS_${profile}_ACCESS_KEY_ID`: AWS access key ID (required).
- `GF_AWS_${profile}_SECRET_ACCESS_KEY`: AWS secret access  key (required).
- `GF_AWS_${profile}_REGION`: AWS region (optional).

## Grafana container with persistent storage (recommended)

```
# create a persistent volume for your data in /var/lib/grafana (database and plugins)
docker volume create grafana-storage

# start grafana
docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  -v grafana-storage:/var/lib/grafana \
  grafana/grafana
```

Note: An unnamed volume will be created for you when you boot Grafana,
using `docker volume create grafana-storage` just makes it easier to find
by giving it a name.

## Grafana container using bind mounts

You may want to run Grafana in Docker but use folders on your host for the database or configuration. When doing so it becomes important to start the container with a user that is also able to access and write to the folder you map into the container.

```bash
mkdir data # creates a folder for your data
ID=$(id -u) # saves your user id in the ID variable

# starts grafana with your user id and using the data folder
docker run -d --user $ID --volume "$PWD/data:/var/lib/grafana" -p 3000:3000 grafana/grafana:5.1
```

## Migration from a previous version of the docker container to 5.1 or later

In 5.1 we switched from running as the `grafana` user to the `nobody` user. Unfortunately this means that files created prior to 5.1 won't have the correct permissions for later versions. We made this change so that it would be easier for you to control what user Grafana is executed as (see examples below).

Version | User    | User ID
--------|---------|---------
< 5.1   | grafana | 104
>= 5.1  | nobody  | 65534

There are two possible solutions to this problem. Either you start the new container as the root user and changes ownership from `104` to `65534` or you start the upgraded container as user `104`.

### Running docker as a different user

```bash
docker run --user 104 --volume "<your volume mapping here>" grafana/grafana:5.1
```

#### docker-compose.yml with custom user
```yaml
version: "2"

services:
  grafana:
    image: grafana/grafana:5.1
    ports:
      - 3000:3000
    user: "104"
```

### Modifying permissions

Always be careful when modifying permissions.

```bash
$ docker run -ti --user root --volume "<your volume mapping here>" --entrypoint bash grafana/grafana:5.1

# in the container you just started:
chown -R root:root /etc/grafana && \
  chmod -R a+r /etc/grafana && \
  chown -R nobody:nogroup /var/lib/grafana && \
  chown -R nobody:nogroup /usr/share/grafana
```
