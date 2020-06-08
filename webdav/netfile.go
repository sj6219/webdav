// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webdav

import (
	"path/filepath"
	"strings"
)

// slashClean is equivalent to but slightly more efficient than
// path.Clean("/" + name).
type NetFile struct {
	Server string
	Path   string
}

func ResolvePath(dir, name string) NetFile {
	// This implementation is based on Dir.Open's code in the standard net/http package.
	i := strings.IndexByte(name, '/')
	var prefix string
	if i >= 0 {
		prefix = name[:i]
	} else {
		prefix = name
	}

	i = strings.IndexByte(prefix, '@')
	if i < 0 {
		return NetFile{Server: "", Path: filepath.Join(dir, filepath.FromSlash(slashClean(name)))}
	}

	return NetFile{Server: prefix[:i], Path: name[i+1:]}
}

func (p *NetFile) String() string {
	if len(p.Server) == 0 {
		return p.Path
	}
	if len(p.Path) == 0 {
		return "."
	}
	i := strings.IndexByte(p.Path, '@')
	if i < 0 {
		return "\\\\" + p.Server + "\\" + p.Path
	}
	return "\\\\" + p.Path[:i] + "\\" + p.Path[i+1:]
}
