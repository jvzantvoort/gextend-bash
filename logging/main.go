package logging

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type LogMessage struct {
	Tag       string    `json:"tag"`
	File      string    `json:"-"`
	Priority  string    `json:"priority"`
	SkipEmpty bool      `json:"-"`
	StdErr    bool      `json:"-"`
	Message   string    `json:"message"`
	Time      time.Time `json:"time"`
}

type LogMessages struct {
	messages []LogMessage
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

	l.Tag = GetString(*cmd, "tag")
	prio := GetString(*cmd, "priority")
	l.SetLevel(prio)
	l.Message = strings.Join(args, " ")

}

func (l LogMessage) Print() error {
	fileh, err := os.OpenFile(l.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer fileh.Close()
	if err != nil {
		return err
	}

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

	return retv
}

func NewLogMessages(inputfile string) *LogMessages {
	retv := &LogMessages{}
	filehandle, err := os.Open(inputfile)
	defer filehandle.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(filehandle)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		obj := &LogMessage{}
		json.Unmarshal([]byte(scanner.Bytes()), &obj)
		retv.messages = append(retv.messages, *obj)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return retv
}
