package main

import (
	"encoding/json"
	"fmt"
	gonadlan "github.com/devbrain/go-nadlan"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func iterateFiles() {

	_ = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".json") {
			fmt.Println(path)
			f, ferr := os.Open(path)
			if ferr != nil {
				return ferr
			}
			defer f.Close()
			body, ferr := io.ReadAll(f)
			if ferr != nil {
				return ferr
			}
			var yad2Data gonadlan.Yad2RawData
			jerr := json.Unmarshal(body, &yad2Data)
			if jerr != nil {
				return jerr
			}
			for _, fi := range yad2Data.Data.Feed.FeedItems {
				for _, obj := range fi.Row4 {
					fmt.Println(obj)
				}
			}
			gonadlan.ParseYad2RawData(&yad2Data, true)
		}
		return nil
	})

}

func main() {
	iterateFiles()
	//page := 1
	//lastPage := 0
	//homes := mapset.NewSet[string]()
	//for {
	//	fmt.Println(page)
	//	_, lp, err := gonadlan.GetYad2Data(page, 9000, false)
	//	if err != nil {
	//		fmt.Printf("Error %v\n", err)
	//		break
	//	}
	//
	//	if lastPage == 0 {
	//		lastPage = lp
	//	}
	//	page++
	//	if page > lastPage {
	//		break
	//	}
	//}
	//fmt.Println(homes)
}
