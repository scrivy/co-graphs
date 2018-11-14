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
	case "rrdupdate":
		if len(os.Args) != 4 {
			cliHelp()
			os.Exit(1)
		}
		rrdUpdate(os.Args[2], os.Args[3])
	default:
		cliHelp()
	}
}

func cliHelp() {
	fmt.Println(`
wrong

rrdupdate [.rrd file] [.csv file]`)
}

const timeFormat = "2006-01-02 15:04:05"

func rrdUpdate(rrdFilename, csvFilename string) {
	csvFile, err := os.Open(csvFilename)
	if err != nil {
		panic(err)
	}
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
		err = exec.Command("rrdtool", "update", rrdFilename, update).Run()
		if err != nil {
			fmt.Println(update, err)
		}
	}
}
