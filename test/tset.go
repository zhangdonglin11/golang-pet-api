package main

import "fmt"

func main() {
	fmt.Println(createId(11.12, "22"))
}
func createId(uid, toUid interface{}) string {
	sprintf := fmt.Sprintf("%v->%v", uid, toUid)
	return sprintf
}
