package logging

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jvzantvoort/gextend-bash/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type LogMessage struct {
	Tag                  string    `json:"tag"`
	File                 string    `json:"-"`
	Priority             string    `json:"priority"`
	SkipEmpty            bool      `json:"-"`
	StdErr               bool      `json:"-"`
	Message              string    `json:"message"`
	Time                 time.Time `json:"time"`
	config.ConfigLogging `json:"-"`
}

type LogMessages struct {
	Messages []LogMessage
}

func GetString(cmd cobra.Command, name string) string {
	retv, _ := cmd.Flags().GetString(name)
	if len(retv) != 0 {
		log.Debugf("Found %s as %s", name, retv)
	}
	return retv
}

func (l LogMessage) MakeString() []byte {
	retv := l.Message
	if len(l.Tag) != 0 {
		retv = fmt.Sprintf("[%s] %s", l.Tag, retv)
	}
	retv = fmt.Sprintf("%s %s", l.Priority, retv)
	retv = fmt.Sprintf("%s %s", l.Time.Format(time.RFC3339), retv)
	retv += "\n"
	return []byte(retv)
}

func (l LogMessage) MakeJSONString() ([]byte, error) {
	dst, err := json.Marshal(l)
	return dst, err
}

func (l *LogMessage) SetLevel(level string) {
	ulevel := strings.ToUpper(level)
	if len(ulevel) == 0 {
		ulevel = l.Priority
	}
	switch ulevel {
	case "EMERG":
		l.Priority = ulevel
	case "ALERT":
		l.Priority = ulevel
	case "CRIT":
		l.Priority = ulevel
	case "ERR":
		l.Priority = ulevel
	case "WARNING":
		l.Priority = ulevel
	case "NOTICE":
		l.Priority = ulevel
	case "INFO":
		l.Priority = ulevel
	case "DEBUG":
		l.Priority = ulevel
	case "PANIC":
		l.Priority = "EMERG"
	case "ERROR":
		l.Priority = "ERR"
	case "WARN":
		l.Priority = "WARNING"
	default:
		log.Errorf("Invalid priority: %s", ulevel)
		l.Priority = "NOTICE"
	}
}

func (l *LogMessage) ImportArgs(cmd *cobra.Command, args []string) {

	l.StdErr, _ = cmd.Flags().GetBool("stderr")
	l.SkipEmpty, _ = cmd.Flags().GetBool("skip-empty")
	l.File = GetString(*cmd, "file")
	if len(l.File) == 0 {
		l.File, _ = l.ConfigLogging.LogfilePath()
	}
	log.Debugf("file is %s", l.File)

	l.Tag = GetString(*cmd, "tag")
	prio := GetString(*cmd, "priority")
	l.SetLevel(prio)
	l.Message = strings.Join(args, " ")

}

// mkdir create directory
func (l LogMessage) mkdir(path string) {
	mode := 0755
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
	if err := os.MkdirAll(path, mode_oct); err != nil {
		log.Errorf("directory cannot be created: %s", path)
	}

}

func (l LogMessage) Print() error {
	// make parent dirs
	l.mkdir(filepath.Dir(l.File))

	fileh, err := os.OpenFile(l.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fileh.Close()

	msg := l.MakeString()
	jmsg, _ := l.MakeJSONString()

	if l.StdErr {
		_, err = os.Stderr.Write(msg)
	}

	if err != nil {
		return err
	}
	xmsg := []byte("\n")
	jmsg = append(jmsg, xmsg...)

	_, err = fileh.Write(jmsg)
	return err

}

func NewLogMessage(level string) *LogMessage {
	retv := &LogMessage{}
	retv.Time = time.Now()
	retv.Priority = strings.ToUpper(level)
	cfg := config.NewConfigLogging()
	retv.ConfigLogging = *cfg

	return retv
}

func NewLogMessages(inputfile string) *LogMessages {
	retv := &LogMessages{}
	filehandle, err := os.Open(inputfile)
	if err != nil {
		log.Fatal(err)
	}
	defer filehandle.Close()

	scanner := bufio.NewScanner(filehandle)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		obj := &LogMessage{}
		if err := json.Unmarshal([]byte(scanner.Bytes()), &obj); err != nil {
			log.Errorf("Error: %s", err)
		}
		retv.messages = append(retv.messages, *obj)
	}

	// Sort in ascending order
	sort.Slice(retv.messages, func(i, j int) bool {
		return retv.messages[i].Time.Before(retv.messages[j].Time)
	})

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return retv
}
