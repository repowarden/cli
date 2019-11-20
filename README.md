# Warden [![CI Status](https://circleci.com/gh/repowarden/cli.svg?style=shield)](https://app.circleci.com/pipelines/github/repowarden/cli) [![Software License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/repowarden/cli/trunk/LICENSE)

The RepoWarden CLI `warden` is a tool to audit your git repositories based on policy.
Check the default branch, license, and labels across multiple repositories across multiple organizations.

Warden is very early stage right now with limited features.


## Table of Contents

- [Compatibility](#compatibility)
- [Installation](#installation)
  - [Linux](#linux)
  - [macOS](#macos)
  - [Windows](#windows)
- [Configuring](#configuring)
- [Features](#features)


## Compatibility

### Operating Systems

Designed to work on Linux, macOS, and Windows computers.

### VCS Providers

GitHub is supported.
GitLab support is on the roadmap.


## Installation

### Linux

#### Debian Package (.deb)
You can install `warden` on an Debian/Apt based computer by downloading the `.deb` file to the desired system.

For graphical systems, you can download it from the [GitHub Releases page][gh-releases].
Many distros allow you to double-click the file to install.
Via terminal, you can do the following:

```bash
wget https://github.com/repowarden/cli/releases/download/v0.1.0/warden_0.1.0_amd64.deb
sudo dpkg -i warden_0.1.0_amd64.deb
```

`0.1.0` and `amd64` may need to be replaced with your desired version and CPU architecture respectively.

#### Binary Install
You can download and run the raw `warden` binary from the [GitHub Releases page][gh-releases] if you don't want to use any package manager.
Simply download the tarball for your OS and architecture and extract the binary to somewhere in your `PATH`.
Here's one way to do this with `curl` and `tar`:

```bash
dlURL="https://github.com/repowarden/cli/releases/download/v0.1.0/warden-v0.1.0-linux-amd64.tar.gz"
curl -sSL $dlURL | sudo tar -xz -C /usr/local/bin warden
```

`0.1.0` and `amd64` may need to be replaced with your desired version and CPU architecture respectively.

### macOS

There are two ways you can install `warden` on a macOS system.

#### Brew (recommended)

Installing Warden via brew is a simple one-liner:

```bash
brew install repowarden/tap/warden
```

#### Binary Install
You can download and run the raw `warden` binary from the [GitHub Releases page][gh-releases] if you don't want to use any package manager.
Simply download the tarball for your OS and architecture and extract the binary to somewhere in your `PATH`.
Here's one way to do this with `curl` and `tar`:

```bash
dlURL="https://github.com/repowarden/cli/releases/download/v0.1.0/warden-v0.1.0-macos-amd64.tar.gz"
curl -sSL $dlURL | sudo tar -xz -C /usr/local/bin warden
```

`0.1.0` and `amd64` may need to be replaced with your desired version and CPU architecture respectively.

### Windows

Warden supports Windows 10 by downloading and installing the binary.
Chocolately support is likely coming in the future.
If there's a Windows package manager you'd like support for (including Chocolately), please open and Issue and ask for it.

#### Binary Install (exe)
You can download and run the `warden` executable from the [GitHub Releases page][gh-releases].
Simply download the zip for architecture and extract the exe.


## Configuring

**credentials** - the credentials file is `~/.config/warden/creds.yaml`.
The key `githubtoken` should be set to a token that has enough permissions to do what you need.

**policies** - the policy file, `warden.yml`, should be in the current directory.
You can get started by copying over the example one: `cp example.warden.yml warden.yml`


## Features

Currently `warden` can audit the following items:

- license
- labels
- default branch

Run `warden help` to see all commands available.


## License

This repository is licensed under the MIT license.
The license can be found [here](./LICENSE).



[gh-releases]: https://github.com/repowarden/cli/releases
