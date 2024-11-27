package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pierreyves258/atorch"
)

func main() {
	var tty, file, delimiter string
	var current, voltage float64

	flag.StringVar(&tty, "p", "/dev/ttyUSB0", "Serial port")
	flag.StringVar(&file, "o", "/dev/stdout", "CSV file")
	flag.StringVar(&delimiter, "d", ",", "CSV delimiter")
	flag.Float64Var(&voltage, "v", 5.5, "Voltage cut off")
	flag.Float64Var(&current, "c", 1.65, "Load current")

	flag.Parse()

	dl24, err := atorch.NewPX100(tty)
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

	_, err = f.WriteString("time" + delimiter + "voltage" + delimiter + "current" + delimiter + "capacity\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = dl24.SetData(atorch.Reset, nil, true)
	if err != nil {
		log.Fatalln(err)
	}

	err = dl24.SetData(atorch.SetCurrent, current, false)
	if err != nil {
		log.Fatalln(err)
	}

	err = dl24.SetData(atorch.SetCutoff, voltage, false)
	if err != nil {
		log.Fatalln(err)
	}

	err = dl24.SetData(atorch.SetOutput, true, false)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		voltage, err := dl24.GetData(atorch.GetVoltage)
		if err != nil {
			fmt.Println(err)
			continue
		}

		current, err := dl24.GetData(atorch.GetCurrent)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ison, err := dl24.GetData(atorch.GetIsOn)
		if err != nil {
			continue
		}

		capacity, err := dl24.GetData(atorch.GetCharge)
		if err != nil {
			continue
		}

		if !ison.(bool) {
			fmt.Printf("Capacity %f\n", capacity)
			return
		}

		str := fmt.Sprintf("%s%s%f%s%f%s%f\n", time.Now().Format("2006-01-02 15:04:05"), delimiter, voltage, delimiter, current, delimiter, capacity)

		_, err = f.WriteString(str)
		if err != nil {
			fmt.Println("ERROR WRITE", err)
		}

		time.Sleep(1 * time.Second)
	}
}
