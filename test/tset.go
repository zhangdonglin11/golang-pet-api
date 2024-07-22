package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now())
}
func createId(uid, toUid interface{}) string {
	sprintf := fmt.Sprintf("%v->%v", uid, toUid)
	return sprintf
}
