/*
The MIT License (MIT)

Copyright (c) 2014-2017 DutchCoders [https://github.com/dutchcoders/]
Copyright (c) 2018-2020 Andrea Spacca.
Copyright (c) 2020- Andrea Spacca and Stefan Benten.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package server

import (
	// _ "transfer.sh/app/handlers"
	// _ "transfer.sh/app/utils"

	"fmt"
	"io"
	"net/http"
	"time"

	clamd "github.com/dutchcoders/go-clamd"

	"github.com/gorilla/mux"
)

func (s *Server) scanHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	filename := sanitize(vars["filename"])

	contentLength := r.ContentLength
	contentType := r.Header.Get("Content-Type")

	s.logger.Printf("Scanning %s %d %s", filename, contentLength, contentType)

	var reader io.Reader

	reader = r.Body

	c := clamd.NewClamd(s.ClamAVDaemonHost)

	abort := make(chan bool)
	defer close(abort)
	response, err := c.ScanStream(reader, abort)
	if err != nil {
		s.logger.Printf("%s", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	select {
	case s := <-response:
		w.Write([]byte(fmt.Sprintf("%v\n", s.Status)))
	case <-time.After(time.Second * 60):
		abort <- true
	}
}
