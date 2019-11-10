package main

import "flag"

func main() {
	sourcePort := flag.Int("port", 0, "Source port number")
	dest := flag.String("dest", "", "Destination multiaddr string")
	debug := flag.Bool("debug", false, "Debug generates the same node ID on every execution")
	flag.Parse()

	NewNode().Run(*sourcePort, *dest, *debug)
}
