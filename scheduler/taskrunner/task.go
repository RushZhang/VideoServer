package taskrunner

import (
	"video_server/scheduler/dbops"
	"log"
	"github.com/pkg/errors"
	"sync"
	"video_server/scheduler/ossops"
)

/*
	Runner是生产消费者模型，可复用
	task是只针对这个项目，不可复用
 */


//func deleteVideo(vid string) error {
//	log.Println("delete函数：", VIDEO_PATH + vid)
//	err := os.Remove(VIDEO_PATH + vid)
//	if err != nil || !os.IsNotExist(err) {
//		log.Printf("Deleting video error: %v", err)
//		return err
//	}
//	return nil
//}

func deleteVideo(vid string) error {
	ossfn := "videos/" + vid
	bn := "rush-videos2"  //bucket name
	ok := ossops.DeleteObject(ossfn, bn)

	if !ok {
		log.Printf("Deleting video error, oss operation failed")
		return errors.New("Deleting video error")
	}

	return nil
}


func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil {
		//log.Printf("Video clear dispatcher err: %v", err)
		return err
	}
	if len(res) == 0 {
		return errors.New("All tasks finished, 没记录可以再删除了")
	}

	for _, id := range res {
		dc <- id
	}
	return nil
}


func VideoClearExecutor(dc dataChan) error {
	//因为有多个线程在删，所以每个线程得到的错误先记录下来，统一返回
	errMap := &sync.Map{}
	var err error

	forloop:
		for {
			select {
			case vid := <- dc:
			/*
				这里使用go func()会有点风险。比如
				dispatcher发出了: 1, 2, 3这三个要删除。 executor也接受到了1，2，3，
				但是新建的go routine只在数据库删了1时executor就发了READYTODISPATCH的消息
				这样dispatcher搜索数据库又得到了2，3，4. 所以2，3要被重复删除
				这个可以通过sync.groupwait来解决，但是有点复杂，没必要做
			 */
				go func(id interface{}) {
					//log.Println("要删除的id：", id)
					//这里把vid当做参数传进go routine就是闭包的讲究了
					if err := deleteVideo(id.(string)); err != nil {
						log.Println(err)
						errMap.Store(id, err)
						return
					}
					if err := dbops.DelVideoDeletionRecord(id.(string)); err != nil {
						errMap.Store(id, err)
						return
					}
				}(vid)

			default:
				break forloop

			}
		}


	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			return false
		}
		return true
	})

	//只要随便返回一个错误就代表整个过程不顺利
	return err


}

