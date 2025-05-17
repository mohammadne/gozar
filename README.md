# Gozar

![Go Version](https://img.shields.io/badge/Golang-1.24-66ADD8?style=for-the-badge&logo=go)
![App Version](https://img.shields.io/github/v/tag/mohammadne/gozar?sort=semver&style=for-the-badge&logo=github)
![Repo Size](https://img.shields.io/github/repo-size/mohammadne/gozar?logo=github&style=for-the-badge)

**Gozar** helps you configure an Xray proxy to bypass internet filtering and send a "hello" packet to the world — effortlessly and effectively.

## Introduction

**Gozar** is a lightweight, CLI-based Go application that simplifies the process of generating Xray configuration files and setting up your Xray server. It’s fast, minimal, and astonishingly easy to use.

## Generation

You can generate Xray configuration files using one of the following methods:

### Go (locally)

```sh
git clone ssh://git@github.com/mohammadne/gozar
cd gozar
go run main.go executer
```

### Podman

```bash
podman run --name "gozar" ghcr.io/mohammadne/gozar:v0.2.1
```

Then, in another shell:

```bash
podman exec -it gozar sh
podman cp gozar:/app/outputs ./outputs
podman rm -f gozar
```

## Setup

> **Note**: Docker must be installed on both your local machine and the target server where you intend to run the Xray proxy.

An Ansible role is included to automate Docker installation on the target server.

After ensuring Docker is installed, run the setup script to configure the environment and install required tools:

```bash
./scripts/setup.sh
```

## Contribution 📬

Pull requests, bug reports, and feature suggestions are welcome! Feel free to open an issue or contribute to improve **Gozar**.
