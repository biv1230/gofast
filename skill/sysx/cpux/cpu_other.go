//go:build !linux
// +build !linux

package cpux

func RefreshCpu() uint64 {
	return 0
}
