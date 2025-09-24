package tailfile

import (
	"LogCollector/common"
	"github.com/sirupsen/logrus"
)

type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask       //所有tailTask任务
	collectEntryList []common.CollectEntry      //所有配置项
	confChan         chan []common.CollectEntry //等待新配置的通道
}

var (
	ttMgr *tailTaskMgr
)

func Init(allConf []common.CollectEntry) (err error) {
	//allConf中存了若干个日志的收集项
	//针对每一个日志收集项创建一个对应的tailObj

	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask),
		confChan:         make(chan []common.CollectEntry),
		collectEntryList: allConf,
	}
	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic)
		err = tt.Init()
		if err != nil {
			logrus.Errorf("create tailObj for path:%s failed,err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create tailObj for path:%s success", conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt
		go tt.run()
	}

	go ttMgr.watch() //在后台等新的配置来

	return
}

func (t *tailTaskMgr) watch() {
	for {
		newConf := <-t.confChan
		logrus.Info("get new conf from etcd,conf:%v", newConf)
		for _, conf := range newConf {
			if t.isExist(conf) {
				continue
			}
			tt := newTailTask(conf.Path, conf.Topic)
			err := tt.Init()
			if err != nil {
				logrus.Errorf("create tailObj for path:%s failed,err:%v", conf.Path, err)
				continue
			}
			logrus.Infof("create tailObj for path:%s success", conf.Path)
			ttMgr.tailTaskMap[tt.path] = tt
			go tt.run()
		}
		//原来有现在没有的tailTask要停掉
		for key, task := range t.tailTaskMap {
			var found bool
			for _, conf := range t.collectEntryList {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				//没找到这个task，停掉
				logrus.Infof("the task collect path:%s need to stop", task.path)
				delete(t.tailTaskMap, key)
				task.cancel()
			}
		}
	}
}

func (t *tailTaskMgr) isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Topic]
	return ok
}

func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
