package main

import (
	"bufio"
	"bytes"
	"os"
	"reflect"
	"testing"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

func TestMain(t *testing.T) {
	if os.Getenv("ACG_CONFIG") == "" {
		println("ACG_CONFIFG is not set. Exiting.")
		return
	}

	v := reflect.ValueOf(clowder.LoadedConfig)

	f := &bytes.Buffer{}

	w := bufio.NewWriter(f)
	w.WriteString("#!/bin/bash\n\n")

	req_print("CLOWDER", v, w)

	w.Flush()

	dat, err := os.ReadFile("vars.test")
	if err != nil {
		t.Error("coulnd't read file")
	}

	str := string(dat)

	if f.String() != str {
		t.Fatal("strings didn't match")
	}
}
