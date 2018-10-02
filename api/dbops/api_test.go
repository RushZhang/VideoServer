package dbops

import (
	"testing"
	"strconv"
	"time"
	"fmt"
)

/*
	涉及到数据库的单元测试：
	1、init（包括db的login，truncate table）
	2、运行测试
	3、清除数据（truncate table）
 */


var tempvid string


func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

//因为要保证执行顺序，所以要用到子test
func TestUserWorkFlow(t *testing.T) {
	t.Run("ADD", testAddUser)
	t.Run("GET", testGetUser)
	t.Run("DEL", testDeleteUser)
	t.Run("REGET", testRegetUser)

}

/*
	==================User========================
 */
func testAddUser(t *testing.T) {
	err := AddUserCredential("rush", "rush123")
	if err != nil {
		t.Errorf("ErrorChan of AddUser: %v", err)
	}
}

func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("rush")
	if err != nil {
		t.Errorf("ErrorChan of GetUser: %v", err)
	}
	if pwd != "rush123" {
		t.Errorf("得到的数据错误")
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("rush", "rush123")
	if err != nil {
		t.Errorf("ErrorChan of DeleteUser: %v", err)
	}
}

//reget是为了测试删除后删的到底彻底不
func testRegetUser(t *testing.T) {
	pwd, err := GetUserCredential("rush")
	if err != nil {
		t.Errorf("ERROR of RegetUser: %v", err)
	}
	if pwd != "" {
		t.Errorf("delete user test failed")
	}
}




/*
	=====================Video=======================
 */


func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("PrepareUser", testAddUser)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DelVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}


func testAddVideoInfo(t *testing.T) {
	vi, err := AddNewVideo(1, "my-video")
	if err != nil {
		t.Errorf("ErrorChan of AddVideoInfo: %v", err)
	}
	tempvid = vi.Id
}

func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(tempvid)
	//fmt.Println(video)
	if err != nil {
		t.Errorf("ErrorChan of GetVideoInfo: %v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tempvid)
	if err != nil {
		t.Errorf("ErrorChan of DeleteVideoInfo: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	vi, err := GetVideoInfo(tempvid)
	if err != nil || vi != nil{
		t.Errorf("ErrorChan of RegetVideoInfo: %v", err)
	}
}



/*
	==============Test===================
 */
func TestComments(t *testing.T) {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddCommnets", testAddComments)
	t.Run("ListComments", testListComments)
}

func testAddComments(t *testing.T) {
	vid := "12345"
	aid := 1
	content := "I like this video"

	err := AddNewComments(vid, aid, content)

	if err != nil {
		t.Errorf("ErrorChan of AddComments: %v", err)
	}

	vid = "12000"
	aid = 1
	content = "Name?"

	err = AddNewComments(vid, aid, content)

	if err != nil {
		t.Errorf("ErrorChan of AddComments: %v", err)
	}
}

func testListComments(t *testing.T) {
	vid := "12345"
	from := 151476480
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))

	res, err := ListComments(vid, from, to)
	if err != nil {
		t.Errorf("ErrorChan of ListComments: %v", err)
	}

	for i, ele := range res {
		fmt.Printf("comment: %d, %v \n", i, ele)
	}
}




