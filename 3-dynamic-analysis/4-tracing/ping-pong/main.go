// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"os"
	"runtime/trace"
	"time"
)

func main() {
	_ = trace.Start(os.Stdout)
	defer trace.Stop()

	table := make(chan int)
	go player(table, "ping")
	go player(table, "pong")

	table <- 0
	time.Sleep(time.Second)
	ball := <-table
	close(table)
	log.Printf("played %d turns", ball)

	// runtime.GC()
	// _ = pprof.WriteHeapProfile(os.Stdout)
}

func player(table chan int, name string) {
	for ball := range table {
		log.Printf("%d\t%s", ball, name)
		table <- ball + 1
	}
}
