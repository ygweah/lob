package main

import (
	"flag"
	"lob/api"
	"log"
)

func main() {
	addrFlag := flag.String("addr", "", "host:port")
	fileFlag := flag.String("file", "", "data file in JSON")
	flag.Parse()

	log.Printf("addr: %v, file: %v", *addrFlag, *fileFlag)

	adrs, err := api.LoadFromFile(*fileFlag)
	if err != nil {
		log.Fatalf("cannot load file: %v", err)
	}
	log.Printf("loaded address count: %d", len(adrs))

	m := api.NewSimpleAddressManager(adrs)
	s := api.NewAddressService(*addrFlag, m)
	if err := s.Start(); err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
	s.Wait()
}
