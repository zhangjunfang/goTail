package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var (
	c           = make(chan os.Signal)
	once        sync.Once
	fs          []string
	fileAndSeek map[string]int64
	shut        []*os.File
)

func initData() {
	signal.Notify(c, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	fs = make([]string, 16)
	shut = make([]*os.File, 16)
	once = sync.Once{}
	once.Do(initData)
	fileAndSeek = make(map[string]int64)
	go SignalTest()
}

func initFile(v string, f *os.File) (*leveldb.DB, error) {
	name := f.Name()
	v = strings.Replace(v, name, strings.Replace(name, ".", "", -1), -1) + "/d"
	db, err := leveldb.OpenFile(v, nil)
	if err != nil {
		return nil, err
	}
	val, err := db.Get([]byte(v), nil)
	if err == errors.ErrNotFound {
		f.Seek(0, 0)
		return db, nil
	} else {
		db.Close()
		return nil, err
	}
	t, err := strconv.Atoi(string(val))
	if err != nil {
		db.Close()
		return nil, err
	}
	f.Seek(int64(t), 0)
	return db, nil
}
func SignalTest() {
	select {
	case <-c:
		{
			fmt.Println("非正常终止程序！！！！！")
			for _, v := range shut {
				v.Close()
			}
			close(c)
		}
	}
}
func ReadFileData() {
	for _, v := range fs {
		v = strings.Replace(v, "\\", "/", -1)
		go ReadFile(v)
	}
	time.Sleep(2 * time.Second)
}
func ReadFile(v string) {
	f, err := os.Open(v) //defer f.Close()
	if err != nil {
		return
	}
	shut = append(shut, f)
	var temp int64
	db, err := initFile(v, f)
	fmt.Println("--87---", db, err)
	if err != nil {
	}
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		temp = temp + int64(n)
		if err != nil {
			db.Put([]byte(v), []byte(strconv.Itoa(int(temp))), nil)
			runtime.Gosched()
			//			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Println("----99---1--", temp, n, string(buf[:n]))
	}
}

func FileIterator(dir string) {
	filepath.Walk(dir, WalkFunc)
}
func WalkFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && strings.HasSuffix(path, "/d") {
		fs = append(fs, path)
	}
	return nil
}

type WriteSource func(s string)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	dir := "d:/test" //需要从命令行读取监控路径
	FileIterator(dir)
	ReadFileData()
	for {

	}
}
