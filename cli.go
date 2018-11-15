package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		cliHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "update":
		if len(os.Args) != 4 {
			cliHelp()
			os.Exit(1)
		}
		update(os.Args[2], os.Args[3])
	case "graph":
		if len(os.Args) != 4 {
			cliHelp()
			os.Exit(1)
		}
		graph(os.Args[2], os.Args[3])
	default:
		cliHelp()
	}
}

func cliHelp() {
	fmt.Println(`
wrong

update [.rrd file] [.csv file]
graph [.rrd file] [24 hour start 2006-01-02T15:04:05Z07:00]`)
}

const timeFormat = "2006-01-02 15:04:05"

func update(rrdFilename, csvFilename string) {
	csvFile, err := os.Open(csvFilename)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	scanner := bufio.NewScanner(csvFile)
	for scanner.Scan() {
		record := strings.Split(scanner.Text(), ",")
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		if len(record) < 3 {
			fmt.Println("not enough columns")
			continue
		}
		if record[0] == "EasyLog USB" { // skip header
			continue
		}
		recordedAtStr := record[1]
		value := record[2]
		recordedAt, err := time.Parse(timeFormat, recordedAtStr)
		if err != nil {
			panic(err)
		}
		recordedAt = recordedAt.Add(time.Hour * 8) // adjust to utc

		update := fmt.Sprintf("%d:%s", recordedAt.Unix(), value)
		//		fmt.Println(update)
		output, err := exec.Command("rrdtool", "update", rrdFilename, update).CombinedOutput()
		if len(output) > 0 {
			fmt.Printf("%s\n", output)
		}
		if err != nil {
			fmt.Println(update, err)
		}
	}
}

func graph(rrdFilename, timeStr string) {
	start, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		panic(err)
	}
	outputFilename := start.Format("2006_01_02") + ".png"
	output, err := exec.Command(
		"rrdtool", "graphv", outputFilename,
		"DEF:co_ppm="+rrdFilename+":co_ppm:AVERAGE",
		"--start", fmt.Sprintf("%d", start.Unix()),
		"--end", "start+24h",
		"LINE1:co_ppm#0000FF:co (ppm)",
		"-w", "1200", "-h", "400",
	).CombinedOutput()
	if len(output) > 0 {
		fmt.Printf("%s\n", output)
	}
	if err != nil {
		panic(err)
	}
}
