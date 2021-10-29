package filecompare

import (
	"fmt"
	"io/fs"
	"time"

	log "github.com/sirupsen/logrus"
)

type CfgAttribute struct {
	Name string
}

type Configuration struct {
	Attributes []CfgAttribute `yaml:"attributes"`
}

func Differs(logger *log.Entry, config *Configuration, srcFile, dstFile fs.FileInfo) (bool, string, error) {
	logger = logger.WithFields(log.Fields{"Component": "FileCompare", "Package": "filecompare", "Function": "Compare", "File": srcFile.Name()})
	for _, attribute := range config.Attributes {
		switch attribute.Name {
		case "size":
			if srcFile.Size() != dstFile.Size() {
				logger.Warnf("the file size has chnaged: %d != %d", srcFile.Size(), dstFile.Size())
				return true, attribute.Name, nil
			}
		case "modify":
			if srcFile.ModTime().After(dstFile.ModTime().Add(time.Second)) { // TODO - evalulate if his comparision can give the wrong answer due to the second added
				logger.Warnf("the source file is newer: %s != %s", srcFile.ModTime().Format(time.RFC3339), dstFile.ModTime().Format(time.RFC3339))
				return true, attribute.Name, nil
			}
		default:
			return false, "", fmt.Errorf("do not know how to compare file attribute: %s", attribute)
		}
	}
	return false, "", nil
}
