package config

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/jvzantvoort/gextend-bash/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// SectLogging missing godoc.
type SectLogging struct {
	OutputDir  string    `ini:"outputdir"`
	OutputFile string    `ini:"outputfile"`
	FileMode   int       `ini:"mode" comment:"mode of the logfiles"`
	MaxLines   int       `ini:"max_lines"`
	Now        time.Time `ini:"-"`
}

// ConfigLogging missing godoc.
type ConfigLogging struct {
	AppName        string    `ini:"-"`
	ConfigFile     string    `ini:"-"`
	TemplateFields []string  `ini:"-"`
	Hostname   string    `ini:"-"`
	Now            time.Time `ini:"-"`
	Config         `ini:"-"`
	SectLogging    `ini:"main"`
}

func TmplLookupEnv(variablename string, args ...string) string {
	retv, ok := os.LookupEnv(variablename)
	if ok {
		return retv
	}
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

// prefix returns a prefix for logging and messages based on function name.
func (cl ConfigLogging) prefix() string {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return fmt.Sprintf("%s.%s", cl.Config.AppName, elements[len(elements)-1])
}

// Parse missing godoc.
func (cl ConfigLogging) Parse(templatestring string) (string, error) {
	var retv string
	buf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"env": TmplLookupEnv,
	}

	tmpl, err := template.New("template").Funcs(funcMap).Parse(templatestring)
	if err != nil {
		log.Errorf("Error: %s", err)
		return templatestring, err
	}

	err = tmpl.Execute(buf, cl)
	if err != nil {
		log.Errorf("Template string: %s", templatestring)
		log.Errorf("Error: %s", err)
		utils.ExitOnError(err)
	}
	retv = buf.String()
	if retv == templatestring {
		log.Debugf("  no changes on variable")
	}

	return retv, nil
}

// GetOutputDir missing godoc.
func (cl ConfigLogging) GetOutputDir() (string, error) {
	retv, err := cl.Parse(cl.OutputDir)
	if err != nil {
		return retv, err
	}
	return ExpandHome(retv)
}

// GetOutputFile missing godoc.
func (cl ConfigLogging) GetOutputFile() (string, error) {
	return cl.Parse(cl.OutputFile)
}

// LogfilePath missing godoc.
func (cl ConfigLogging) LogfilePath() (string, error) {
	dirn, err := cl.GetOutputDir()
	if err != nil {
		log.Errorf("GetOutputDir failed: %s", err)
		return "", err
	}
	filen, err := cl.GetOutputFile()
	if err != nil {
		return "", err
	}
	return path.Join(dirn, filen), nil
}

// Write missing godoc.
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
	mode_oct := os.FileMode(ConstFileMode)
	err = os.Chmod(cl.ConfigFile, mode_oct)
	if err != nil {
		return err
	}
	return nil
}

// FileExists missing godoc.
func (cl ConfigLogging) FileExists() bool {
	_, err := os.Stat(cl.ConfigFile)
	if err == nil {
		return true
	} else {
		return false
	}
}

// CreateIfEmpty missing godoc.
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

// Read missing godoc.
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

// Initialize missing godoc.
func (cl *ConfigLogging) Initialize() {
	// Setup logging
	log_prefix := "gextend-bash.Initialize"
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	// initialize Config object
	ncfg := NewConfig()
	cl.Config = *ncfg
	cl.AppName = cl.Config.AppName
	log.Debugf("  set AppName to %s", cl.AppName)

	// top level elements (mostly for templating)
	cl.ConfigDir = cl.Config.ConfigDir
	log.Debugf("  set ConfigDir to %s", cl.ConfigDir)

	cl.ConfigFile = path.Join(cl.ConfigDir, "logging.ini")
	log.Debugf("  set ConfigFile to %s", cl.ConfigFile)

	log.Debugf("  try to read config file")
	retv := cl.Read()
	if retv != nil {
		log.Debugf("  @@ try to read config file, failed")
		sect_logging := NewSectLogging()
		log.Debugf("  Initialzed SectLogging")
		cl.SectLogging = *sect_logging
		log.Debugf("  Initialzed SectLogging")
		retv = cl.Write()
		if retv != nil {
			log.Errorf("Error: %s", retv)
		}
	} else {
		log.Debugf("  try to read config file, success")
	}
	cl.Hostname = utils.ShortHostname()
	cl.SectLogging.InitializeVariableFields()

}

// NewConfigLogging missing godoc.
func NewConfigLogging() *ConfigLogging {
	retv := &ConfigLogging{}
	retv.Initialize()
	retv.Now = time.Now()
	return retv
}

//
//
// -----------------------------------------------------------------------------

// prefix returns a prefix for logging and messages based on function name.
func (sl SectLogging) prefix() string {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return fmt.Sprintf("%s.%s", ApplicationName, elements[len(elements)-1])
}

// InitializeVariableFields update variable fields
func (sl *SectLogging) InitializeVariableFields() {

	// Setup logging
	log_prefix := sl.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

}

// Initialize update variable fields
func (sl *SectLogging) Initialize() {

	// Setup logging
	log_prefix := sl.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	sl.OutputDir = "~/Logs"
	sl.OutputFile = "common.log"
	sl.FileMode = ConstFileMode
	sl.MaxLines = 10000

	sl.InitializeVariableFields()
}

func NewSectLogging() *SectLogging {
	retv := &SectLogging{}
	retv.Initialize()
	return retv
}
