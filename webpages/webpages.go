package webpages

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
)

const (
	basePath string = "/home/leron/workspace/src/github.com/dimitarvalkanov7/temp/results{v}.json"
)

var (
	mutex sync.Mutex
	data  = make([]PageContent, 20)
)

type PageContent struct {
	URL   string         `json:"URL"`
	Words map[string]int `json:"Words"`
}

func GetData() []PageContent {
	numOfFiles := 3
	var wg sync.WaitGroup
	for i := 1; i <= numOfFiles; i++ {
		var version string
		if i < 10 {
			version = "0" + strconv.Itoa(i)
		} else {
			version = strconv.Itoa(i)
		}
		path := basePath
		fullPath := strings.Replace(path, "{v}", version, -1)
		wg.Add(1)
		go func(fullPath string) {
			loadScrappedDataInMemory(fullPath)
			defer wg.Done()
		}(fullPath)
	}

	wg.Wait()
	return data
}

func loadScrappedDataInMemory(fullPath string) {
	mutex.Lock()
	file, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Println(err)
		return
	}

	temp := make([]PageContent, 0)
	_ = json.Unmarshal([]byte(file), &temp)
	data = append(data, temp...)
	mutex.Unlock()
}
