// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
//go:build !linux

package breaker

func Report(string) {
}

func SetReporter(func(string)) {
}
