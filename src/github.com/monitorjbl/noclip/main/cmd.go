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

	zip:="/Users/thundermoose/Downloads/commons-lang-2.6/org/apache/commons/lang/StringUtils.class.zip"
	//zip:="/Users/thundermoose/.m2/repository/commons-lang/commons-lang/2.6/commons-lang-2.6.jar"
	for _, c := range bytecode.LoadJar(zip) {
		fmt.Printf("%v\n", c.ToString())
	}
}