package main

import (
	"flag"
	"fmt"
	"github.com/monitorjbl/noclip/bytecode"
)

func main() {
	var classpath, main_class string
	flag.StringVar(&classpath, "classpath", "", "Directory with all JARs used by application")
	flag.StringVar(&main_class, "main", "", "Main class")
	flag.Parse()

	for _, c := range bytecode.LoadJar("/Users/thundermoose/.m2/repository/com/monitorjbl/json-view/0.14/json-view-0.14.jar") {
		fmt.Printf("%v\n", c.ToString())
	}
}