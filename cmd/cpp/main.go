package main

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	cmd "github.com/bhojpur/build/cmd/cpp/commands"
	"github.com/tj/go-spin"
)

const logo = `Bhojpur Build - C/C++ to Go source code interface engine
Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.
All rights reserved.

See https://github.com/bhojpur/build for more examples and documentation.
`

func init() {
	if *cmd.Debug {
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetFlags(0)
	}
	flag.Usage = func() {
		fmt.Println(logo)
		fmt.Printf("Usage: buildc2go package1.yml [package2.yml] ...\n\n")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		fmt.Println()
		log.Fatalln("[ERR] no package configuration files have been provided.")
	}
}

func main() {
	s := spin.New()

	var wg sync.WaitGroup
	doneChan := make(chan struct{})
	for _, cfgPath := range getConfigPaths() {
		if *cmd.Fancy {
			wg.Add(1)
			go func() {
				for {
					select {
					case <-doneChan:
						doneChan = make(chan struct{})
						fmt.Printf("\r  \033[36mprocessing %s\033[m done.\n", cfgPath)
						wg.Done()
						return
					default:
						fmt.Printf("\r  \033[36mprocessing %s\033[m %s", cfgPath, s.Next())
						time.Sleep(100 * time.Millisecond)
					}
				}
			}()
		}

		var t0 time.Time
		if *cmd.Debug {
			t0 = time.Now()
		}
		process, err := cmd.NewProcess(cfgPath, *cmd.OutputPath)
		if err != nil {
			log.Fatalln("[ERR]", err)
		}
		process.Generate(*cmd.NoCGO)
		if err := process.Flush(*cmd.NoCGO); err != nil {
			log.Fatalln("[ERR]", err)
		}
		if *cmd.Debug {
			fmt.Printf("done in %v\n", time.Now().Sub(t0))
		}
		if *cmd.Fancy {
			close(doneChan)
			wg.Wait()
		}
	}
}

func getConfigPaths() (paths []string) {
	for _, path := range flag.Args() {
		if info, err := os.Stat(path); err != nil {
			log.Fatalln("[ERR] cannot locate the specified path:", path)
		} else if info.IsDir() {
			if path, ok := configFromDir(path); ok {
				paths = append(paths, path)
				continue
			}
			log.Fatalln("[ERR] cannot find any config file in:", path)
		}
		paths = append(paths, path)
	}
	return
}

func configFromDir(path string) (string, bool) {
	possibleNames := []string{"c2go.yaml", "c2go.yml"}
	if base := filepath.Base(path); len(base) > 0 {
		possibleNames = append(possibleNames,
			fmt.Sprintf("%s.yaml", base), fmt.Sprintf("%s.yml", base))
	}
	for _, name := range possibleNames {
		path := filepath.Join(path, name)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, true
		}
	}
	return "", false
}
