/* otto Source Code
 * Copyright (C) 2013 Lumen LLC. 
 *
 * otto is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * 
 * otto is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 * 
 * You should have received a copy of the GNU Affero General Public License
 * along with otto. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	lmnhtml "lumenlearning.com/util/html"
	lmnhttp "lumenlearning.com/util/http"
	lmnweb "lumenlearning.com/util/canvas/web"
)

var username *string = flag.String("username", "somefakeuser@email", "The username of the Canvas user account we'll be using for scraping.")
var password *string = flag.String("password", "supersecretpass!", "The password for the user account.")
var idFile *string = flag.String("idfile", "", "The path to the file that contains the iDs of the courses and assignments we want to find outcomes for. The file is comma-separated (course_id,assignment_id).")
var login1 *string = flag.String("login1", "https://lumen.instructure.com/login", "The first URL to load that will give us the needed authenticy_token information for the actual login.")
var login2 *string = flag.String("login2", "https://lumen.instructure.com/login?nonldap=true", "The second URL to POST to in order to complete the login process.")
var userAgentString *string = flag.String("useragent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_8) AppleWebKit/537.13+ (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2", "The User-Agent string to send in all HTTP requests.")

func main() {
	// Get the command line arguments
	flag.Parse()

	// Open up a Canvas web session using the supplied credentials. We'll use
	//   the resulting UserAgentClient for all of our subsequent page pulls.
	sessionClient, err := lmnweb.CanvasWebLogin(*username, *password, *login1, *login2, *userAgentString)
	if err != nil {
		log.Fatalf("lmnweb.CanvasWebLogin => %v", err.Error())
	}

	// Read in the file containing the course and assignment IDs that we'll
	//   need to go look at. Expecting one course_id and assignment_id per line, separated by a comma.  Input file must not have a header line.  Example: 
	//     12345,34967
	//     12345,35102
	//     23456,18182
	//     24243,67182
	//     etc.

	// Print the output CSV header line
	fmt.Println("\"course_id\",\"assignment_id\",\"outcome_id\"")

	// Open the file that has all the course and assignment IDs
	idFile, err := os.Open(*idFile)
	if err != nil {
		log.Fatalf("os.Open => %v", err.Error())
	}
	idScan := bufio.NewScanner(idFile)
	for idScan.Scan() {
		// Parse the course and assignment IDs from the line of text
		line := idScan.Text()
		ids := strings.Split(line, ",")
		courseId := ids[0]
		asnId := ids[1]

		// Create the URL for this assignment
		url := fmt.Sprintf("https://lumen.instructure.com/courses/%v/assignments/%v/rubric", courseId, asnId)

		// Fetch the URL using our session client
		log.Printf("GET %v\n", url)
		pageContent, err := lmnhttp.GetPageContentWithClient(url, sessionClient)
		if err != nil {
			log.Fatalf("lmnhttp.GetPageContent => %v", err.Error())
		}

		// We find the outcome IDs in the HTML from each page by looking for
		//   elements matching the outcome selector.
		outcomeSelector := "span.learning_outcome_id"
		outcomeNodes, err := lmnhtml.FindNodes(pageContent, outcomeSelector)
		if err != nil {
			log.Fatalf("lmnhtml.FindNodes => %v", err.Error())
		}

		if len(outcomeNodes) > 0 {
			for _, o := range outcomeNodes {
				// Write data to stdout as CSV
				outcomeId := lmnhtml.GetNodeText(o)
				outcomeId = strings.TrimSpace(outcomeId)

				if len(outcomeId) > 0 {
					fmt.Printf("%v,%v,%v\n", courseId, asnId, outcomeId)
				}
			}
		} else {
			log.Printf("courses/%v/assignments/%v/rubric has no outcomes.\n", courseId, asnId)
		}
	}
}
