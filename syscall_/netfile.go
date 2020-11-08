// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscall_

import (
	"syscall"
	"strings"
	// "path"
	// "path/filepath"
	// "strings"
	// "syscall"
	"unsafe"
)

func Debug() {
}

type HandleOp interface {
	GetHandle() syscall.Handle
	GetDebugHandle() syscall.Handle
	GetFileType() (uint32, error)
	CloseHandle_() (err error)
	FindClose_() (err error)
	FindNextFile(data *syscall.Win32finddata) (err error)
	Seek_(offset int64, whence int) (newoffset int64, err error)
	ReadFile_(buf []byte, done *uint32, overlapped *syscall.Overlapped) (err error)
	TransmitFile_(s syscall.Handle, bytesToWrite uint32, bytsPerSend uint32, overlapped *syscall.Overlapped, transmitFileBuf *syscall.TransmitFileBuffers, flags uint32) (err error)
	CancelIoEx_(o *syscall.Overlapped) (err error)
}

// func (h Handle) GetHandle() Handle {
// 	return h
// }
// func (h Handle) GetDebugHandle() Handle {
// 	return h
// }
// func (h Handle) GetFileType() (uint32, error) {
// 	return GetFileType(h)
// }

var dll syscall.Handle

type NetFileName struct {
	Server string
	Path   string
}

func LoadLib(name string) {
	dll, _ = LoadLibrary(name)
	//syscall.LoadLib(name)
}

func GetProc(name string) uintptr {
	f, _ := GetProcAddress(dll, name)	
	return f
}

func DecomposeFromPtr(p *uint16) NetFileName {
	var i int
	for i = 0; ; i++ {
		ch := *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(i*2)))
		if ch == 0 {
			break
		}
	}
	slice := make([]uint16, i+1)
	for i = 0; ; i++ {
		ch := *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(i*2)))
		slice[i] = ch
		if ch == 0 {
			break
		}
	}
	netname := Decompose(syscall.UTF16ToString(slice))
	if netname.Server == "" {
		Debug()
	}
	return netname
}

func Decompose(name string) NetFileName {
	// This implementation is based on Dir.Open's code in the standard net/http package.
	i := 0
	if strings.HasPrefix(name, "\\\\") {
		i = 2
	}
	j := strings.IndexByte(name[i:], '\\')
	if j < 0 {
		j = len(name) - i
	}
	prefix := name[i : i+j]

	k := strings.IndexByte(prefix, '#')
	if k < 0 {
		return NetFileName{Server: "", Path: name}
	}

	return NetFileName{Server: name[i : i+k], Path: name[i+k+1:]}
}

func DebugDecompose(name string) NetFileName {
	netname := Decompose(name)
	if netname.Server != "" {
		Debug()
	}
	return netname
}

func (name NetFileName) DebugString() string {
	if len(name.Server) == 0 {
		return string(name.Path)
	}
	if len(name.Path) == 0 {
		return "."
	}
	i := strings.IndexByte(name.Path, '/')
	if i < 0 {
		i = len(name.Path)
	}
	prefix := name.Path[:i]
	j := strings.IndexByte(prefix, '#')
	if j < 0 {
		return "\\\\" + name.Server + "\\" + name.Path
	}
	return "\\\\" + name.Path[:j] + "\\" + name.Path[j+1:]
}

func (name NetFileName) String() string {
	if len(name.Server) == 0 {
		return string(name.Path)
	}
	return "\\\\" + name.Server + "#" + name.Path
}
