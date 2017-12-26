package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strconv"
)

func GenSign(srcMap map[string]interface{}, secret string) (string, error) {
	//TODO GET Secret from redis
	srcMap["AppKey"] = secret
	keys := make([]string, len(srcMap))
	i := 0
	for k, _ := range srcMap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	signStr := ""
	var value interface{}
	var v string
	for _, k := range keys {
		if k == "Sign" {
			continue
		}
		value = srcMap[k]
		switch value := value.(type) {

		case bool:
			v = strconv.FormatBool(value)
		case int:
			v = strconv.Itoa(value)
		case uint:
			v = strconv.FormatUint(uint64(value), 10)
		case int8:
			v = strconv.Itoa(int(value))
		case int16:
			v = strconv.Itoa(int(value))
		case int32:
			v = strconv.Itoa(int(value))
		case int64:
			v = strconv.FormatInt(value, 10)
		case float32:
			v = strconv.FormatFloat(float64(value), 'f', -1, 64)
		case float64:
			v = strconv.FormatFloat(value, 'f', -1, 64)
		case string:
			v = value
		default:
			continue

		}
		signStr = signStr + k + v
	}
	fmt.Println("Sign str is :", signStr)
	t := sha1.New()
	io.WriteString(t, signStr)
	out := fmt.Sprintf("%x", t.Sum(nil))
	return out, nil
}
