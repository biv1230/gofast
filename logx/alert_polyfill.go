//go:build !linux
// +build !linux

package logx

func Report(string) {
}

func SetReporter(func(string)) {
}
