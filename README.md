# GoZix Redigo

[documentation-img]: https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff
[documentation-url]: https://pkg.go.dev/github.com/gozix/redigo/v4
[license-img]: https://img.shields.io/github/license/gozix/redigo.svg?style=for-the-badge
[license-url]: https://github.com/gozix/redigo/blob/master/LICENSE
[release-img]: https://img.shields.io/github/tag/gozix/redigo.svg?label=release&color=24B898&logo=github&style=for-the-badge
[release-url]: https://github.com/gozix/redigo/releases/latest
[build-status-img]: https://img.shields.io/github/actions/workflow/status/gozix/redigo/go.yml?logo=github&style=for-the-badge
[build-status-url]: https://github.com/gozix/redigo/actions
[go-report-img]: https://img.shields.io/badge/go%20report-A%2B-green?style=for-the-badge
[go-report-url]: https://goreportcard.com/report/github.com/gozix/redigo
[code-coverage-img]: https://img.shields.io/codecov/c/github/gozix/redigo.svg?style=for-the-badge&logo=codecov
[code-coverage-url]: https://codecov.io/gh/gozix/redigo

[![License][license-img]][license-url]
[![Documentation][documentation-img]][documentation-url]

[![Release][release-img]][release-url]
[![Build Status][build-status-img]][build-status-url]
[![Go Report Card][go-report-img]][go-report-url]
[![Code Coverage][code-coverage-img]][code-coverage-url]

The bundle provide a Redigo integration to GoZix application.

## Installation

```shell
go get github.com/gozix/redigo/v4
```

## Dependencies

* [viper](https://github.com/gozix/viper)

## Configuration example

```json
{
  "redis": {
    "default": {
      "host": "127.0.0.1",
      "port": "6379",
      "db": 0,
      "password": "somepassword",
      "max_idle": 3,
      "max_active": 100,
      "idle_timeout": "4m"
    }
  }
}
```

"password" field is optional and ignored if empty
"db" field is optional. Default is 0

## Documentation

You can find documentation on [pkg.go.dev][documentation-url] and read source code if needed.

## Questions

If you have any questions, feel free to create an issue.
