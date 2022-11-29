# Scale HTTP Adapters

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-brightgreen.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Tests](https://github.com/loopholelabs/scale-http-adapters/actions/workflows/test.yml/badge.svg)](https://github.com/loopholelabs/scale-http-adapters/actions/workflows/test.yml)

This library contains the definitions and source for the official Scale HTTP Adapters for [Scale Functions](https://scale.sh). These 
adapters are used to integrate Scale Functions with the following HTTP servers (organized by language):

- [Go](https://golang.org)
  - [net/http](https://pkg.go.dev/net/http)
  - [fasthttp](https://pkg.go.dev/github.com/valyala/fasthttp)
  - [fiber](https://pkg.go.dev/github.com/gofiber/fiber/v2)

**This library requires Go1.18 or later.**

## Important note about releases and stability

This repository generally follows [Semantic Versioning](https://semver.org/). However, **this library is currently in
Beta** and is still considered experimental. Breaking changes of the library will _not_ trigger a new major release. The
same is true for selected other new features explicitly marked as
**EXPERIMENTAL** in [the changelog](/CHANGELOG.md).

## Usage and Documentation

Usage instructions and documentation for using the Scale HTTP Adapters is available at [https://scale.sh/docs](https://scale.sh/docs).

## Contributing

Bug reports and pull requests are welcome on GitHub at [https://github.com/loopholelabs/scale-http-adapters][gitrepo]. For more
contribution information check
out [the contribution guide](https://github.com/loopholelabs/scale-http-adapters/blob/master/CONTRIBUTING.md).

## License

The Scale HTTP Adapters project is available as open source under the terms of
the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).

## Code of Conduct

Everyone interacting in the Scale HTTP Adapters project’s codebases, issue trackers, chat rooms and mailing lists is expected to follow the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).

## Project Managed By:

[![https://loopholelabs.io][loopholelabs]](https://loopholelabs.io)

[gitrepo]: https://github.com/loopholelabs/scale-http-adapters
[loopholelabs]: https://cdn.loopholelabs.io/loopholelabs/LoopholeLabsLogo.svg
[loophomepage]: https://loopholelabs.io
