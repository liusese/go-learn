package main

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func main() {
	// uri, _ := url.Parse(path.Join("http://www.baidu.com", "/hello", "/world1"))
	// fmt.Println(uri.String())
	var m map[string]string
	fmt.Println(len(m))

	jsonStr := `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
	val := gjson.Get(jsonStr, "age")
	fmt.Println(val.String())
}
