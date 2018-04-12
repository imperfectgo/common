// +build !commonlog_noautoinit

package log

import "flag"

func init() {
	if err := AddFlags(flag.CommandLine); err != nil {
		panic(err)
	}
}
