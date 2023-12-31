# superdnsgo

## 1. Overview

superdnsgo is a Go client for superdns, it supports service discovery like dns and dynamic configuration.

## 2. Quick Start

```
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

	endpoint, cluster, err := superdnsgo.LookupEndpoint(context.Background(), "www.superdns.com", nil)
	if err != nil {
		fmt.Printf("superdnsgo lookup endpoint: %v\n", err)
		return
	}
	fmt.Printf("cluster=%s, endpoint=%v\n", cluster, endpoint)
}
```

