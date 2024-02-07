package main

import (
	"context"
	"fmt"

	"github.com/ironzhang/superdnsgo"
)

func main() {
	err := superdnsgo.AutoSetup()
	if err != nil {
		fmt.Printf("superdnsgo auto setup: %v\n", err)
		return
	}

	addr, cluster, err := superdnsgo.Lookup(context.Background(), "www.superdns.com", nil)
	if err != nil {
		fmt.Printf("superdnsgo lookup endpoint: %v\n", err)
		return
	}
	fmt.Printf("cluster=%s, address=%v\n", cluster, addr)
}
