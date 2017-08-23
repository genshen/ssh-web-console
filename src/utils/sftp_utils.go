package utils

import (
	"github.com/pkg/sftp"
	"mime/multipart"
	"github.com/genshen/webConsole/src/models"
	"os"
	"path"
)

func UploadFile(user models.UserInfo, srcFile multipart.File, header *multipart.FileHeader) error {
	sshEntity := SSH{
		Node: Node{
			Host: user.Host,
			Port: user.Port,
		},
	}
	_, err := sshEntity.Connect(user.Username, user.Password)
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
