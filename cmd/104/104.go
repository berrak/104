/*
Copyright 2018 The 104 Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	grt "github.com/berrak/pkg103greet"
	spe "github.com/bgentry/speakeasy"
)

/*
The purpose of this example is:
- Show usage of multiple files used in main package
- Import a previous built package, which is installed on build host, but not in Debian official stretch release
- Import a Debian golang package (speakeasy) and use it. This is not installed on build host.
- Debanize this application, see docs folder.
To run:
- go get -u -v github.com/bgentry/speakeasy
- go get -u -v github.com/berrak/pkg103greet
- Run this application with 'cd cmd/104 && go run 104.go one.go'
*/
func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s. %s", getOne(), "Enter your name: ")
	yourName, _ := reader.ReadString('\n')

	// Randomize source to different greeting each run
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	i := r1.Intn(5)

	// Greet user with random languages 'Hello'
	langCode := [5]string{"en", "it", "fr", "po", "ru"}
	greeting := grt.GreetMe(langCode[i])
	fmt.Printf("%s %s", greeting, yourName)

	// Ask password of this user
	password, err := spe.Ask("2. Please enter a password: ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Password result: %q\n", password)

}
