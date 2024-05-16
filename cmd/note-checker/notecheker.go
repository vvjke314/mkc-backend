package main

import (
	"fmt"

	"github.com/vvjke314/mkc-backend/internal/notechecker"
)

func main() {
	nc := notechecker.NewNoteChecker()
	err := nc.Init()
	if err != nil {
		fmt.Println("[notecheker]: can't initialize notechecker")
	}

	err = nc.Run()
	if err != nil {
		msg := fmt.Sprintf("[notecheker]: failed run %s", err.Error())
		nc.Log(msg)
	}
}
