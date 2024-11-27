package main

import (
	"flag"

	"github.com/pierreyves258/atorch"
	"github.com/pierreyves258/atorch/cmd/dl24_api/router"
	"github.com/rs/zerolog/log"
)

func main() {
	var tty string

	flag.StringVar(&tty, "p", "/dev/serial/by-path/platform-xhci-hcd.1.auto-usb-0:1.1:1.0-port0", "Serial port")

	flag.Parse()

	dl24, err := atorch.NewPX100(tty)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot init atorch load")
	}
	defer dl24.Destroy()

	app := router.Init(dl24)

	app.Run()
}
