package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func GetFileMd5(filePath string) string {
	file, inerr := os.Open(filePath)
	if inerr == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		file.Close()
		return hex.EncodeToString(md5h.Sum(nil))
	}
	return ""
}
func IsFileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)

	// _, err := os.Stat(name)
	// if os.IsNotExist(err) {
	// 	return false, nil
	// }
	// return err != nil, err
}
