package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pierreyves258/dl24/px100"
)

func main() {
	var tty, file, delimiter string

	flag.StringVar(&tty, "p", "/dev/ttyUSB0", "Serial port")
	flag.StringVar(&file, "o", "/dev/stdout", "CSV file")
	flag.StringVar(&delimiter, "d", ",", "CSV delimiter")
	flag.Parse()

	dl24, err := px100.NewPX100(tty)
	if err != nil {
		log.Fatal(err)
	}
	defer dl24.Destroy()

	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = f.WriteString("time" + delimiter + "voltage" + delimiter + "current\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = dl24.SetData(px100.Reset, nil)
	if err != nil {
		log.Fatalln(err)
	}

	err = dl24.SetData(px100.SetCurrent, 1.65)
	if err != nil {
		log.Fatalln(err)
	}

	err = dl24.SetData(px100.SetCutoff, 5.5)
	if err != nil {
		log.Fatalln(err)
	}

	err = dl24.SetData(px100.SetOutput, true)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		voltage, err := dl24.GetData(px100.GetVoltage)
		if err != nil {
			fmt.Println(err)
			continue
		}

		current, err := dl24.GetData(px100.GetCurrent)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ison, err := dl24.GetData(px100.GetIsOn)
		if err != nil {
			continue
		}

		if !ison.(bool) {
			return
		}

		str := fmt.Sprintf("%s%s%f%s%f\n", time.Now().Format("2006-01-02 15:04:05"), delimiter, voltage, delimiter, current)

		_, err = f.WriteString(str)
		if err != nil {
			fmt.Println("ERROR WRITE", err)
		}

		time.Sleep(1 * time.Second)
	}
}
