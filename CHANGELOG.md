<!-- markdownlint-disable MD024 -->

# Changelog

All notable changes in Changelog will be documented in this file.

The format is based on [Keep a Changelog].
This project attempts to adhere to [Semantic Versioning].

## [UNRELEASED]

### Added

- Shell completion ([#13]).
- Global flag `--old-parser` (`-x`) to disable the new experimental parser and use the old text
  parser implementation.
- `--single-line` (`-s`) flag for `get` command.
  With this flag set, `get` will output changes for specific version(s) without breaking
  sentences in a paragraph into their own lines ([#3]).
- New experimental `Parse` function in package `keepachangelog`.
  This parses `CHANGELOG.md` using Abstract Syntax Tree (AST) parsing.
- Types in package `keepachangelog` now implements `String` function, which will generate the
  appropriate Markdown string for that type and its contents.
- Types `Version` and `Section` in package `keepachangelog` now implements `SingleLineString`
  functions, which will generate the appropriate Markdown string for that type on its contents.
  Unlike `String`, `SingleLineString` does not break sentences into multiple lines ([#3]).

### Changed

- Changelog is now installable and importable under `giiena.me/changelog`.
  Past versions are (prior to v0.4.0) must be installed and imported from
  `github.com/ghifari160/changelog`.
- Changelog is no longer available as an NPM package. Distribution through NPM required wrapper
  script, which is now broken.
  We do not have the bandwidth to fix it.
  NPM distribution may return in the future, but for now the supported installation methods going
  forward are as outlined in [README.md](https://github.com/GiienaAstrella/changelog#installation).
- All commands now use the new experimental AST parser.
- All commands now properly support reference style links
  (`[Text]`, `[Text][key]`, and `[Text][]`).
  `get` will intelligently output references used in the version body.
  `promote` now optionally accepts link and title for the version page, which will be added to the
  references.
  It will also preserves utilized references.
  `prepare` will properly preserves utilized references ([#4]).
- `urfave/cli` library has been upgraded to `v3`.

### Deprecated

- Implementations of `markdown.Unmarshaler` interface in all types in package `keepachangelog`.
- Implementations of `markdown.Marshaler` interface in all types in package `keepachangelog`.
- Package `markdown`.
  For marshaling, the type should implement [`fmt.Stringer`](https://pkg.go.dev/fmt#Stringer).
  For parsing, refer to the target type's parsing function.

### Removed

### Fixed

### Security

## [0.3.2] - 2025-08-07

### Fixed

- Fixed an issue where sublists are not preserved on operations ([#2]).

## [0.3.1] - 2025-07-25

### Fixed

- Fixed a bug where Changelog binary would be installed relative to the current working directory.
  The installation location is now relative to the module directory in `node_modules`.

## [0.3.0] - 2025-07-22

### Changed

- When installing from NPM, Changelog precompiled binary is now installed to `vendor/changelog/bin/changelog` (`vendor\changelog\bin\changelog.exe` on Windows).
- When installing from NPM, installation of the Changelog binary is now delayed until first execution.
  This way, we're no longer relying on installation hooks (which can be skipped by the user) to properly install the Changelog binary.
- `get` and `promote` commands now strips leading `v` from version targets.
  Running `changelog get v0.2.0` will match `0.2.0`, and running `changelog promote v0.3.0` will create a version section with `0.3.0` as the heading.

## [0.2.0] - 2025-04-08

### Changed

- Default export for NodeJS wrapper is now the wrapper object itself.
  See `src/wrapper.ts` for more information.
- NodeJS wrapper pre-installation hook now imports the default export for NodeJS wrapper module.
- When installing on unsupported platforms through NPM, the wrapper will attempt to build from source archive.
  Note: building from source requires [Go](https://go.dev).

### Fixed

- When installing through NPM on Windows, the pre-installation hook now downloads the correct archive.
- NodeJS wrapper no longer prints the path to temporary files when downloading archives.
- Command `prepare` now prints all default and user provided sections to the changelog file.
  Previously, a changelog file without an existing `unreleased` version causes the command to output an empty version.

## [0.1.0] - 2025-03-26

### Added

- Package `markdown`.
  It implements encoding and decoding of data from and to Markdown formatted representation.
- Package `keepachangelog`.
  It implements types and functions to assist in maintaining a Changelog based on the [Keep a Changelog] format.
- Command `get`, which shows changes for a specific version(s).
- Command `version`, which prints the app version.
- Command `promote`, which promotes unreleased draft to be the next release version.
- Command `prepare`, which prepares the changelog for the next release cycle.
- NodeJS wrapper.
  Changelog can be installed through npm (`npm install @ghifari160/changelog`).
  On supported platforms, the pre-install hook download and install the precompiled binary for that platform.
  It can also be imported as a module, which will return the path to the changelog binary.
  Note: installation will silently fail of installed with `--ignore-scripts`.

[#13]: https://github.com/GiienaAstrella/changelog/issues/13
[#2]: https://github.com/GiienaAstrella/changelog/issues/2
[#3]: https://github.com/GiienaAstrella/changelog/issues/3
[#4]: https://github.com/GiienaAstrella/changelog/issues/4
[0.1.0]: https://github.com/GiienaAstrella/changelog/releases/tag/0.1.0
[0.2.0]: https://github.com/GiienaAstrella/changelog/releases/tag/0.2.0
[0.3.0]: https://github.com/GiienaAstrella/changelog/releases/tag/0.3.0
[0.3.1]: https://github.com/GiienaAstrella/changelog/releases/tag/0.3.1
[0.3.2]: https://github.com/GiienaAstrella/changelog/releases/tag/0.3.2
[Keep a Changelog]: https://keepachangelog.com/en/1.1.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
[UNRELEASED]: https://github.com/GiienaAstrella/changelog/compare/0.3.2...HEAD
