# Common

- [log](./log): A logging wrapper around [zap](https://github.com/uber-go/zap).
- [version](./version): Version information and metrics.
- [debugutil](.debugutil): Utils for debugging purpose.

## Caveats

When using this package together with `github.com/prometheus/common/log`, you may encounter panic on program initializing.
In order to workaround this issue, please just add `commonlog_noautoinit` to your build tags.
