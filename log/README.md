# log package

Currently, this documentation is unfinished, so you may need to read the source code before use.

## Flags

The log package comes with following flags auto registered by default.

- `log.level`
- `log.format`

**HINT** You can disable the auto-register feature by passing `commonlog_noautoinit` to your build tags.

### `log.level`

Only log messages with the given severity or above. Valid levels:

- debug
- info
- warn
- error
- fatal

### `log.format`

The `log.format` have a common format (optional parts marked by squared brackets):

`logger:<encoder>[?param=value[&param2=value2]]`

A full example:

`logger:json?outputPaths=/var/log/test.log&disableCaller=false&disableStacktrace=false`

#### Encoders

Currently, only two encoders are supported:

1. JSON (Default)
2. Console

#### Parameters

- `development`
- `disableCaller`
- `disableStacktrace`
- `outputPaths`
- `lumberjack`
