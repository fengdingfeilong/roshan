package roshantool

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func loginfo(s string, err error) {
	if InnerLog != nil {
		InnerLog("roshan: "+s, err)
	}
}

//GetFileMD5 get file md5
func GetFileMD5(path string) string {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		loginfo("handler.GetFileList.getFileMD5", err)
		return ""
	}
	md5 := md5.New()
	_, err = io.Copy(md5, file)
	if err != nil {
		loginfo("handler.GetFileList.getFileMD5", err)
		return ""
	}
	s := hex.EncodeToString(md5.Sum(nil))
	return s
}
