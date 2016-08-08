

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/laktak/hjson-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func getContent(file string) ([]byte) {
	if data, err := ioutil.ReadFile(file); err == nil {
		return data
	} else { panic(err) }
}

func getTestContent(name string) ([]byte) {
	p := fmt.Sprintf("./assets/%s_test.hjson", name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		p = fmt.Sprintf("./assets/%s_test.json", name)
	}
	return getContent(p)
}

func getResultContent(name string) ([]byte,[]byte) {
	p1 := fmt.Sprintf("./assets/sorted/%s_result.json", name)
	p2 := fmt.Sprintf("./assets/sorted/%s_result.hjson", name)
	return getContent(p1), getContent(p2)
}

func fixJson(data []byte) ([]byte) {
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
	data = bytes.Replace(data, []byte("\\u0008"), []byte("\\b"), -1)
	data = bytes.Replace(data, []byte("\\u000c"), []byte("\\f"), -1)
	return data
}

func run(t *testing.T, file string) {
	name := strings.TrimSuffix(file, "_test" + filepath.Ext(file))
	t.Logf("running %s", name);
	shouldFail := strings.HasPrefix(file, "fail")

	testContent := getTestContent(name)
	var data interface{}
	if err := hjson.Unmarshal(testContent, &data); err != nil {
		if !shouldFail { panic(err) } else { return }
	} else if shouldFail { panic(errors.New(name + " should_fail!")) }

	rjson, rhjson := getResultContent(name)

	actualHjson, _ := hjson.Marshal(data)
	actualJson, _ := json.MarshalIndent(data, "", "  ")
	actualJson = fixJson(actualJson)

	// add fixes where go's json differs from javascript
	switch name {
	case "kan":
		actualJson = []byte(strings.Replace(string(actualJson), "    -0,", "    0,", -1))
	case "pass1":
		actualJson = []byte(strings.Replace(string(actualJson), "1.23456789e+09", "1234567890", -1))
	}

	hjsonOK := bytes.Equal(rhjson, actualHjson)
	jsonOK := bytes.Equal(rjson, actualJson)
	if !hjsonOK {
		t.Logf("%s\n---hjson expected\n%s\n---hjson actual\n%s\n---\n", name, rhjson, actualHjson)
	}
	if !jsonOK {
		t.Logf("%s\n---json expected\n%s\n---json actual\n%s\n---\n", name, rjson, actualJson)
	}
	if !hjsonOK || !jsonOK { panic("fail!") }
}

func TestHjson(t *testing.T) {

	files := strings.Split(string(getContent("assets/testlist.txt")), "\n")

	for _, file := range files {
		run(t, file)
	}
}