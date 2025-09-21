package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

func main() {
	filename := "./xx.log"
	config := tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: false,
		Poll:      true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
	}

	tails, err := tail.TailFile(filename, config)

	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		msg *tail.Line
		ok  bool
	)

	for {
		msg, ok = <-tails.Lines
		if !ok {
			fmt.Println("tailTest err")
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Println(msg.Text)
	}
}
