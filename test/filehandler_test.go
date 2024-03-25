package main

import (
	"testing"

	"github.com/google/uuid"
	"github.com/vvjke314/mkc-backend/internal/pkg/filehandler"
)

func TestFilehandler(t *testing.T) {
	projName := uuid.New().String()
	filehandler.CreateDir(projName)

	filehandler.CreateFile("25ec5ab9-27df-42b2-b28e-5bd0137df668/hello.txt", []byte{'a', 'b', 'c', 'f'})
	//filehandler.RemoveFile("25ec5ab9-27df-42b2-b28e-5bd0137df668/hello.txt")
}
