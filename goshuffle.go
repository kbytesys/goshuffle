package main

import "errors"
import "fmt"
import "flag"
import "os"
import (
	"./utils"
)

func get_device_directory(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("I need the device dir!!! Use me with -h for help!")
	}

	src, err := os.Stat(args[0])
	if err != nil {
		return "", err
	}

	if src.IsDir() {
		return args[0], nil
	} else {
		return "", errors.New("The destination isn't a directory!")
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "    %s [flags] device_directory \n\n", "goshuffle")
		fmt.Fprintf(os.Stderr, " Flag lists:\n")
		flag.PrintDefaults()
	}

	initPrt := flag.Bool("init", false, "init the directory/device for shuffle")
	randomizePrt := flag.Bool("randomize", false, "randomize the files into the device")

	flag.Parse()

	var device_dir, err = get_device_directory(flag.Args())

	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", err)
		os.Exit(1)
	}

	if *initPrt {
		utils.InitDevice(device_dir)
	}

	if *randomizePrt {
		utils.Randomize(device_dir)
	}
}
