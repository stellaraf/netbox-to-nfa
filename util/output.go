package util

import (
	"encoding/json"
	"fmt"
	"os"
)

// Check checks an error, and if present, logs an error message to the console & exits.
func Check(m string, e error, format ...interface{}) {
	args := append(format, e)
	if e != nil {
		fmt.Println(fmt.Errorf(m+"\n%v", args...))
		os.Exit(1)
	}
}

func PrettyStruct(s interface{}) string {
	i, _ := json.MarshalIndent(s, "", "\t")
	return string(i)
}
