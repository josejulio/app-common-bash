package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

func main() {
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
	fmt.Print(f.String())
}

func exportVariable(w *bufio.Writer, name string, value string) {
	if strings.Contains(value, "\n") {
		w.WriteString(fmt.Sprintf(`
read -r -d '' %s <<'EOF'
%s
EOF`, name, value))
		w.WriteString("\n\n")
	} else {
		w.WriteString(fmt.Sprintf("export %s=\"%s\"\n", name, value))
	}
}

func req_print(prefix string, ob reflect.Value, w *bufio.Writer) {
	ob = reflect.Indirect(ob)

	if !ob.IsValid() {
		return
	}

	if ob.Type().String() == "bool" {
		exportVariable(w, strings.ToUpper(prefix), strconv.FormatBool(reflect.Value(ob).Bool()))
		return
	}

	if ob.Type().String() == "string" {
		exportVariable(w, strings.ToUpper(prefix), reflect.Value(ob).String())
		return
	}

	if ob.Type().String() == "int" {
		exportVariable(w, strings.ToUpper(prefix), strconv.FormatInt(reflect.Value(ob).Int(), 10))
		return
	}

	if ob.Kind() == reflect.Slice {
		for i := 0; i < ob.Len(); i++ {
			req_print(prefix+"_"+strconv.Itoa(i), ob.Index(i), w)
		}
		return
	}

	for i := 0; i < ob.NumField(); i++ {
		newObj := reflect.ValueOf(ob.Field(i).Interface())
		req_print(prefix+"_"+ob.Type().Field(i).Name, newObj, w)
	}
}
