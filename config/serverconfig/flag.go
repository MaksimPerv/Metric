package serverconfig

import "flag"

var FlagRunAddr string

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")

	flag.Parse()
}
