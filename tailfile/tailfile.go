package tailfile

import (
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

var TailObj *tail.Tail
var err error

func Init(filename string) {
	config := tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: false,
		Poll:      true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
	}

	TailObj, err = tail.TailFile(filename, config)

	if err != nil {
		logrus.Error("tailfile: create tailObj for path:%s failed, err:%v", filename, err)
		return
	}
	return
}
