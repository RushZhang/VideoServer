package session

import (
	"sync"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/utils"
	"time"
)


/*
	处理session的策略是先写到内存的sessionMap，再写到数据库
 */




//sync.Map是线程安全的Map，普通Map是线程不安全
//这个简单的session来模拟缓存
var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}


func nowInMilli() int64 {
	return time.Now().UnixNano()/1000000
}


//把db的session加载到内存
func LoadSessionsFromDB() {
	tempMap, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}

	tempMap.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})

}

func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	curtime := nowInMilli()
	ttl := curtime + 30 * 60 * 1000//Serverside session valid time: 这里定义30分钟
	ss := &defs.SimpleSession{
		Username: un,
		TTL: ttl,
	}
	//先写缓存
	sessionMap.Store(id, ss)
	//再写数据库
	dbops.InsertSession(id, ttl, un)

	return id
}



func IsSessionExpired(sid string) (string, bool) {
	tempSession, ok := sessionMap.Load(sid)
	ct := nowInMilli()
	if ok {
		if tempSession.(*defs.SimpleSession).TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}
		return tempSession.(*defs.SimpleSession).Username, false
	} else {
		//在有loadbalance的情况下，有可能不同节点操作了map后内存会不一致
		//所以如果读不出来，不一定真的不存在，可以强制去读数据库
		ss, err := dbops.RetrieveSession(sid)
		if err != nil || ss == nil {
			return "", true
		}
		if ss.TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}
		sessionMap.Store(sid, ss)
		return ss.Username, false
	}

	return "", true

}



func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}