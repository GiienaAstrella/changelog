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

## Why NodeJS

Many of my projects and many DDV projects utilize NodeJS.
Creating a wrapper around the Go tool is a no-brainer.

While I could write this tool in Typescript, that would require installing NodeJS (or some other JavsScript environment) for projects that may not need it (e.g. pure Go projects).
The most obvious solution here is to write the tool in Go and write a JS wrapper around the tool.

## Installation

There are various ways to install Changelog.

### Manual

You can download the pre-compiled binaries for supported platforms from the [release] page.
You can also download such binaries from the official distribution server: `https://projects.gassets.space/changelog/VERSION/changelog-PLATFORM-VERSION.tar.gz`, where `VERSION` is the version number (i.e. `0.1.0`) and `PLATFORM` is the platform identifier pair (i.e. `darwin-amd64`).

Simply extract this archive.

### NPM

On systems with NPM installed, you can install Changelog with NPM.

``` text
npm install @ghifari160/changelog
```

**Note:** on older versions, the installation will *silently fail* if installed with
`--ignore-scripts`.
These older versions rely on installation hooks to install the Changelog binary.
As of Changelog v0.3.0, it is safe to use `--ignore-scripts`.
The binary installation happens on first execution instead.

If installed locally in a project, you can run the tool through NPM

``` text
npx changelog command [options] [arguments...]
```

### From source

On systems with Go installed, run the tool install from source.

``` text
go install giiena.me/changelog@latest
```

[Keep a Changelog]: https://keepachangelog.com/en/1.1.0
[DDV]: https://github.com/DiamondDrakeVentures
[release]: https://github.com/GiienaAstrella/changelog/releases/latest
