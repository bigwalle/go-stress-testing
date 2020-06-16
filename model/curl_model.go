/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-19
* Time: 09:51
 */

package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// curl参数解析
type CURL struct {
	Data map[string][]string
}

func (c *CURL) getDataValue(keys []string) []string {
	var (
		value = make([]string, 0)
	)

	for _, key := range keys {
		var (
			ok bool
		)

		value, ok = c.Data[key]
		if ok {
			break
		}
	}

	return value
}

// 从文件中解析curl
func ParseTheFile(path string) (curl *CURL, err error) {

	if path == "" {
		err = errors.New("路径不能为空")

		return
	}

	curl = &CURL{
		Data: make(map[string][]string),
	}

	file, err := os.Open(path)
	if err != nil {
		err = errors.New("打开文件失败:" + err.Error())

		return
	}

	defer func() {
		file.Close()
	}()

	dataBytes, err := ioutil.ReadAll(file)
	if err != nil {
		err = errors.New("读取文件失败:" + err.Error())

		return
	}
	data := string(dataBytes)

	for len(data) > 0 {
		if strings.HasPrefix(data, "curl") {
			data = data[5:]
		}

		data = strings.TrimSpace(data)
		var (
			key   string
			value string
		)

		index := strings.Index(data, " ")
		if index <= 0 {
			break
		}
		key = strings.TrimSpace(data[:index])
		data = data[index+1:]
		data = strings.TrimSpace(data)

		// url
		if !strings.HasPrefix(key, "-") {
			key = strings.Trim(key, "'")
			curl.Data["curl"] = []string{key}

			// 去除首尾空格
			data = strings.TrimFunc(data, func(r rune) bool {
				if r == ' ' || r == '\\' || r == '\n' {
					return true
				}

				return false
			})
			continue
		}

		if strings.HasPrefix(data, "-") {
			continue
		}

		var (
			endSymbol = " "
		)

		if strings.HasPrefix(data, "'") {
			endSymbol = "'"
			data = data[1:]
		}

		index = strings.Index(data, endSymbol)
		if index <= -1 {
			break
		}
		value = data[:index]
		data = data[index+1:]

		// 去除首尾空格
		data = strings.TrimFunc(data, func(r rune) bool {
			if r == ' ' || r == '\\' || r == '\n' {
				return true
			}

			return false
		})

		curl.Data[key] = append(curl.Data[key], value)

		// break

	}

	// for key, value := range curl.Data {
	// 	fmt.Println("key:", key, "value:", value)
	// }

	return
}

func (c *CURL) String() (url string) {
	curlByte, _ := json.Marshal(c)

	return string(curlByte)
}

// GetUrl
func (c *CURL) GetUrl() (url string) {

	keys := []string{"curl"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {

		return
	}

	url = value[0]

	return
}

// GetMethod
func (c *CURL) GetMethod() (method string) {
	method = "GET"

	var (
		postKeys = []string{"--d", "--data", "--data-binary $", "--data-binary"}
	)
	value := c.getDataValue(postKeys)

	if len(value) >= 1 {
		return "POST"
	}

	keys := []string{"-X", "--request"}
	value = c.getDataValue(keys)

	if len(value) <= 0 {

		return
	}

	method = strings.ToUpper(value[0])

	return
}

// GetHeaders
func (c *CURL) GetHeaders() (headers map[string]string) {
	headers = make(map[string]string, 0)

	keys := []string{"-H", "--header"}
	value := c.getDataValue(keys)

	for _, v := range value {
		getHeaderValue(v, headers)
	}

	return
}

// GetHeaders
func (c *CURL) GetHeadersStr() string {
	headers := c.GetHeaders()
	bytes, _ := json.Marshal(&headers)

	return string(bytes)
}

// 获取body
func (c *CURL) GetBody() (body string) {

	keys := []string{"--data", "-d", "--data-raw", "--data-binary"}
	value := c.getDataValue(keys)

	if len(value) <= 0 {

		return
	}

	// body = strings.NewReader(value[0])
	body = value[0]

	return
}
