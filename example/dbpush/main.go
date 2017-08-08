package main

import "sync"
import "strings"
import "github.com/go-pg/pg"
import . "github.com/Vonng/gopher/db/pg"
import log "github.com/Sirupsen/logrus"

type User struct {
	ID   string `sql:",pk"`
	Name string
}

var Users sync.Map // Users 内部数据缓存

func LoadAllUser() {
	var users []User
	Pg.Query(&users, `SELECT ID,name FROM users;`)
	for _, user := range users {
		Users.Store(user.ID, user)
	}
}

func LoadUser(id string) {
	user := User{ID: id}
	Pg.Select(&user)
	Users.Store(user.ID, user)
}

func PrintUsers() string {
	var buf []string
	Users.Range(func(key, value interface{}) bool {
		buf = append(buf, key.(string));
		return true
	})
	return strings.Join(buf, ",")
}

// ListenUserChange 会监听PostgreSQL users数据表中的变动通知
func ListenUserChange() {
	go func(c <-chan *pg.Notification) {
		for notify := range c {
			action, id := notify.Payload[0], notify.Payload[1:]
			switch action {
			case 'I':
				fallthrough
			case 'U':
				LoadUser(id);
			case 'D':
				Users.Delete(id)
			}
			log.Infof("[NOTIFY] Action:%c ID:%s Users: %s", action, id, PrintUsers())
		}
	}(Pg.Listen("users_chan").Channel())
}

// MakeSomeChange 会向数据库写入一些变更
func MakeSomeChange() {
	go func() {
		Pg.Insert(&User{"001", "张三"})
		Pg.Insert(&User{"002", "李四"})
		Pg.Insert(&User{"003", "王五"})  // 插入
		Pg.Update(&User{"003", "王麻子"}) // 改名
		Pg.Delete(&User{ID: "002"})    // 删除
	}()
}

func main() {
	Pg = NewPg("postgres://localhost:5432/postgres")
	Pg.Exec(`TRUNCATE TABLE users;`)
	LoadAllUser()
	ListenUserChange()
	MakeSomeChange()
	<-make(chan struct{})
}
