package main

import (
	"flag"
	"fmt"
)

func main() {
	var classpath, main_class string
	flag.StringVar(&classpath, "classpath", "", "Directory with all JARs used by application")
	flag.StringVar(&main_class, "main", "", "Main class")
	flag.Parse()
	fmt.Printf("Hello, world: %v\n", classpath)

	for _,c := range loadJar("/Users/thundermoose/.m2/repository/com/monitorjbl/json-view/0.14/json-view-0.14.jar") {
		fmt.Printf("%v: %v\n", c.CanonicalName, c.Size)
	}
}