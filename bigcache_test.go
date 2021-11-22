package gocache

import (
	"fmt"
	"testing"
)

func TestSerializeGOB(t *testing.T) {
	bs, err := serializeGOB("ZhouJinke")
	if err != nil {
		fmt.Println(err)
	}

	value, err := deserializeGOB(bs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(value.(string))
	if value.(string) == "ZhouJinke" {
		fmt.Println("func serializeGOB and deserializeGOB ok")
	} else {
		fmt.Println("func serializeGOB and deserializeGOB error")
	}
}
