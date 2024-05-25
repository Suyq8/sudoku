package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sudoku/pkg/sudoku"

	"google.golang.org/grpc"
)

// Usage String
const USAGE_STRING = "go run main.go -p 8081 -d"

// Exit codes
const EX_USAGE int = 64

func main() {
	// Custom flag Usage message
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of %s:\n", USAGE_STRING)
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "  -%s: %v\n", f.Name, f.Usage)
		})
	}

	// Parse command-line argument flags
	port := flag.Int("p", 8080, "(default = 8080) Port to accept connections")
	debug := flag.Bool("d", false, "Output log statements")
	flag.Parse()

	addr := "localhost:" + strconv.Itoa(*port)

	// Disable log outputs if debug flag is missing
	if !(*debug) {
		log.SetFlags(0)
		log.SetOutput(io.Discard)
	}

	log.Fatal(startServer(addr))
}

func startServer(hostAddr string) error {
	grpcServer := grpc.NewServer()
	srvSolver := sudoku.NewSolver()
	sudoku.RegisterSudokuServer(grpcServer, srvSolver)

	listener, err := net.Listen("tcp", hostAddr)
	if err != nil {
		return err
	}

	if err = grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}
