package dialog

import (
	"os"
	"path/filepath"
	"syscall"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func driveMask() uint32 {
	dll, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		fyne.LogError("Error loading kernel32.dll", err)
		return 0
	}
	handle, err := syscall.GetProcAddress(dll, "GetLogicalDrives")
	if err != nil {
		fyne.LogError("Could not find GetLogicalDrives call", err)
		return 0
	}

	ret, _, err := syscall.Syscall(uintptr(handle), 0, 0, 0, 0)
	if err != syscall.Errno(0) { // for some reason Syscall returns something not nil on success
		fyne.LogError("Error calling GetLogicalDrives", err)
		return 0
	}

	return uint32(ret)
}

func listDrives() []string {
	var drives []string
	mask := driveMask()

	for i := 0; i < 24; i++ {
		if mask&1 == 1 {
			letter := string('A' + i)
			drives = append(drives, letter+":")
		}
		mask >>= 1
	}

	return drives
}

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
	var places []fyne.CanvasObject

	for _, drive := range listDrives() {
		places = append(places, widget.NewButton(drive, func() {
			f.setDirectory(drive + string(os.PathSeparator))
		}))
	}
	return places
}

func isHidden(file os.FileInfo, dir string) bool {
	path := filepath.Join(dir, file.Name())

	point, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		fyne.LogError("Error making string pointer", err)
		return false
	}
	attr, err := syscall.GetFileAttributes(point)
	if err != nil {
		fyne.LogError("Error getting file attributes", err)
		return false
	}

	return attr&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}