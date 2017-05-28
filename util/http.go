// Copyright 2017 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"gopkg.in/xmlpath.v2"
)

var errElementNotFound = errors.New("element not found in document")

func HTTPGet(link string) ([]byte, error) {
	var buf []byte
	cj, err := cookiejar.New(nil)
	tr := &http.Transport{
		TLSHandshakeTimeout:   20 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Jar:       cj,
	}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	// arbitrairly chosen 5MiB
	if resp.ContentLength > 5*1024*1024 || resp.ContentLength < 0 {
		buf = make([]byte, 5*1024*1024)
	} else if resp.ContentLength == 0 {
		return []byte{}, nil
	} else {
		buf = make([]byte, resp.ContentLength)
	}

	i, err := io.ReadFull(resp.Body, buf)
	if err == io.ErrUnexpectedEOF {
		buf = buf[:i]
	} else if err != nil {
		return []byte{}, err
	}

	return buf, nil
}

func HTTPGetXpath(link, xpathStr string) (string, error) {
	cj, err := cookiejar.New(nil)
	tr := &http.Transport{
		TLSHandshakeTimeout:   20 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Jar:       cj,
	}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	lr := io.LimitReader(resp.Body, 5*1024*1024)

	root, err := xmlpath.ParseHTML(lr)
	if err != nil {
		return "", err
	}

	xpath := xmlpath.MustCompile(xpathStr)
	if value, ok := xpath.String(root); ok {
		return value, nil
	}

	return "", errElementNotFound
}
