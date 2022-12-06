package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	inputtext := ""
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)

		for {
			line, err := reader.ReadString('\n')

			if err != nil {

				break
			}

			inputtext += line
		}

	} else {
		inputstring := flag.String("html", "", "use --html ")
		flag.Parse()
		if *inputstring == "" {
			println("please use --html to set file")
			return
		} else {

			inputtext = string(*inputstring)
		}
	}
	justpoint := flag.Bool("point", false, "use --point ")
	flag.Parse()
	inputtext = strings.ReplaceAll(strings.ToLower(inputtext), "\"", "'")
	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err.Error())
	}
	exPath := filepath.Dir(ex)

	configfilepath := exPath + "/scripts/"
	files, err := ioutil.ReadDir(configfilepath)
	if err != nil {
		log.Fatal(err)
	}
	out := []outpoint{}
	for _, file := range files {
		//fmt.Println(file.Name(), file)
		jsonFile, err := os.Open(configfilepath + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)

		// we initialize our Users array
		paths := moduls{}

		// we unmarshal our byteArray which contains our
		xml.Unmarshal(byteValue, &paths)

		switch paths.Module.Type {
		case "contains":
			for _, v := range paths.Module.Values {
				matchcount := iscontaions(inputtext, v.Value)
				if matchcount != 0 {
					//fmt.Printf("%d %s Found \n", matchcount, paths.Module.Title)
					add := outpoint{}
					add.Match = paths.Module.Title
					add.MatchCount = matchcount
					add.Module = paths.Module
					out = append(out, add)
				} else {
					//fmt.Printf("%s Not Found \n", paths.Title)
				}
			}

		case "regex":
			for _, v := range paths.Module.Values {
				rec := isregex(inputtext, v.Value)
				if len(rec) > 0 {

					//fmt.Printf("%s %s \n", paths.Module.Title, v)
					add := outpoint{}
					add.Match = paths.Module.Title
					add.MatchCount = len(rec)
					add.Module = paths.Module

					for _, v := range rec {

						add.MtachValue = add.MtachValue + "[ " + v + " ]"
					}
					out = append(out, add)
				} else {
					//fmt.Printf("%s Not Found \n", paths.Title)
				}
			}

		}

		defer jsonFile.Close()

	}
	printwithpoint(out, *justpoint)
}

func printwithpoint(outpoint []outpoint, showpoint bool) {
	point := 0
	for _, v := range outpoint {
		point += v.Module.Point * v.MatchCount
		if !showpoint {
			fmt.Printf("%s %s|  %s | matchcount: %d point : %d\n", v.Match, v.Module.Name, v.MtachValue, v.MatchCount, v.Module.Point*v.MatchCount)

		}

	}
	if !showpoint {
		fmt.Printf("%s %d", "page point ", point)
	} else {
		fmt.Printf("%d", point)
	}

}

func iscontaions(input string, match string) int {

	Compare_exsist := strings.Count(input, match)
	return Compare_exsist
}
func isregex(input string, match string) []string {

	Compare_regex := regexp.MustCompile(match)
	matches := Compare_regex.FindAllStringSubmatch(input, -1)
	retxx := []string{}
	for _, v := range matches {
		retxx = append(retxx, v[0])
	}
	return retxx
}

type module struct {
	Type   string  `xml:"type"`
	Name   string  `xml:"name"`
	Title  string  `xml:"title"`
	Point  int     `xml:"point"`
	Values []value `xml:"values"`
}
type value struct {
	Value string `xml:"value>text"`
}

type moduls struct {
	XMLName xml.Name `xml:"modules"`
	Module  module   `xml:"module"`
}

type outpoint struct {
	Match      string
	MatchCount int
	Module     module
	MtachValue string
}
