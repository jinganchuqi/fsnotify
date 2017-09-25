package main

import (
	"time"
	"fmt"
	"os/exec"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"golang.org/x/text/encoding/simplifiedchinese"
	"bytes"
	"golang.org/x/text/transform"
	"log"
	"github.com/fsnotify/fsnotify"
	"os"
)

type configType struct {
	rootPath string
	callback string
	except   []interface{}
	depth    int
	reload  string
	isReturn  bool
}

const (
	EVENT_PATH = "./"
	EVENT_CONFIG = "event.json"
)

func initConfig(config *configType) (error) {
	conf, err := ioutil.ReadFile(EVENT_PATH + EVENT_CONFIG)
	js, err := simplejson.NewJson([]byte(string(conf)))
	if err == nil {
		callback, err := js.GetIndex(0).Get("callback").String()
		checkErr(err)
		path, err := js.GetIndex(0).Get("path").String()
		checkErr(err)
		depth, err := js.GetIndex(0).Get("depth").Int()
		checkErr(err)
		except, err := js.GetIndex(0).Get("except").Array()
		checkErr(err)
		reload, err := js.GetIndex(0).Get("reload").String()
		isReturn, err := js.GetIndex(0).Get("isReturn").Bool()
		checkErr(err)
		config.rootPath = path
		config.callback = callback
		config.except = except
		config.depth = depth
		config.reload = reload
		config.isReturn = isReturn
	}
	return err
}

var config configType = configType{}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[E]", r)
		}
	}()
	initConfig(&config)
	//callCmd(config.callback)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println()
				fmt.Println("--------")
				action := event.Op.String()
				if EVENT_CONFIG == string(event.Name) {
					err := initConfig(&config)
					if err == nil {
						fmt.Println("重启中")
						callCmd(config.reload)
					} else {
						checkErr(err)
					}
				} else {
					isExist, _ := inArray(config.except, event.Name)
					if !isExist {
						callbackStr:=config.callback
						if config.isReturn{
							callbackStr +=  " " + action + " " + string(event.Name)
						}
						fmt.Println("notify开始: ---",)
						callCmd(callbackStr)
						fmt.Println("notify结束: ---",)
					}
				}
				switch action {
				case "RENAME":
					watcher.Remove(event.Name)
				case "CREATE", "WRITE":
					fileInfo, err := os.Stat(event.Name)
					if err == nil {
						if fileInfo.IsDir() {
							watcher.Add(event.Name)
							fmt.Println("addWatch:", event.Name)
						}
					} else {
						checkErr(err)
					}
				case "REMOVE":
					watcher.Remove(event.Name)
				}
				fmt.Println("event:", event)
				fmt.Println("--------")
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	fmt.Println("addWatch: ----", config.rootPath)
	fmt.Println("addWatch: ----", EVENT_PATH + EVENT_CONFIG)
	watcher.Add(EVENT_PATH + EVENT_CONFIG)
	walkDir(config.rootPath,"", 0, config, func(dir string,subName string) {
		err = watcher.Add(dir)
		checkErr(err)
	})
	<-done
}

func inArray(array []interface{}, key string) (bool, int) {
	for k, v := range array {
		if v == key {
			return true, k
		}
	}
	return false, 0
}

func walkDir(dirPath string,subName string, depth int, config configType, callback func(dir string,subName string)) {
	if depth > config.depth {
		return
	}
	files, err := ioutil.ReadDir(dirPath)
	if err == nil {
		callback(dirPath,subName)
	} else {
		checkErr(err)
	}
	for _, file := range files {
		isExist,_:=inArray(config.except,file.Name())
		if !isExist {
			if file.IsDir() {
				walkDir(dirPath+"\\"+file.Name(),file.Name(), depth + 1, config, callback)
				continue
			}
		}
	}
}

func callCmd(callback string) {
	cmd := exec.Command("cmd.exe", "/c", callback)
	out, err := cmd.Output()
	cmd.Process.Kill()
	checkErr(err)
	out, err = GbkToUtf8([]byte(out))
	fmt.Println("callback: ----",string(out))
}

func Task(callback func(t time.Time)) {
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for t := range ticker.C {
			callback(t)
		}
	}()
	<-make(chan int)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("err:", err)
	}
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
