package watchlog

import (
	"awesomeFilebeat/etcd"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

var taskMgr *watchLogMgr

type watchLogMgr struct {
	newConfChan       chan []*etcd.LogEntry
	watchYamlFilePath string
}

func Init(yamlFilePath string) {
	taskMgr = &watchLogMgr{
		newConfChan:       make(chan []*etcd.LogEntry), // 无缓冲区的通道
		watchYamlFilePath: yamlFilePath,
	}
	go taskMgr.run()
}

func (t *watchLogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			fmt.Println(newConf)
			t.rewriteLogConfig(newConf)
		}
	}
}

func (t *watchLogMgr) rewriteLogConfig(newConf []*etcd.LogEntry) {
	// 读取 YAML 文件内容
	data, err := ioutil.ReadFile(t.watchYamlFilePath)
	if err != nil {
		panic(err)
	}

	// 定义一个 Config 类型的变量，用于存储解析后的结果
	var c Configs
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}

	// 清空切片
	c.Filebeat.Inputs = []Inputs{}

	for _, conf := range newConf {
		begin := strings.LastIndex(conf.Path, "/")
		end := strings.LastIndex(conf.Path, ".")
		types := conf.Path[begin : end+1]
		log_types := fmt.Sprintf("%s-%s", conf.Topic, types)

		input := Inputs{Type: "log", Enabled: true, Paths: []string{conf.Path},
			Fields: Fields{LogTopics: conf.Topic, LogTypes: log_types}}

		c.Filebeat.Inputs = append(c.Filebeat.Inputs, input)
	}

	// 写回源文件
	yamlData, err := yaml.Marshal(c)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	// 将转换后的 YAML 内容写入文件
	err = ioutil.WriteFile(t.watchYamlFilePath, yamlData, 0644)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
}

// NewConfChan 向外暴露 watchLogMgr 的 newConfChan
func NewConfChan() chan<- []*etcd.LogEntry {
	return taskMgr.newConfChan
}
