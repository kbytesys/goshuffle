package utils

import "os"
import "errors"
import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"
)

type randomize_helper struct {
	targetdir string
}

func (r *randomize_helper) walk(path string, info os.FileInfo, err error) error {
	if strings.Contains(path, ".goshuffle") {
		return filepath.SkipDir
	}

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return nil
	}

	if err != nil {
		return err
	}

	if filepath.Ext(path) == ".mp3" || filepath.Ext(path) == ".MP3" {
		_, filename := filepath.Split(path)

		match, _ := regexp.MatchString("^s[0-9]{5}-", filename)
		if match {
			filename = filename[7:len(filename)]
		}

		filename = generate_random_prefix() + filename

		err := os.Rename(path, r.targetdir+string(os.PathSeparator)+filename)

		if err != nil {
			return err
		}

	}

	return nil
}

func generate_random_prefix() string {
	var letters = []rune("0123456789")

	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("s%s-", string(b))
}

func cleanwalk(path string, info os.FileInfo, err error) error {
	if strings.Contains(path, ".goshuffle") {
		return filepath.SkipDir
	}

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {

		fis, err := ioutil.ReadDir(path)

		if err != nil {
			return err
		}

		for _, child := range fis {
			if child.IsDir() {
				err = filepath.Walk(path+string(os.PathSeparator)+child.Name(), cleanwalk)

				if err != nil {
					return err
				}
			}
		}

		fis, err = ioutil.ReadDir(path)

		if err != nil {
			return err
		}

		if len(fis) == 0 {
			err = os.Remove(path)
		}
		if err != nil {
			return err
		}

		return filepath.SkipDir
	}

	return nil
}

func InitDevice(dir string) {
	// check if the dir contains already the .goshuffle directory
	sdir, err := os.Stat(dir + string(os.PathSeparator) + ".goshuffle")
	if err != nil {
		if os.IsNotExist(err) {
			create_err := os.Mkdir(dir+string(os.PathSeparator)+".goshuffle", 755)
			if create_err != nil {
				panic(create_err)
			}

			fmt.Fprintf(os.Stdout, "Device at directory %s initialized, exiting...\n", dir)
			os.Exit(0)
		}
	} else {
		fmt.Fprintf(os.Stdout, "Device at directory %s already initialized, exiting...\n", dir)
		os.Exit(0)
	}

	if !sdir.IsDir() {
		panic(errors.New("Can't create the control directory!"))
	}
}

func Randomize(dir string) {
	// check if the device is initialized
	_, err := os.Stat(dir + string(os.PathSeparator) + ".goshuffle")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stdout, "Device at directory %s not initialized, exiting...\n", dir)
			os.Exit(0)
		}
		panic(err)
	}

	var rh = new(randomize_helper)
	rh.targetdir = dir

	err = filepath.Walk(dir, rh.walk)
	if err != nil {
		panic(err)
	}

	// clean empty directories
	err = filepath.Walk(dir, cleanwalk)
	if err != nil {
		panic(err)
	}

}
