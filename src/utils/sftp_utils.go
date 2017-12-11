package utils

import (
	"os"
	"path"
	"mime/multipart"
	"github.com/pkg/sftp"
	//"github.com/genshen/webConsole/src/models"
)

type SftpNode Node // struct alias.

func UploadFile(user SftpNode, username, password string ,srcFile multipart.File, header *multipart.FileHeader) error {
	sshEntity := SSH{
		Node: Node{
			Host: user.Host,
			Port: user.Port,
		},
	}
	_, err := sshEntity.Connect(username, password)
	if err != nil {
		return err
	} else {
		defer sshEntity.Close()

		client, err := sftp.NewClient(sshEntity.Client)
		if err != nil {
			return err
		}

		var fullPath string;
		if wd, err := client.Getwd(); err == nil {
			fullPath = path.Join(wd, "/tmp/")
			if _, err := client.Stat(fullPath); err != nil {
				if (os.IsNotExist(err)) {
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

		dstFile, err := client.Create(path.Join(fullPath, header.Filename))
		if err != nil {
			return err
		}
		defer srcFile.Close()
		defer dstFile.Close()

		_, err = dstFile.ReadFrom(srcFile)
		if err != nil {
			return err
		}
		return nil
	}
}
