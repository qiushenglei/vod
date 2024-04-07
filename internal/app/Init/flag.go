package Init

import "flag"

var Env string
var HttpConfFile string
var LalConfFile string

func InitFlag() {
	flag.Parsed()
	flag.StringVar(&Env, "Env", "local", "environment")
	flag.StringVar(&HttpConfFile, "httpconf", "vod.conf.json", "use http conf HttpConfFile")
	flag.StringVar(&LalConfFile, "lalconf", "lal.conf.json", "use lal conf HttpConfFile")
}
