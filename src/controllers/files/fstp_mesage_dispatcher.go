package files

import (
	"github.com/pkg/sftp"
	"os"
	"path"
)

func DispatchSftpMessage(messageType int, message []byte, client *sftp.Client) error {
	var fullPath string
	if wd, err := client.Getwd(); err == nil {
		fullPath = path.Join(wd, "/tmp/")
		if _, err := client.Stat(fullPath); err != nil {
			if os.IsNotExist(err) {
				if err := client.Mkdir(fullPath); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	} else {
		return err
	}

	//dstFile, err := client.Create(path.Join(fullPath, header.Filename))
	//if err != nil {
	//	return err
	//}
	//defer srcFile.Close()
	//defer dstFile.Close()
	//
	//_, err = dstFile.ReadFrom(srcFile)
	//if err != nil {
	//	return err
	//}
	return nil
}
