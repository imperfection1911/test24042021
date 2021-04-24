package main

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"sort"
)

var BackupRoot = "/tmp/backup"

func main() {
	project := flag.String("p", "", "project name")
	flag.Parse()
	projectPath := BuildPath(*project)
	err := CheckDirExist(projectPath)
	if err != nil {
		log.Fatal(err)
	}
	err = RemoveBackups(projectPath)
	if err != nil {
		log.Fatal(err)
	}
}

func BuildPath(project string) (projectPath string) {
	projectPath = path.Join(BackupRoot, project)
	return
}

func CheckDirExist(projectPath string) (err error) {
	info, err := os.Stat(projectPath)
	if err != nil {
		return
	}
	if !info.IsDir() {
		err = errors.New("project is not a directory")
	}
	return
}

func RemoveBackups(projectPath string) (err error) {
	var backups []os.FileInfo
	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		matched, err := filepath.Match("*.tar", filepath.Base(path))
		if err != nil {
			return err
		} else if matched {
			backups = append(backups, info)
		}
		return nil
	})
	if err != nil {
		return
	}
	if len(backups) <= 5 {
		log.Info("nothing to do")
		return
	}
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].ModTime().After(backups[j].ModTime())
	})
	for _, i := range backups[5:] {
		err = os.Remove(path.Join(projectPath, i.Name()))
		if err != nil {
			return
		}
		log.Info(fmt.Sprintf("%s removed", i.Name()))
	}
	return
}
