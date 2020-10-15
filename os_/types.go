// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_

import (
	"syscall"
	"github.com/hacdias/webdav/v3/syscall_"
	//"time"
	//"unsafe"
)

// // Getpagesize returns the underlying system's memory page size.
// func Getpagesize() int { return syscall.Getpagesize() }

// File represents an open file descriptor.
type File struct {
	*file // os specific
}

// HandleOp

func (f *File) GetHandle() syscall.Handle {
	syscall_.Debug()
	return f.pfd.Sysfd
}

func (f *File) GetDebugHandle() syscall.Handle {
	syscall_.Debug()
	return f.pfd.Sysfd
}


func (file *File) FindNextFile(data *syscall.Win32finddata) (err error) {
	syscall_.Debug()
	netname := syscall_.Decompose(file.name)
	return syscall_.FindNextFile_(netname.Server, file.GetHandle(), data)
}

func (file *File) 	Seek_(offset int64, whence int) (newoffset int64, err error) {
	syscall_.Debug()
	netname := syscall_.Decompose(file.name)
	return syscall_.Seek_(netname.Server, file.GetHandle(), offset, whence)
}

func (file *File) ReadFile_(buf []byte, done *uint32, overlapped *syscall.Overlapped) (err error) {
	syscall_.Debug()
	netname := syscall_.Decompose(file.name)
	return syscall_.ReadFile_(netname.Server, file.GetHandle(), buf, done, overlapped)
}

func (file *File) 	TransmitFile_(s syscall.Handle, bytesToWrite uint32, bytsPerSend uint32, overlapped *syscall.Overlapped, transmitFileBuf *syscall.TransmitFileBuffers, flags uint32) (err error) {
	syscall_.Debug()
	//netname := syscall_.Decompose(file.name)
	//return syscall_.TransmitFile_(netname.Server,  s , file.GetHandle(), bytesToWrite, bytsPerSend, overlapped, transmitFileBuf, flags)
	return nil
}

func (file *File) 	CancelIoEx_(o *syscall.Overlapped) (err error) {
	syscall_.Debug()
	netname := syscall_.Decompose(file.name)
	return syscall_.CancelIoEx_(netname.Server,  file.GetHandle(), o)
}

// // A FileInfo describes a file and is returned by Stat and Lstat.
// type FileInfo interface {
// 	Name() string       // base name of the file
// 	Size() int64        // length in bytes for regular files; system-dependent for others
// 	Mode() FileMode     // file mode bits
// 	ModTime() time.Time // modification time
// 	IsDir() bool        // abbreviation for Mode().IsDir()
// 	Sys() interface{}   // underlying data source (can return nil)
// }

// // A FileMode represents a file's mode and permission bits.
// // The bits have the same definition on all systems, so that
// // information about files can be moved from one system
// // to another portably. Not all bits apply to all systems.
// // The only required bit is ModeDir for directories.
// type FileMode uint32

// // The defined file mode bits are the most significant bits of the FileMode.
// // The nine least-significant bits are the standard Unix rwxrwxrwx permissions.
// // The values of these bits should be considered part of the public API and
// // may be used in wire protocols or disk representations: they must not be
// // changed, although new bits might be added.
// const (
// 	// The single letters are the abbreviations
// 	// used by the String method's formatting.
// 	ModeDir        FileMode = 1 << (32 - 1 - iota) // d: is a directory
// 	ModeAppend                                     // a: append-only
// 	ModeExclusive                                  // l: exclusive use
// 	ModeTemporary                                  // T: temporary file; Plan 9 only
// 	ModeSymlink                                    // L: symbolic link
// 	ModeDevice                                     // D: device file
// 	ModeNamedPipe                                  // p: named pipe (FIFO)
// 	ModeSocket                                     // S: Unix domain socket
// 	ModeSetuid                                     // u: setuid
// 	ModeSetgid                                     // g: setgid
// 	ModeCharDevice                                 // c: Unix character device, when ModeDevice is set
// 	ModeSticky                                     // t: sticky
// 	ModeIrregular                                  // ?: non-regular file; nothing else is known about this file

// 	// Mask for the type bits. For regular files, none will be set.
// 	ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice | ModeCharDevice | ModeIrregular

// 	ModePerm FileMode = 0777 // Unix permission bits
// )

// func (m FileMode) String() string {
// 	const str = "dalTLDpSugct?"
// 	var buf [32]byte // Mode is uint32.
// 	w := 0
// 	for i, c := range str {
// 		if m&(1<<uint(32-1-i)) != 0 {
// 			buf[w] = byte(c)
// 			w++
// 		}
// 	}
// 	if w == 0 {
// 		buf[w] = '-'
// 		w++
// 	}
// 	const rwx = "rwxrwxrwx"
// 	for i, c := range rwx {
// 		if m&(1<<uint(9-1-i)) != 0 {
// 			buf[w] = byte(c)
// 		} else {
// 			buf[w] = '-'
// 		}
// 		w++
// 	}
// 	return string(buf[:w])
// }

// // IsDir reports whether m describes a directory.
// // That is, it tests for the ModeDir bit being set in m.
// func (m FileMode) IsDir() bool {
// 	return m&ModeDir != 0
// }

// // IsRegular reports whether m describes a regular file.
// // That is, it tests that no mode type bits are set.
// func (m FileMode) IsRegular() bool {
// 	return m&ModeType == 0
// }

// // Perm returns the Unix permission bits in m.
// func (m FileMode) Perm() FileMode {
// 	return m & ModePerm
// }

func (fs *fileStat) Name() string { return fs.name }
func (fs *fileStat) IsDir() bool  { return fs.Mode().IsDir() }

// // SameFile reports whether fi1 and fi2 describe the same file.
// // For example, on Unix this means that the device and inode fields
// // of the two underlying structures are identical; on other systems
// // the decision may be based on the path names.
// // SameFile only applies to results returned by this package's Stat.
// // It returns false in other cases.
// func SameFile(fi1, fi2 FileInfo) bool {
// 	fs1, ok1 := fi1.(*fileStat)
// 	fs2, ok2 := fi2.(*fileStat)
// 	if !ok1 || !ok2 {
// 		return false
// 	}
// 	return sameFile(fs1, fs2)
// }
