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

	kingpin.Parse()
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1").Author("Jorge Luis Betancourt")
	kingpin.CommandLine.Help = "Tool for testing a set of sample loglines against a dissect pattern."

	if *pattern == "" {
		os.Exit(1)
	}

	fmt.Printf("pattern = %+v\n", *pattern)
	processor, err := dissect.New(*pattern)
	if err != nil {
		panic(err)
	}

	fmt.Println("🔎Enter a sample per line to evaluate against the pattern")
	reader := bufio.NewReader(os.Stdin)
	for {
		sample, _ := reader.ReadString('\n')
		sample = strings.TrimSuffix(sample, "\n")

		if sample == "" {
			return
		}

		tokens, err := processor.Dissect(sample)
		if err != nil {
			panic(err)
		}

		payload, err := json.Marshal(tokens)

		fmt.Printf("%v\n", string(payload))
	}
}
