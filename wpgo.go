// wpgo - command-line tool for wordpress rest api

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/automattic/go/jaguar"
	"github.com/postfix/goconf"
)

var c *goconf.ConfigFile
var args []string
var blog_id, token string

var cmds = []string{"read", "post", "stats", "upload"}

// read config and parse args
func init() {
	var configFilename string

	// get user home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	cfgfile := usr.HomeDir + "/.wpgo.conf"
	flag.StringVar(&configFilename, "config", cfgfile, "specify a config file")
	flag.Parse()
	args = flag.Args()

	// confirm file exists
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		log.Fatal("Config file ~/.wpgo.conf does not exists")
	}

	c, err = goconf.ReadConfigFile(configFilename)

	if len(args) < 1 {
		usage()
	}
}

// route command and args
func main() {
	var err error
	// confirm params and config all set
	blog, cmd, param := parse_args()
	token, err = c.GetString(blog, "token")
	if err != nil {
		log.Fatalln("No auth token configured in ~/.wpgo.conf")
	}

	blog_id, err = c.GetString(blog, "blog_id")
	if err != nil {
		log.Fatalln("No blog_id configured in ~/.wpgo.conf")
	}

	switch cmd {
	case "read":
		if param == "" {
			get_latest()
		} else {
			get_single_post(param)
		}

	case "post":
		if param == "" {
			log.Fatalln("No filename specifed to post")
		}
		do_post(param)

	case "stats":
		get_stats(param)

	case "upload":
		if param == "" {
			log.Fatalln("No filename specifed to post")
		}
		upload_media(param)

	default:
		usage()
	}

}

// help
func usage() {
	fmt.Println("Usage: wpgo [blog] [command] [param]")
	os.Exit(1)
}

func parse_args() (blog, cmd, param string) {

	// make sure we have at least one arg
	if len(args) < 1 {
		usage()
	}

	// if first position is a command, then set default blog
	if elemExists(args[0], cmds) {
		blog = "default"
		cmd = args[0]
		if len(args) > 1 {
			param = strings.Join(args[1:], " ")
		}
		return blog, cmd, param
	}

	blog = args[0]
	cmd = args[1]
	if len(args) >= 3 {
		param = strings.Join(args[2:], " ")
	}
	return blog, cmd, param
}

func getApiFetcher(endpoint string) (j jaguar.Jaguar) {
	apiurl := "https://public-api.wordpress.com/rest/v1"
	url := strings.Join([]string{apiurl, "sites", blog_id, endpoint}, "/")

	j = jaguar.New()
	j.Header.Add("Authorization", "Bearer "+token)
	j.Url(url)
	return j
}

func elemExists(s string, a []string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}
	return false
}
