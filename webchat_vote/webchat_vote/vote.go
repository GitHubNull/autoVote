package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"compress/gzip"
	"encoding/binary"

	simplejson "github.com/bitly/go-simplejson"
)

//var voicecode = "1291528"
var RootPath = "/root/webchat_vote/"
var ConfigPath = RootPath + "config/"
var voicecodeFile = ConfigPath +  "voicecode.txt"
var openidsFile = ConfigPath + "openids.txt"
var logFile = RootPath + "vote_log.txt"

func voicecodeReader(fileName string) (string, error) {
	file, err := os.Open(fileName)
	codeStr := ""
	if nil != err {
		fmt.Println("L29", err.Error())
		return codeStr, err
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	line, err := buf.ReadString('\n')
	if "" != line {
		line = strings.TrimSpace(line)
		//fmt.Println(line)dd
		codeStr = line
	}
	if nil != err {
		fmt.Println("L",err.Error())
	}

	return codeStr, nil

}

func openidsReader(fileName string) (*list.List, error) {
	file, err := os.Open(fileName)
	openidsList := list.New()
	if nil != err {
		fmt.Println("L52", err.Error())
		return openidsList, err
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if "" != line {
			line = strings.TrimSpace(line)
			openidsList.PushBack(line)
		}
		//fmt.Println(line)
		if nil != err {
			if io.EOF == err {
				//fmt.Println("quit")
				break
			} else {
				fmt.Println("L",err.Error())
				return openidsList, err
			}
		}

	}
	return openidsList, nil
}

func logVoteInfo(info string) error {
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 06666)
	if nil != err {
		fmt.Println("L81", err.Error())
		return err
	}
	defer file.Close()
	file.WriteString(info + "\n")
	return nil
}

func vote(openid string, voicecode string) {
	client := &http.Client{}
	data := "openid=" + openid + "&voicecode=" + voicecode
	ContentLength := strconv.Itoa(len(data))
	req, err := http.NewRequest("POST", "https://wx.qingxuanwenhua.com/weixinmp/vote.php", strings.NewReader(data))
	if err != nil {
		// handle error
	}

	req.Header.Set("charset", "utf-8")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("referer", "https://servicewechat.com/wx91914c13ad712c6e/50/page-frame.html")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.0; MI 5 Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044033 Mobile Safari/537.36 MicroMessenger/6.6.6.1300(0x26060636) NetType/WIFI Language/zh_CN MicroMessenger/6.6.6.1300(0x26060636) NetType/WIFI Language/zh_CN")
	req.Header.Set("Content-Length", ContentLength)
	req.Header.Set("Host", "wx.qingxuanwenhua.com")
	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)
	if nil != err {
		fmt.Println("L109",err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("L116",err.Error())
		return
	}
	b := new(bytes.Buffer)
    	binary.Write(b, binary.LittleEndian, body)
	r, err := gzip.NewReader(b)
	if err != nil {
		// handle error
		fmt.Println("L123",err.Error())
		return
	}
	defer r.Close()
	undata, err := ioutil.ReadAll(r)
	if err != nil {
		// handle error
		fmt.Println("L132",err.Error())
		return
	}
	body = undata
	fmt.Println(string(body))
	fmt.Println((body))
	js, err := simplejson.NewJson(body)
	if nil != err {
		fmt.Println("L123",err.Error())
		return
	}
	check, err := js.Get("check").Int()
	if nil != err {
		fmt.Println("L128",err.Error())
		return
	}
	status, err := js.Get("status").Bool()
	if nil != err {
		fmt.Println("L133",err.Error())
		return
	}
	if 0 == check || false == status {
		info := time.Now().Format("2006-01-02 15:04:05") + " [" + openid + "] vote for [" + voicecode + "] " + "error."
		fmt.Println(info)
		err = logVoteInfo(info)
		if nil != err {
			fmt.Println("L141",err.Error())
		}

	} else {
		info := time.Now().Format("2006-01-02 15:04:05") + " [" + openid + "] vote for [" + voicecode + "] " + "success."
		err = logVoteInfo(info)
		fmt.Println(info)
		if nil != err {
			fmt.Println("L149",err.Error())
		}
	}
}

func main() {
	voicecodeStr, err := voicecodeReader(voicecodeFile)
	if nil != err {
		fmt.Println("L157",err.Error())
		return
	}

	openidsList, err := openidsReader(openidsFile)
	if nil != err {
		fmt.Println("L163",err.Error())
		return
	}
	for e := openidsList.Front(); nil != e; e = e.Next() {
		fmt.Println(e.Value)
		vote(e.Value.(string), voicecodeStr)
	}
}
