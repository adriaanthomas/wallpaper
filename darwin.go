// +build darwin

package wallpaper

import "os/exec"
import "os/user"
import "path/filepath"
import "strconv"
import "strings"

// Get gets the current wallpaper.
func Get() (wallpaper string, err error) {
	output, err := exec.
		Command("osascript", "-e", `tell application "Finder" to get POSIX path of (get desktop picture as alias)`).
		Output()

	if err != nil {
		return
	}

	// is calling strings.TrimSpace() necessary?
	wallpaper = strings.TrimSpace(string(output))

	return
}

// SetFromFile uses AppleScript to tell Finder to set the desktop wallpaper to specified file.
func SetFromFile(file string) error {
	return exec.
		Command("osascript", "-e", `tell application "Finder" to set desktop picture to POSIX file `+strconv.Quote(file)).
		Run()
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
	usr, err := user.Current()

	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, "Library", "Caches"), nil
}
