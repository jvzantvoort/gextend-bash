package config

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	AppName        string
	HomeDir        string
	ConfigDir      string
	ConfigDirMode int
}

const (
	ApplicationName string = "gextend-bash"
	ConfigDirEnv    string = "XTRA_BASH_CONFIG_DIR"
	ConfigDirName   string = "gextend-bash"
)

func GetHomeDir() (string, error) {
	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir, nil
	}
	return "", err
}

// ExpandHome expand the tilde in a given path.
func ExpandHome(pathstr string) (string, error) {

	if len(pathstr) == 0 {
		return pathstr, nil
	}

	if pathstr[0] != '~' {
		return pathstr, nil
	}
	HomeDir, _ := GetHomeDir()

	return filepath.Join(HomeDir, pathstr[1:]), nil

}

// prefix returns a prefix for logging and messages based on function name.
func (c Config) prefix() string {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return fmt.Sprintf("%s.%s", c.AppName, elements[len(elements)-1])
}

// mkdir create directory
func (c Config) mkdir(path string, mode int) {
	log_prefix := c.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	finfo, err := os.Stat(path)
	// we found something
	if err == nil {
		// already exists
		if finfo.IsDir() {
			log.Debugf("found dir: %s", path)
			return
		} else {
			log.Errorf("found target: %s but it is not a directory", path)
		}
	}
	mode_oct := os.FileMode(mode)
	os.MkdirAll(path, mode_oct)

}

// SetDefaultHomeDir get the user's homedir
func (c *Config) SetDefaultHomeDir() error {
	if len(c.HomeDir) == 0 {
		homedir, err := GetHomeDir()
		c.HomeDir = homedir
		return err
	}
	return nil
}

func (c *Config) SetDefaultConfigDir() {
	if len(c.ConfigDir) != 0 {
		return
	}

	c.SetDefaultHomeDir()

	// check environment variable
	item_path, item_path_set := os.LookupEnv(ConfigDirEnv)

	if item_path_set {
		c.ConfigDir = item_path
	} else {
		c.ConfigDir = path.Join(c.HomeDir, ".config", ConfigDirName)
		c.mkdir(c.ConfigDir, c.ConfigDirMode) // make sure the directory exists
	}
}

func (c *Config) Initialize() {
	c.AppName = ApplicationName
	c.ConfigDirMode = 0755
	c.SetDefaultHomeDir()
	c.SetDefaultConfigDir()
}

func NewConfig() *Config {
	retv := &Config{}
	retv.Initialize()
	return retv
}
