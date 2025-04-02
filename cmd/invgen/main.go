package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
	log "github.com/sirupsen/logrus"
)

type Host struct {
	Name      string
	Continent string
	Region    string
	City      string
	Team      string
	Purpose   string
}

// Geographic hierarchy
var locations = map[string][][]string{
	"US": {
		{"US", "NY", "NewYork"},
		//		{"US", "CA", "sanfransisco"},
		{"US", "TX", "austen"},
	},
	"EU": {
		{"EU", "NL", "AMSterdam"},
		{"EU", "NL", "Eindhoven"},
		//		{"EU", "DE", "BERlin"},
		{"EU", "FR", "PARis"},
	},
	/*
		"ASIA": {
			{"ASIA", "JP", "TOKio"},
			{"ASIA", "IN", "DELi"},
			{"ASIA", "CN", "bejing"},
		},
	*/
}

var Teams = []string{
	"team1",
	"team2",
	"team3",
	"team4",
	"team5",
}
var Purpose = []string{
	"prod",
	"dev",
	"qa",
	"test",
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

}

func LogIfError(msg interface{}) {
	if msg == nil {
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	log.Errorf("%s: return not nil: %s", elements[len(elements)-1], msg)
}
func LogIfWriteError(n int, err error) {
	if err == nil {
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	log.Errorf("%s: return not nil: %s", elements[len(elements)-1], err)
}

func mkdirRecursive(directory string) {
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		LogIfError(os.MkdirAll(directory, 0755))
	}
}

func printMkDirPerListItem(prefix string, fields []string) {
	directories := []string{}
	for _, field := range fields {
		directories = append(directories, fmt.Sprintf("%s_%s", prefix, field))
	}

	for _, field := range directories {
		mkdirRecursive(fmt.Sprintf("group_vars/%s", field))

		// write a file in the directory
		w, err := os.Create(fmt.Sprintf("group_vars/%s/vars.yml", field))
		if err != nil {
			fmt.Println("Error: Cannot create file")
			return
		}
		defer w.Close()
		LogIfWriteError(w.WriteString("---\n"))
		LogIfWriteError(w.WriteString(fmt.Sprintf("%s_token: %s\n", prefix, field)))
		LogIfWriteError(w.WriteString("...\n"))

	}
}

func printMkDirPerLocation() {
	prefix := "l"
	directories := []string{}
	for continent, regions := range locations {
		for _, region := range regions {
			col1 := strings.ToLower(continent)
			col2 := strings.ToLower(region[1])
			col3 := strings.ToLower(region[2])

			directories = append(directories, fmt.Sprintf("%s_%s_%s_%s", prefix, col1, col2, col3))
			directories = append(directories, fmt.Sprintf("%s_%s_%s", prefix, col1, col2))
			directories = append(directories, fmt.Sprintf("%s_%s", prefix, col1))

		}
	}
	for _, field := range directories {
		mkdirRecursive(fmt.Sprintf("group_vars/%s", field))

		// write a file in the directory
		w, err := os.Create(fmt.Sprintf("group_vars/%s/vars.yml", field))
		if err != nil {
			fmt.Println("Error: Cannot create file")
			return
		}
		defer w.Close()
		LogIfWriteError(w.WriteString("---\n"))
		LogIfWriteError(w.WriteString(fmt.Sprintf("%s_token: %s\n", prefix, field)))
		LogIfWriteError(w.WriteString("...\n"))

	}
}

func groupHostsByFields(hosts []Host, prefix string, fields []string) map[string][]string {
	retv := make(map[string][]string)

	for _, host := range hosts {
		index := ""
		for _, field := range fields {
			switch field {
			case "Continent":
				index = host.Continent
			case "Region":
				index = fmt.Sprintf("%s_%s", host.Continent, host.Region)
			case "City":
				index = fmt.Sprintf("%s_%s_%s", host.Continent, host.Region, host.City)
			case "Team":
				index = host.Team
			case "Purpose":
				index = host.Purpose
			}
		}
		index = fmt.Sprintf("%s_%s", prefix, index)
		retv[index] = append(retv[index], host.Name)
	}
	return retv
}

func groupHostsByContinent(hosts []Host) map[string][]string {
	return groupHostsByFields(hosts, "l", []string{"Continent"})
}

func groupHostsByRegion(hosts []Host) map[string][]string {
	return groupHostsByFields(hosts, "l", []string{"Continent", "Region"})
}

func groupHostsByCity(hosts []Host) map[string][]string {
	return groupHostsByFields(hosts, "l", []string{"Continent", "Region", "City"})
}

func groupHostsByTeam(hosts []Host) map[string][]string {
	return groupHostsByFields(hosts, "t", []string{"Team"})
}

func groupHostsByPurpose(hosts []Host) map[string][]string {
	return groupHostsByFields(hosts, "p", []string{"Purpose"})
}

func writeGroup(w *os.File, group map[string][]string, priority int) {
	for continent, names := range group {
		fmt.Fprintf(w, "%s:\n", continent)
		fmt.Fprintf(w, "  vars:\n")
		fmt.Fprintf(w, "    ansible_group_priority: %d\n", priority)
		fmt.Fprintf(w, "  hosts:\n")
		for _, hn := range names {
			fmt.Fprintf(w, "    %s:\n", hn)
		}
	}
}

func RandomFromList(inlist []string) string {
	// rand.Seed(time.Now().UnixNano())
	return inlist[rand.Intn(len(inlist))]
}

// Generate a random location slice
func randomLocation() []string {
	// rand.Seed(time.Now().UnixNano())

	// Get random continent key
	keys := make([]string, 0, len(locations))
	for k := range locations {
		keys = append(keys, k)
	}
	randomKey := keys[rand.Intn(len(keys))]

	// Get random location slice from that key
	place := locations[randomKey]
	return place[rand.Intn(len(place))]
}

func generateHost() *Host {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	items := randomLocation()

	retv := &Host{}
	retv.Name = nameGenerator.Generate()
	retv.Continent = strings.ToLower(items[0])
	retv.Region = strings.ToLower(items[1])
	retv.City = strings.ToLower(items[2])
	retv.Team = RandomFromList(Teams)
	retv.Purpose = RandomFromList(Purpose)

	return retv
}

func writeHost(w *os.File, host *Host) {
	fmt.Fprintf(w, "    %s:\n", host.Name)
	fmt.Fprintf(w, "      continent: %s\n", host.Continent)
	fmt.Fprintf(w, "      region: %s\n", host.Region)
	fmt.Fprintf(w, "      city: %s\n", host.City)
	fmt.Fprintf(w, "      team: %s\n", host.Team)
	fmt.Fprintf(w, "      purpose: %s\n", host.Purpose)
}

func writeHostsFile(outputfile string, hosts []Host) {
	w, err := os.Create(outputfile)
	if err != nil {
		fmt.Println("Error: Cannot create file")
		return
	}
	defer w.Close()

	LogIfWriteError(w.WriteString("---\n"))
	LogIfWriteError(w.WriteString("all:\n"))
	LogIfWriteError(w.WriteString("  hosts:\n"))
	for _, host := range hosts {
		writeHost(w, &host)
	}
	continentList := groupHostsByContinent(hosts)
	regionList := groupHostsByRegion(hosts)
	cityList := groupHostsByCity(hosts)
	teamList := groupHostsByTeam(hosts)
	purposeList := groupHostsByPurpose(hosts)

	writeGroup(w, continentList, 10)
	writeGroup(w, regionList, 20)
	writeGroup(w, cityList, 30)
	writeGroup(w, teamList, 100)
	writeGroup(w, purposeList, 110)

}

func main() {
	hosts := []Host{}
	NumberOfHosts := 5
	var err error

	if len(os.Args) > 1 {
		firstArg := os.Args[1]
		NumberOfHosts, err = strconv.Atoi(firstArg)
		if err != nil {
			fmt.Println("Error: Argument is not a valid integer")
			return
		}
	}

	for i := 1; i < NumberOfHosts; i++ {
		r := generateHost()
		hosts = append(hosts, *r)
	}
	writeHostsFile("hosts.yml", hosts)

	printMkDirPerListItem("t", Teams)
	printMkDirPerListItem("p", Purpose)
	printMkDirPerLocation()
}
