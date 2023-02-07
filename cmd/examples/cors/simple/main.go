/*
Copyright 2023 The Kapital Authors

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
	"flag"
	"log"
	"net/http"
)

// Define a string constant containing the HTML for the webpage. This consists of a <h1> // header tag, and some JavaScript which fetches the JSON from our GET /v1/healthcheck // endpoint and writes it to inside the <div id="output"></div> element.
const html = `
<!DOCTYPE html> 
<html lang="en"> 
	<head>
		<meta charset="UTF-8"> 
	</head>
	<body>
		<h1>Simple CORS</h1> 
		<div id="output"></div> 
		<script>
			document.addEventListener('DOMContentLoaded', function() { fetch("http://localhost:4000/v1/healthcheck").then(
			function (response) { response.text().then(function (text) {
			document.getElementById("output").innerHTML = text; });
			}, function(err) {
			document.getElementById("output").innerHTML = err; }
			); });
		</script>
	</body>
</html>`

func main() {
	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()
	log.Printf("starting server on %s", *addr)

	err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))

	log.Fatal(err)
}
