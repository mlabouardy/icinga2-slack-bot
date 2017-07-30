## Icinga2 Slack Bot [![CircleCI](https://circleci.com/gh/mlabouardy/icinga2-slack-bot/tree/master.svg?style=svg)](https://circleci.com/gh/mlabouardy/icinga2-slack-bot/tree/master) [![Go Report Card](https://goreportcard.com/badge/github.com/mlabouardy/icinga2-slack-bot)](https://goreportcard.com/report/github.com/mlabouardy/icinga2-slack-bot) [![Gitter chat](https://badges.gitter.im/icinga2bot/Lobby.png)](https://gitter.im/icinga2bot/Lobby) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

This bot uses Icinga2 remote API to fetch the status of the services & hosts running in icinga2

## Configuration

- Setup a slack bot token
- Setup icinga2 instance credentials in config.toml

## Available commands

![alt text](https://raw.githubusercontent.com/mlabouardy/icinga2-slack-bot/master/screenshots/help.png)

### Get all hosts

![alt text](https://raw.githubusercontent.com/mlabouardy/icinga2-slack-bot/master/screenshots/hosts.png)

### Filter by host name

![alt text](https://raw.githubusercontent.com/mlabouardy/icinga2-slack-bot/master/screenshots/host.png)

### Get all services

![alt text](https://raw.githubusercontent.com/mlabouardy/icinga2-slack-bot/master/screenshots/services.png)

### Filter by service name

![alt text](https://raw.githubusercontent.com/mlabouardy/icinga2-slack-bot/master/screenshots/service.png)
