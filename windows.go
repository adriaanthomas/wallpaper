// +build windows

package wallpaper

import (
	"os"

	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

// Get gets the current wallpaper.
func Get() (wallpaper string, err error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.READ)

	if err != nil {
		return
	}

	defer func() {
		err = key.Close()
	}()

	wallpaper, _, err = key.GetStringValue("Wallpaper")

	if err != nil {
		return
	}

	return
}

// SetFromFile sets the wallpaper for the current user to specified file by setting HKEY_CURRENT_USER\Control Panel\Desktop\Wallpaper.
//
// Note: this requires you to log out and in again.
func SetFromFile(file string) (err error) {
	systemParametersInfo(spiSetDeskWallPaper, 0, file, spifUpdateIniFile|spifSendWinIniChange)

	return
}

// SetFromURL downloads url and calls SetFromFile.
func SetFromURL(url string) error {
	file, err := downloadImage(url)

	if err != nil {
		return err
	}

	return SetFromFile(file)
}

func getCacheDir() (string, error) {
	return os.TempDir(), nil
}

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	systemParametersInfoProc = user32.NewProc("SystemParametersInfoW")
)

const (
	spiSetDeskWallPaper = 0x0014

	spifUpdateIniFile    = 0x01
	spifSendWinIniChange = 0x02
)

func systemParametersInfo(uiAction uint, uiParam uint, pvParam string, fWinIni uint) (err error) {
	r1, r2, err := systemParametersInfoProc.Call(
		uintptr(uiAction),
		uintptr(uiParam),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(pvParam))),
		uintptr(fWinIni))
	fmt.Printf("%v, %v, %v\n", r1, r2, err)
	return
}
