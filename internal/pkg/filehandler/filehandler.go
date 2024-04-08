package filehandler

import (
	"fmt"
	"os"
)

const Path = "/home/vvjkee/Workspace/university/diplom/storage/"

// type Filehandler struct {
// }

// func NewFileHandler() *Filehandler {
// 	return &Filehandler{}
// }

// CreateDir Создает директорию
func CreateDir(dirName string) error {
	err := os.Mkdir(Path+dirName, 0755) // Создание папки с правами доступа 0755
	if err != nil {
		return fmt.Errorf("[filehandler.CreateDir]: can't create directory %w", err)
	}

	return err
}

// CreateFile Создает файл в локальном хранилище
func CreateFile(fileName string, content []byte) error {
	file, err := os.Create(Path + fileName)
	if err != nil {
		return fmt.Errorf("[os.Create]: can't create file %w", err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("[file.Write]: can't write content to file.Writer: %w", err)
	}
	defer file.Close()

	return err
}

// RemoveFile
// Удаляет файл из локального хранилища
func RemoveFile(fileName string) error {
	err := os.Remove(Path + fileName) // Удаление файла
	if err != nil {
		return fmt.Errorf("[os.Remove]: can't remove file: %w", err)
	}

	return err
}
