# Changelog

A simple tool to manage [Keep a Changelog] changelogs.

``` text
USAGE:
    changelog command [command options] [arguments...]

COMMANDS:
    get      show changes for a specific version(s)
    prepare  prepare changelog for next cycle
    promote  promote unreleased
    version  print the app version
    help, h  Shows a list of commands or help for one command
```

## Why

Communicating changes between versions in a project is important.
There are various ways to do this and various formats.
With multiple projects, this can easily grow out of hands.

Personally, I follow the Keep a Changelog format.
At [DDV], we also use the Keep a Changelog format.

## Why Go

My projects and DDV projects are written in various languages, using various tech stacks, on various platforms.
To avoid rewriting this tool as hacky scripts for every single project (been there, done that), and to keep the dependency list minimal, whatever language I use must compile into native binaries.
Many languages can achieve this goal, but I am most comfortable with Go.

## Installation

There are various ways to install Changelog.

### Manual

You can download the pre-compiled binaries for supported platforms from the [release] page.

Simply extract this archive.

### From source

On systems with Go installed, run the tool install from source.

``` text
go install giiena.me/changelog@latest
```

[Keep a Changelog]: https://keepachangelog.com/en/1.1.0
[DDV]: https://github.com/DiamondDrakeVentures
[release]: https://github.com/GiienaAstrella/changelog/releases/latest
