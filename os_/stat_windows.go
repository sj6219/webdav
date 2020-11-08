// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_

import (
	"os"
	"github.com/hacdias/webdav/v3/internal_/syscall_/windows_"
	"syscall"
	"github.com/hacdias/webdav/v3/syscall_"
	"unsafe"
)

// Stat returns the FileInfo structure describing file.
// If there is an error, it will be of type *PathError.
func (file *File) Stat() (os.FileInfo, error) {
	if file == nil {
		return nil, os.ErrInvalid
	}

	if file.isdir() {
		// I don't know any better way to do that for directory
		return Stat(file.dirinfo.path)
	}
	if isWindowsNulName(file.name) {
		return &devNullStat, nil
	}

	ft, err := file.pfd.GetFileType()
	if err != nil {
		return nil, &os.PathError{"GetFileType", file.name, err}
	}
	switch ft {
	case syscall.FILE_TYPE_PIPE, syscall.FILE_TYPE_CHAR:
		return &fileStat{name: basename(file.name), filetype: ft}, nil
	}

	fs, err := newFileStatFromGetFileInformationByHandle(file.name, file.pfd.Sysfd)
	if err != nil {
		return nil, err
	}
	fs.filetype = ft
	return fs, err
}

// stat implements both Stat and Lstat of a file.
func stat(funcname, name string, createFileAttrs uint32) (os.FileInfo, error) {
	netname := syscall_.Decompose(name)
	if netname.Server == "" {
		syscall_.Debug()
	}

	if len(name) == 0 {
		return nil, &os.PathError{funcname, name, syscall.Errno(syscall.ERROR_PATH_NOT_FOUND)}
	}
	if isWindowsNulName(name) {
		return &devNullStat, nil
	}
	namep, err := syscall.UTF16PtrFromString(fixLongPath(netname.String()))
	if err != nil {
		return nil, &os.PathError{funcname, name, err}
	}

	// Try GetFileAttributesEx first, because it is faster than CreateFile.
	// See https://golang.org/issues/19922#issuecomment-300031421 for details.
	var fa syscall.Win32FileAttributeData
	err = syscall_.GetFileAttributesEx(namep, syscall.GetFileExInfoStandard, (*byte)(unsafe.Pointer(&fa)))
	if err == nil && fa.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT == 0 {
		// Not a symlink.
		fs := &fileStat{
			FileAttributes: fa.FileAttributes,
			CreationTime:   fa.CreationTime,
			LastAccessTime: fa.LastAccessTime,
			LastWriteTime:  fa.LastWriteTime,
			FileSizeHigh:   fa.FileSizeHigh,
			FileSizeLow:    fa.FileSizeLow,
		}
		if err := fs.saveInfoFromPath(netname.String()); err != nil {
			return nil, err
		}
		return fs, nil
	}
	// GetFileAttributesEx fails with ERROR_SHARING_VIOLATION error for
	// files, like c:\pagefile.sys. Use FindFirstFile for such files.
	if err == windows_.ERROR_SHARING_VIOLATION {
		var fd syscall.Win32finddata
		err := syscall_.FindFirstFile_Close(namep, &fd)
		if err != nil {
			return nil, &os.PathError{"FindFirstFile", name, err}
		}
		fs := newFileStatFromWin32finddata(&fd)
		if err := fs.saveInfoFromPath(netname.String()); err != nil {
			return nil, err
		}
		return fs, nil
	}

	// Finally use CreateFile.
	if netname.Server == "" {
		h, err := syscall.CreateFile(namep, 0, 0, nil,
			syscall.OPEN_EXISTING, createFileAttrs, 0)
		if err != nil {
			return nil, &os.PathError{"CreateFile", name, err}
		}
		defer syscall.CloseHandle(h)

		return newFileStatFromGetFileInformationByHandle(name, h)
	} else {
		var d syscall.ByHandleFileInformation
		var reparseTag uint32

		_, _, err := syscall.Syscall6(syscall_.GetProc("_GetFileInformation"), 4, uintptr(unsafe.Pointer(namep)), uintptr(createFileAttrs), uintptr(unsafe.Pointer(&d)), uintptr(unsafe.Pointer(&reparseTag)), 0, 0)
		if err != 0 {
			return nil, &os.PathError{"_GetFileInformation", name, syscall.Errno(err)}
		}
		return &fileStat{
			name:           basename(name),
			FileAttributes: d.FileAttributes,
			CreationTime:   d.CreationTime,
			LastAccessTime: d.LastAccessTime,
			LastWriteTime:  d.LastWriteTime,
			FileSizeHigh:   d.FileSizeHigh,
			FileSizeLow:    d.FileSizeLow,
			vol:            d.VolumeSerialNumber,
			idxhi:          d.FileIndexHigh,
			idxlo:          d.FileIndexLow,
			Reserved0:      reparseTag,
			// fileStat.path is used by os.SameFile to decide if it needs
			// to fetch vol, idxhi and idxlo. But these are already set,
			// so set fileStat.path to "" to prevent os.SameFile doing it again.
		}, nil
	}
}

// statNolog implements Stat for Windows.
func statNolog(name string) (os.FileInfo, error) {
	return stat("Stat", name, syscall.FILE_FLAG_BACKUP_SEMANTICS)
}

// lstatNolog implements Lstat for Windows.
func lstatNolog(name string) (os.FileInfo, error) {
	attrs := uint32(syscall.FILE_FLAG_BACKUP_SEMANTICS)
	// Use FILE_FLAG_OPEN_REPARSE_POINT, otherwise CreateFile will follow symlink.
	// See https://docs.microsoft.com/en-us/windows/desktop/FileIO/symbolic-link-effects-on-file-systems-functions#createfile-and-createfiletransacted
	attrs |= syscall.FILE_FLAG_OPEN_REPARSE_POINT
	return stat("Lstat", name, attrs)
}
