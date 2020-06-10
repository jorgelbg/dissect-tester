// Copyright (c) 2020 Jorge Luis Betancourt. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/elastic/beats/libbeat/processors/dissect"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		pattern = kingpin.Flag("pattern", "Tokenizer pattern to test the samples.").Required().Short('p').String()
	)

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1").Author("Jorge Luis Betancourt")
	kingpin.CommandLine.Help = "Tool for testing a set of sample loglines against a dissect pattern."
	kingpin.Parse()

	processor, err := dissect.New(*pattern)
	if err != nil {
		fmt.Printf("ERROR: Coudln't create the processor: %s", err.Error())
		return
	}

	fmt.Println("Enter a sample per line to evaluate")
	reader := bufio.NewReader(os.Stdin)
	for {
		sample, _ := reader.ReadString('\n')
		sample = strings.TrimSuffix(sample, "\n")

		if sample == "" {
			return
		}

		tokens, err := processor.Dissect(sample)
		if err != nil {
			fmt.Printf("ERROR: %s", err.Error())
			continue
		}

		payload, err := json.Marshal(tokens)
		if err != nil {
			fmt.Printf("ERROR: Could not serialize the list of tokens into JSON: %s",
				err.Error())
		}
		fmt.Printf("%v\n", string(payload))
	}
}
