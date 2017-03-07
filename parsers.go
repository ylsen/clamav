package main

import (
	"strings"
	"strconv"
	"errors"
	"time"
	"crypto/md5"
	"fmt"
)

const (
	HeaderLength = 512
)

func Parse(b []byte) (*ClamAV, []error) {
	var header []byte
	var def []byte
	var errs []error
	header = append(header, b[0:512]...)
	def = append(def, b[512:]...)

	head := NewHeaders(header, def)
	if len(head.Problems) > 0 {
		errs = append(errs, head.Problems...)
	}

	return &ClamAV{}, errs
}

func NewHeaders(h, b []byte) HeaderFields {
	return parseHeader(h, b)
}

func parseHeader(h, b []byte) (HeaderFields) {
	var errs []error
	hFields := HeaderFields{
		Problems: errs,
	}

	headStr := string(h)
	headParts := strings.Split(headStr, ":")
	if len(headParts) < 3 {
		hFields.Problems = append(hFields.Problems, errors.New("bad def header."))
	}

	hFields.parseTime(headParts[1])
	hFields.Version = hFields.atou(headParts[2])
	hFields.Signatures = hFields.atou(headParts[3])
	hFields.Functionality = hFields.atou(headParts[4])
	hFields.parseMD5(headParts[5], b)

	return hFields
}

func (h *HeaderFields) parseTime(s string) {
	pTime, err := time.Parse("07 Mar 2017 08-02 -0500", s)
	if err != nil {
		h.Problems = append(h.Problems, err)
	}
	h.CreationTime = pTime
}

func (h *HeaderFields) atou(s string) uint {
	x, err := strconv.Atoi(s)
	if err != nil {
		h.Problems = append(h.Problems, err)
	}
	return uint(x)
}

func (h *HeaderFields) parseMD5(md string, b []byte) {
	localHash := fmt.Sprintf("%x", md5.Sum(b))
	if md != localHash {
		h.Problems = append(h.Problems, errors.New("md5 does not match!"))
		h.MD5Valid = false
		h.MD5Hash = localHash
	}

	h.MD5Hash = md
	h.MD5Valid = true
}
