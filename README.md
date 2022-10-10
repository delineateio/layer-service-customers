[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <img alt="delineate.io" src="https://github.com/delineateio/.github/blob/master/assets/logo.png?raw=true" height="75" />
  <h2 align="center">delineate.io</h2>
  <p align="center">portray or describe (something) precisely.</p>

  <h3 align="center">Microservices Dynamic Configuration</h3>

  <p align="center">
    Demonstrates of centralised dynamic configuration of microservices without downtime.
    <br />
    <br />
    <a href="https://github.com/delineateio/hashicorp-consul-kv-example/issues">Report Bug</a>
    ·
    <a href="https://github.com/delineateio/hashicorp-consul-kv-example/issues">Request Feature</a>
  </p>
</p>

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [About The Project](#about-the-project)
  - [Problem with Env Variables](#problem-with-env-variables)
- [Built With](#built-with)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Initialisation](#initialisation)
- [Usage](#usage)
  - [Customers Service Implementation](#customers-service-implementation)
  - [Running Natively](#running-natively)
  - [Running in Docker](#running-in-docker)
  - [Testing](#testing)
  - [Full Make Targets](#full-make-targets)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

<!-- ABOUT THE PROJECT -->
## About The Project

This repo provides an example of using [Hashicorp Consul](https://www.consul.io/) KV store for centralised and dynamic service configuration that enables zero downtime reconfiguration of existing running microservices.

### Problem with Env Variables

Config is one of the 12 factors specified building for [12 Factor Apps](https://12factor.net/). Most guidance is to use env variables to achieve externalisation of configuration, such as [this](https://12factor.net/config).

However there are some common challenges with using env variables:

* In highly configurable apps this can result in env variable sprawl
* Env variables to not have any native way to group or hierarchically structure config
* Updates to env variables do not support atomic transactions
* Most hosting platforms require redeployment to update env variables

## Built With

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)
![Python](https://img.shields.io/badge/python-3670A0?style=for-the-badge&logo=python&logoColor=ffdd54)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Nginx](https://img.shields.io/badge/nginx-%23009639.svg?style=for-the-badge&logo=nginx&logoColor=white)

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

The following should be installed locally.

* [Homebrew](https://brew.sh/)
* [Docker](https://www.docker.com/)

### Initialisation

**Note is is assumed most developers are using MacOS, if not then some dependencies will need to be installed manually.**

To initialise the local development environment run `make init`, this will:

* Install `brew` dependencies using [brew bundle](https://github.com/Homebrew/homebrew-bundle)
* Generates a local certificate authority and certificates for `consul`
* Creates a python virtual environment and install [PyPi](https://pypi.org/) dependencies
* Sets up [pre-commit](https://pre-commit.com/) hooks

You can inspect exactly what gets installed by `brew` and `pip` using `make list`.

<!-- USAGE EXAMPLES -->
## Usage

### Customers Service Implementation

An example service written in `go` has been provided.  The *customers* service leverages a set of modules for developing production services:

| Go Module | Purpose |
| --- | ----------- |
| [gin](https://github.com/gin-gonic/gin) | Web framework used for the APIs. |
| [gin swagger](github.com/swaggo/gin-swagger) | Used to expose swagger documentation for the API. |
| [viper](https://github.com/spf13/viper) | Module for configuration management.  |
| [consul structure](https://github.com/mitchellh/consulstructure)| Specific library which exposes an interface to interact with Consul. |
| [zap](https://github.com/uber-go/zap) | High performance logging framework from [Uber](https://www.uber.com/gb/en/about/). |
| [gorm](https://github.com/go-gorm/gorm) | Native `go` Object Relational Mapper (ORM). |
| [retry-go](https://github.com/avast/retry-go) | Simple retry library from [Avast](https://www.avast.com/en-gb). |

The example service API starts on port `1102` by default.  The details of the API endpoints can be explored by browsing the `swagger` documentation where the service is running [here](http://localhost:1102/swagger/index.html).

### Running Natively

To run, test and debug the `go` service natively (not inside `docker`) it is necessary to have the following background services running in advance:

* `consul` running on port `8500`
* `postgres` running on port `5432`

Running`make services` will stand up the required backend services.

When the services are running it is possible to browse the [Consul UI](http://localhost:8500/ui/).  You can subsequently use `make ps` to show the running containers.

In addition `psql` is installed then it's possible to connect to the database using `psql -U postgres -h localhost`.  You will need to `source .env` file to avoid being prompted for the password each time.

### Running in Docker

One of the key advantages of using centralised and dynamic configuration is that config is updated across multiple services.  To demonstrate this a `make` target has been provided that instantiates three instances of the *customers* service in Docker behind an [nginx](https://www.nginx.com/) load balancer.

```shell
# builds the image, scales to 3 by default
make build=true up

# doesn't rebuild the image, scales up 5
make scale=5 up

# make scales back down to 2
make scale=2 up
```

The *customers* service is packaged using [Build Packs](https://buildpacks.io/) so there is no explicit `Dockerfile`.

### Testing

Once the service is running either natively or deployed in Docker a set of tests have been provided.  These tests use [BDD](https://cucumber.io/docs/bdd/) and are implemented by extending [behave](https://pypi.org/project/behave/).  These tests can be run using `make tests`.

### Full Make Targets

This section lists

| Target| Purpose |
| --- | ----------- |
| `make list` | Lists the `brew` and `pip` dependencies that will be installed. |
| `make init` | Full initialises the local development environment.|
| `make graph` | Writes a file called `graph.txt` which contains the `go` dependency tree for inspection. |
| `make build` | Builds the application binary for the application in `./build`. |
| `make services` | Creates the backend services to enable the application to be run natively. |
| `make up` | Deploys the application into Docker and scales to three instances. |
| `make ps` | Lists the primary info and status of the docker containers. |
| `make tests` | Runs the test feature files that have been provided in `./tests`. |
| `make rename` | Only required in exceptional cased to rename the repo and update `readme` badges. |

<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/delineateio/hashicorp-consul-kv-example/issues) for a list of proposed features (and known issues).

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

If you would like to contribute to any Capco Digital OSS projects please read:

* [Code of Conduct](https://github.com/delineateio/.github/blob/master/CODE_OF_CONDUCT.md)
* [Contributing Guidelines](https://github.com/delineateio/.github/blob/master/CONTRIBUTING.md)

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements

* [Best README Template](https://github.com/othneildrew/Best-README-Template/blob/master/README.md)
* [Markdown Badges](https://github.com/Ileriayo/markdown-badges)

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/delineateio/hashicorp-consul-kv-example.svg?style=for-the-badge
[contributors-url]: https://github.com/delineateio/hashicorp-consul-kv-example/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/delineateio/hashicorp-consul-kv-example.svg?style=for-the-badge
[forks-url]: https://github.com/delineateio/hashicorp-consul-kv-example/network/members
[stars-shield]: https://img.shields.io/github/stars/delineateio/hashicorp-consul-kv-example.svg?style=for-the-badge
[stars-url]: https://github.com/delineateio/hashicorp-consul-kv-example/stargazers
[issues-shield]: https://img.shields.io/github/issues/delineateio/hashicorp-consul-kv-example.svg?style=for-the-badge
[issues-url]: https://github.com/delineateio/hashicorp-consul-kv-example/issues
[license-shield]: https://img.shields.io/github/license/delineateio/hashicorp-consul-kv-example.svg?style=for-the-badge
[license-url]: https://github.com/delineateio/hashicorp-consul-kv-example/blob/master/LICENSE
