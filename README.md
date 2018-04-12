# Common

[![GoDoc](https://godoc.org/github.com/imperfectgo/common?status.svg)](https://godoc.org/github.com/imperfectgo/common) 
[![Build Status](https://travis-ci.org/imperfectgo/common.svg?branch=master)](https://travis-ci.org/imperfectgo/common)
[![Go Report Card](https://goreportcard.com/badge/github.com/imperfectgo/common)](https://goreportcard.com/report/github.com/imperfectgo/common)
[![Coverage](https://codecov.io/gh/imperfectgo/common/branch/master/graph/badge.svg)](https://codecov.io/gh/imperfectgo/common)

- [log](./log): A logging wrapper around [zap](https://github.com/uber-go/zap).
- [version](./version): Version information and metrics.
- [debugutil](.debugutil): Utils for debugging purpose.

## Caveats

When using this package together with `github.com/prometheus/common/log`, you may encounter panic on program initializing.
In order to workaround this issue, please just add `commonlog_noautoinit` to your build tags.
