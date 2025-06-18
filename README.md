# Cloud Command Line

This repository contains the CLI tool for performing actions against
our [cloud platform](https://github.com/vaguevoid/cloud-platform) hosted
in production on [play.void.dev](https://play.void.dev)

## Quick Start

We use [Go](https://go.dev) (v1.24.2) and the [just](https://just.systems/) task runner for development.

It is assumed:

  * You have [Go](https://go.dev/doc/install) installed
  * You have [just](https://just.systems/man/en/packages.html) installed
  * You have a [.env](.env.example) file for your local ENV

With these system dependencies in place you can get started with the following development tasks:

```bash
> just list    # list all available tasks
> just run     # run the command line executable
> just test    # run all tests
> just lint    # run lint tools (vet and staticcheck)
> just cover   # run code coverage (CLI)
```

When you `just run` you should see the default CLI output listing all available commands...

```bash
NAME:
   void-cloud - access to the Void Cloud Platform

USAGE:
   void-cloud [global options] [command [command options]]

VERSION:
   0.0.1

COMMANDS:
   login    tell us who you are
   deploy   share your game with others
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --server string, -s string  server endpoint (default: "https://play.void.dev/") [$SERVER]
   --help, -h                  show help
   --version, -v               print the version
```
## Login Command

```bash
NAME:
   void-cloud login - tell us who you are

USAGE:
   void-cloud login

OPTIONS:
   --server URL  server endpoint URL (default: "https://play.void.dev/") [$SERVER]
   --help, -h    show help
```

## Deploy Command

```bash
NAME:
   void-cloud deploy - share your game with others

USAGE:
   void-cloud deploy PATH [LABEL]

OPTIONS:
   --server URL       server endpoint URL (default: "https://play.void.dev/") [$SERVER]
   --org string       organization ID [$ORG]
   --game string      game ID [$GAME]
   --token string     personal access TOKEN [$TOKEN]
   --concurrency int  deploy CONCURRENCY (default: 8) [$CONCURRENCY]
   --help, -h         show help
```

> See the [justfile](./justfile) for all available tasks
