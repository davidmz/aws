package aws

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Keys struct {
	AccessKey string
	SecretKey string
}

// Загрузка из заданного файла
func LoadKeysFromFile(filePath string) (k *Keys) {
	k = &Keys{}
	if f, e := os.Open(filePath); e == nil {
		scanner := bufio.NewScanner(f)
		if scanner.Scan() {
			k.AccessKey = strings.TrimSpace(scanner.Text())
		}
		if scanner.Scan() {
			k.SecretKey = strings.TrimSpace(scanner.Text())
		}
		f.Close()
	}
	if k.AccessKey == "" || k.SecretKey == "" {
		k = nil
	}
	return
}

// Поиск реквизитов для соединения
// 1. Переменные окружения
// 2. файл .awssecret в текущем каталоге
// 3. файл .awssecret в домашнем каталоге
func LoadKeys() (k *Keys) {
	k = &Keys{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
	}
	if k.AccessKey == "" {
		k = nil
		dirs := []string{"."}

		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			homeDir = os.Getenv("HOMEPATH") // Windows
		}
		if homeDir != "" {
			dirs = append(dirs, homeDir)
		}
		for _, dir := range dirs {
			if d, e := filepath.Abs(dir); e == nil {
				k = LoadKeysFromFile(filepath.Join(d, ".awssecret"))
				if k != nil {
					break
				}
			}
		}
	}
	return
}
