//go:build !windows

package file

func BitmapFilename(filename string) string {
	return filename + ".bitmap"
}

func SetHidden(path string) error {
	return nil
}
