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
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func main() {
	http.HandleFunc("/", handler)
	log.Printf("listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if name, ok := isGopher(r.URL.Path[1:]); ok {
		fmt.Fprintf(w, "hello, %s", name)
		return
	}
	fmt.Fprintln(w, "hello, stranger")
}

func isGopher(email string) (string, bool) {
	re := regexp.MustCompile("^([[:alpha:]]+)@golang.org$")
	match := re.FindStringSubmatch(email)
	if len(match) == 2 {
		return match[1], true
	}
	return "", false
}
