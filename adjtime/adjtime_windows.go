package adjtime

import (
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	FileTime1970 = 0x019db1ded53e8000
)

var (
	kerneldll32              = windows.NewLazySystemDLL("kernel32.dll")
	procFileTimeToSystemTime = kerneldll32.NewProc("FileTimeToSystemTime")
	procSetSystemTime        = kerneldll32.NewProc("SetSystemTime")
)

func errnoErr(e syscall.Errno) error {
	if e == 0 {
		return syscall.EINVAL
	}
	return e
}

func adjtime(delta time.Duration) (err error) {
	// Ref: https://github.com/ntp-project/ntp/blob/9c75327c3796ff59ac648478cd4da8b205bceb77/ports/winnt/libntp/SetSystemTime.c#L29
	var currTime windows.Filetime
	windows.GetSystemTimeAsFileTime(&currTime)
	currTimeULL := (*uint64)(unsafe.Pointer(&currTime))
	*currTimeULL += 10 * uint64(delta.Microseconds())
	var sysTime windows.Systemtime
	if err = fileTimeToSystemTime(&currTime, &sysTime); err != nil {
		return err
	}
	return setSystemTime(&sysTime)
}

func fileTimeToSystemTime(fileTime *windows.Filetime, systemTime *windows.Systemtime) (err error) {
	r1, _, er := syscall.SyscallN(procFileTimeToSystemTime.Addr(), uintptr(unsafe.Pointer(fileTime)), uintptr(unsafe.Pointer(systemTime)))
	if r1 == 0 {
		err = errnoErr(er)
	}
	return err
}

func setSystemTime(time *windows.Systemtime) (err error) {
	r1, _, er := syscall.SyscallN(procSetSystemTime.Addr(), uintptr(unsafe.Pointer(time)))
	if r1 == 0 {
		err = errnoErr(er)
	}
	return err
}
