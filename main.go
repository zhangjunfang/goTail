// DDDDD project main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

var c = make(chan os.Signal)
var temp int64 = 0
var dir, _ = os.Getwd()
var db, _ = leveldb.OpenFile(strings.Replace(dir, "\\", "/", -1)+"/.d", nil)
var once sync.Once
var key string
var f *os.File

func initData() {
	signal.Notify(c, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	initFile()
}
func initFile() {
	v, _ := db.Get(key, nil)
	t, _ := strconv.Atoi(string(v))
	temp = int64(t)
	f.Seek(temp, 0)
}
func SignalTest() {
	select {
	case <-c:
		{
			db.Put(key, []byte(strconv.Itoa(int(temp))), nil)
		}
	}

}
func main() {
	f, _ = os.Open("test.txt")
	buf := make([]byte, 32)
	key = []byte(strings.Replace(dir, "\\", "/", -1) + "/" + f.Name())
	once.Do(initData)
	go SignalTest()
	for {
		n, err := f.Read(buf)
		temp = temp + int64(n)
		if err != nil {
			db.Put(key, []byte(strconv.Itoa(int(temp))), nil)
			initFile()
			time.Sleep(2 * time.Second)
			continue
		}
		fmt.Println(string(buf[0:n]))
		fmt.Println(temp)
	}
}
