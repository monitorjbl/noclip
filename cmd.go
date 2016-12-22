package main

import (
	"flag"
	"github.com/monitorjbl/noclip/bytecode"
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("main")

func configureLogging(debug bool) {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.6s} %{id:03x}%{color:reset} %{message}`)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(backendFormatter)

	if debug {
		leveledBackend.SetLevel(logging.DEBUG, "")
	} else {
		leveledBackend.SetLevel(logging.INFO, "")
	}

	logging.SetBackend(leveledBackend)
}

func main() {
	var classpath, main_class string
	var debug bool
	flag.StringVar(&classpath, "classpath", "", "Directory with all JARs used by application")
	flag.StringVar(&main_class, "main", "", "Main class")
	flag.BoolVar(&debug, "debug", false, "Turn on debug logging")
	flag.Parse()

	configureLogging(debug)
	jars := []string{
		"/Users/thundermoose/Downloads/commons-lang-2.6/org/apache/commons/lang/StringUtils.class.zip",
		"/Users/thundermoose/.m2/repository/commons-lang/commons-lang/2.6/commons-lang-2.6.jar",
		"/Users/thundermoose/.m2/repository/org/springframework/org.springframework.aop/3.2.3.RELEASE/org.springframework.aop-3.2.3.RELEASE.jar",
		"/Users/thundermoose/.m2/repository/com/monitorjbl/json-view/0.14/json-view-0.14.jar",

	}
	//zip:="/Users/thundermoose/Downloads/commons-lang-2.6/org/apache/commons/lang/StringUtils.class.zip"
	//zip:="/Users/thundermoose/.m2/repository/commons-lang/commons-lang/2.6/commons-lang-2.6.jar"
	//zip:="/Users/thundermoose/.m2/repository/org/springframework/org.springframework.aop/3.2.3.RELEASE/org.springframework.aop-3.2.3.RELEASE.jar"
	//zip := "/Users/thundermoose/.m2/repository/com/monitorjbl/json-view/0.14/json-view-0.14.jar"

	for _, jar := range jars {
		for _, c := range bytecode.LoadJar(jar) {
			log.Infof("%v\n", c.ToString())
		}
	}
}