package config

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type SectLogging struct {
	OutputDir string `ini:"outputdir"`
	OutputFile string `ini:"outputfile"`
	FileMode  int    `ini:"mode" comment:"mode of the logfiles"`
	MaxLines  int    `ini:"max_lines"`
}

type ConfigLogging struct {
	AppName        string `ini:"-"`
	ConfigFile     string `ini:"-"`
	ConfigFileMode int    `ini:"-"`
	ConfigDirMode  int    `ini:"-"`
	Config         `ini:"-"`

	SectLogging `ini:"main"`
}

// prefix returns a prefix for logging and messages based on function name.
func (cl ConfigLogging) prefix() string {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return fmt.Sprintf("%s.%s", cl.Config.AppName, elements[len(elements)-1])
}

func (cl ConfigLogging) Parse(templatestring string) (string, error) {

	var retv string
	buf := new(bytes.Buffer)

	tmpl, err := template.New("template").Parse(templatestring)
	if err != nil {
		return templatestring, err
	}

	err = tmpl.Execute(buf, cl)
	if err != nil {
		return templatestring, err
	}
	retv = buf.String()
	if retv == templatestring {
		log.Debugf("  no changes on variable")
	}

	return retv, nil
}

func (cl ConfigLogging) GetOutputDir() (string, error) {
	retv, err := cl.Parse(cl.OutputDir)
	if err != nil {
		return retv, err
	}
	return ExpandHome(retv)
}

func (cl ConfigLogging) GetOutputFile() (string, error) {
	return cl.Parse(cl.OutputFile)
}

func (cl ConfigLogging) LogfilePath() (string, error) {
	dirn, err := cl.GetOutputDir()
	if err != nil {
		return "", err
	}
	filen, err := cl.GetOutputFile()
	if err != nil {
		return "", err
	}
	return path.Join(dirn, filen), nil
}

func (cl ConfigLogging) Write() error {

	// Setup logging
	log_prefix := cl.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	log.Debugf("configfile: %s", cl.ConfigFile)

	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, &cl)
	if err != nil {
		return err
	}

	err = cfg.SaveTo(cl.ConfigFile)
	if err != nil {
		return err
	}
	mode_oct := os.FileMode(cl.ConfigFileMode)
	os.Chmod(cl.ConfigFile, mode_oct)
	return nil
}

func (cl ConfigLogging) FileExists() bool {
	_, err := os.Stat(cl.ConfigFile)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (cl ConfigLogging) CreateIfEmpty() error {

	// Setup logging
	log_prefix := cl.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	if cl.FileExists() {
		log.Debugf("%s: configfile already exists", log_prefix)
		return nil
	}

	cl.Initialize()

	return cl.Write()
}

func (cl *ConfigLogging) Read() error {

	// Setup logging
	log_prefix := cl.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	var cfg *ini.File
	cfg, err := ini.Load(cl.ConfigFile)
	if err != nil {
		return err
	}

	return cfg.MapTo(cl)
}

func (cl *ConfigLogging) Initialize() {
	ncfg := NewConfig()
	cl.Config = *ncfg
	cl.AppName = cl.Config.AppName

	// top level elements (mostly for templating)
	cl.ConfigDir = cl.Config.ConfigDir
	cl.ConfigDirMode = cl.Config.ConfigDirMode
	cl.ConfigFile = path.Join(cl.ConfigDir, "logging.ini")
	cl.ConfigFileMode = 0644

	// Setup logging
	log_prefix := cl.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	retv := cl.Read()
	if retv == nil {
		return
	}
	fmt.Printf("%v\n", retv)

	sect_logging := &SectLogging{}
	sect_logging.FileMode = 0644
	sect_logging.MaxLines = 10000
	cl.SectLogging = *sect_logging
}

func NewConfigLogging() *ConfigLogging {
	retv := &ConfigLogging{}
	retv.Initialize()
	return retv
}
