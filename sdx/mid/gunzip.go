// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"compress/gzip"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst"
	"net/http"
	"strings"
)

func Gunzip(c *fst.Context) {
	if strings.Contains(c.ReqRaw.Header.Get(cst.HeaderContentEncoding), "gzip") {
		reader, err := gzip.NewReader(c.ReqRaw.Body)
		if err != nil {
			c.AbortDirect(http.StatusBadRequest, "Can't unzip body!!!")
			return
		}
		c.ReqRaw.Body = reader
	}
}
