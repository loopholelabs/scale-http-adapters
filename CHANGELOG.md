# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.3.8] - 2023-03-19

### Dependencies

- Bumping `scale` version to `v0.3.15`
- Bumping `scale-signature-http` version to `v0.3.7`

## [v0.3.7] - 2023-03-12

### Features

- Fixing a bug in the adapters where the header cases would not be properly normalized to lower case

### Dependencies

- Bumping `scale` version to `v0.3.14`
- Bumping `scale-signature` version to `v0.2.11`
- Bumping `scale-signature-http` version to `v0.3.7`
- Bumping `scalefile` version to `v0.1.9`

## [v0.3.6] - 2023-02-28

### Features

- Adding a `WithScale` helper to the `NextJS` adapter which enables the webpack loader to be used with `NextJS`

### Dependencies

- Bumping Scale version to `v0.3.12`

## [v0.3.5] - 2023-02-20

### Changes

- Bumping Scale version to `v0.3.11`

## [v0.3.4] - 2023-02-20

### Changes

- Bumping Scale version to `v0.3.10`

## [v0.3.3] - 2023-02-20

### Changes

- Bumping Scale version to `v0.3.9`

## [v0.3.2] - 2023-02-19

### Changes

- Bumping Scale version to `v0.3.8`

## [v0.3.1] - 2023-02-19

### Changes

- Bumping Scale version to `v0.3.7` which fixes a bug in the `Go` runtime where passing in `nil` as the `Next` function would cause a panic

## [v0.3.0] - 2023-02-18

### Changes

- Bumping Scale version to `v0.3.6`
- Bumping Scale HTTP Signature version to `v0.3.4`
- Bumping Scalefile version to `v0.1.7`
- Updating NextJS Adapter to work with new Scalefile format

## [v0.2.1] - 2023-02-01

### Fixes

- Fixing npm publish github action to properly publish package

## [v0.2.0] - 2023-02-01

### Changes

- Bumping Scale runtime version to `v0.2.1`
- Bumping Scale HTTP Signature version to `v0.2.3`
- Bumping Scalefile version to `v0.1.4`

### Features

- Adding NextJS Adapter

## [v0.1.0] - 2022-11-29

### Features

- Initial release of the Scale HTTP Adapters library.

[unreleased]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.3.8...HEAD
[v0.3.8]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.3.8
[v0.3.7]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.3.7
[v0.3.6]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.3.6
[v0.3.5]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.3.5
[v0.3.4]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.3.4
[v0.3.3]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.3.3
[v0.3.2]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.3.2
[v0.3.1]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.3.1
[v0.3.0]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.3.0
[v0.2.1]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.2.1
[v0.2.0]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.2.0
[v0.1.0]: https://github.com/loopholelabs/scale-http-adapters/compare/v0.1.0
