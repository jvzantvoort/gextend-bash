// Package config provides configuration data globally used
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/jvzantvoort/gextend-bash/utils"
	log "github.com/sirupsen/logrus"
)

const (
	SettingsFile string = "settings.json"
	RegistryFile string = "registry.json"
)

// UserConfig

type UserConfig struct {
	MailAddress string `json:"mailaddress"`
	Company     string `json:"company"`
	Copyright   string `json:"copyright"`
	License     string `json:"license"`
	User        string `json:"user"`
	Username    string `json:"username"`
}

// MainConfig configuration for gextend-bash
type MainConfig struct {
	ForceInit      bool
	HomeDir        string
	ArchiveDir     string
	ConfigDir      string
	TemplatesDir   string
	ConfigDirPerms int
	AppVersion     string     `json:"version"`
	UserConfig     UserConfig `json:"user"`
}

func DefaultString(args ...string) string {
	for _, def := range args {
		if len(def) != 0 {
			return def
		}
	}
	return ""
}

func (u *UserConfig) SetMailAddress(args ...string) error {
	newval := DefaultString(args...)
	if len(newval) == 0 && len(u.MailAddress) == 0 {
		return fmt.Errorf("MailAddress not set")
	}
	if len(newval) != 0 {
		u.MailAddress = newval
	}
	return nil
}

func (u *UserConfig) SetCompany(args ...string) error {
	newval := DefaultString(args...)
	if len(newval) == 0 && len(u.Company) == 0 {
		return fmt.Errorf("Company not set")
	}
	if len(newval) != 0 {
		u.Company = newval
	}
	return nil
}

func (u *UserConfig) SetCopyright(args ...string) error {
	newval := DefaultString(args...)
	if len(newval) == 0 && len(u.Copyright) == 0 {
		return fmt.Errorf("Copyright not set")
	}
	if len(newval) != 0 {
		u.Copyright = newval
	}
	return nil
}

func (u *UserConfig) SetLicense(args ...string) error {
	newval := DefaultString(args...)
	if len(newval) == 0 && len(u.License) == 0 {
		return fmt.Errorf("License not set")
	}
	if len(newval) != 0 {
		u.License = newval
	}
	return nil
}

func (u *UserConfig) SetUser(args ...string) error {
	newval := DefaultString(args...)
	if len(newval) == 0 && len(u.User) == 0 {
		return fmt.Errorf("User not set")
	}
	if len(newval) != 0 {
		u.User = newval
	}
	return nil
}

func (u *UserConfig) SetUsername(args ...string) error {
	newval := DefaultString(args...)
	if len(newval) == 0 && len(u.Username) == 0 {
		return fmt.Errorf("Username not set")
	}
	if len(newval) != 0 {
		u.Username = newval
	}
	return nil
}

// GetHomeDir get the user's homedir
func (m *MainConfig) GetHomeDir() string {
	if len(m.HomeDir) != 0 {
		return m.HomeDir
	}
	m.HomeDir = utils.GetHomeDir()

	return m.HomeDir
}

func (m MainConfig) ConfigFile(name string) string {
	return path.Join(m.ConfigDir, name)
}

func (m *MainConfig) GetConfigDir() string {
	// return cached value
	if len(m.ConfigDir) != 0 {
		return m.ConfigDir
	}

	// check environment variable
	item_path, item_path_set := os.LookupEnv("GOPROJ_CONFIG_DIR")

	if item_path_set {
		m.ConfigDir = item_path
		return m.ConfigDir
	}

	m.GetHomeDir()

	m.ConfigDir = path.Join(m.HomeDir, ".config", "gextend-bash")

	if runtime.GOOS == "windows" {
		m.ConfigDir = path.Join(m.HomeDir, "GOProj")
	}

	return m.ConfigDir
}

func (m *MainConfig) GetTemplatesDir() string {
	if len(m.TemplatesDir) != 0 {
		return m.TemplatesDir
	}
	m.TemplatesDir = path.Join(m.GetConfigDir(), "templates.d")
	return m.TemplatesDir
}

func (m *MainConfig) GetArchiveDir() string {
	// return cached value
	if len(m.ArchiveDir) != 0 {
		return m.ArchiveDir
	}

	// check environment variable
	item_path, item_path_set := os.LookupEnv("GOPROJ_ARCHIVE_DIR")

	if item_path_set {
		m.ArchiveDir = item_path
		return m.ArchiveDir
	}

	m.GetHomeDir()

	m.ArchiveDir = path.Join(m.HomeDir, "Archive", "gextend-bash")

	if runtime.GOOS == "windows" {
		m.ArchiveDir = path.Join(m.HomeDir, "GOProjArchive")
	}

	return m.ArchiveDir
}

// Init initialize the MainConfig struct
func (m *MainConfig) Init() {

	m.GetHomeDir()

	m.GetConfigDir()

	m.GetArchiveDir()

	m.GetTemplatesDir()

	if runtime.GOOS == "windows" {
		m.ConfigDirPerms = 0777
	} else {
		m.ConfigDirPerms = 0755
	}

	if !m.ForceInit {
		m.ReadFromFile(SettingsFile)
	}
}

func (m MainConfig) Save() {
	log.Debugf("Write to %s", SettingsFile)
	defer log.Debugf("Write to %s, end", SettingsFile)
	m.WriteToFile(SettingsFile)
}

// CreateDirs create the main config dir
func (m MainConfig) CreateDirs() {
	mode := m.ConfigDirPerms

	err := utils.MkdirP(m.ConfigDir, int(mode))
	if err != nil {
		log.Errorf("Failed to create dir %s, %v", m.ConfigDir, err)
	}

	err = utils.MkdirP(m.ArchiveDir, int(mode))
	if err != nil {
		log.Errorf("Failed to create dir %s, %v", m.ArchiveDir, err)
	}

	err = utils.MkdirP(m.TemplatesDir, int(mode))
	if err != nil {
		log.Errorf("Failed to create dir %s, %v", m.TemplatesDir, err)
	}
}

// File handling
//
// Read
func (m *MainConfig) Read(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	return nil
}

// Write
func (m MainConfig) Write(writer io.Writer) error {
	jsonString, err := json.MarshalIndent(m, "", "  ")
	if err == nil {
		fmt.Fprint(writer, string(jsonString))
		fmt.Fprintf(writer, "\n")
	}
	return err

}

// ReadFromFile
func (m *MainConfig) ReadFromFile(name string) error {

	fileh, err := os.Open(m.ConfigFile(name))
	defer fileh.Close()

	if err != nil {
		return err
	}

	return m.Read(fileh)
}

// WriteToFile
func (m MainConfig) WriteToFile(name string) error {
	log.Debugf("save to: %s, start", m.ConfigFile(name))
	defer log.Debugf("save to: %s, end", m.ConfigFile(name))

	fileh, err := os.OpenFile(m.ConfigFile(name), os.O_WRONLY|os.O_CREATE, 0644)
	defer fileh.Close()

	if err != nil {
		return err
	}
	return m.Write(fileh)
}

func (m *MainConfig) ResetConfig() {
	m.ForceInit = true
	m.Init()
	m.CreateDirs()
}

func (m *MainConfig) Get(name string) (string, error) {
	switch name {
	case "version":
		return m.AppVersion, nil
	case "user.mailaddress":
		return m.UserConfig.MailAddress, nil
	case "user.company":
		return m.UserConfig.Company, nil
	case "user.copyright":
		return m.UserConfig.Copyright, nil
	case "user.license":
		return m.UserConfig.License, nil
	case "user.user":
		return m.UserConfig.User, nil
	case "user.username":
		return m.UserConfig.Username, nil
	default:
		return "", fmt.Errorf("Illegal field: %s", name)
	}
}

func (m MainConfig) Fields() []string {
	fields := []string{}
	fields = append(fields, "version")
	fields = append(fields, "user.mailaddress")
	fields = append(fields, "user.company")
	fields = append(fields, "user.copyright")
	fields = append(fields, "user.license")
	fields = append(fields, "user.user")
	fields = append(fields, "user.username")
	return fields
}

// NewMainConfig initialize a MainConfig and initialize it.
//
//	mc := config.NewMainConfig()
//	fmt.Printf("config dir: %s\n", mc.ConfigDir)
func NewMainConfig() *MainConfig {
	retv := &MainConfig{}
	retv.ForceInit = false

	retv.Init()

	retv.CreateDirs()

	return retv

}
