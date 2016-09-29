package parser

import (
	"bytes"
	"regexp"
	"strings"
)

// There can be 3 types of parsers for MySQL output
/*
Type 1: with \G or with `mysql --auto-vertical-output=true` or `mysql -E`
starts with ***** - has 27 and then no. row and then 27 stars
*************************** 10. row ***************************
   key: value
otherk: fkldasfd asflsa
   key: value
*************************** 1. row ***************************
...

all key names end at the same location and values start at the location
*/
/*
Type 2: regular mysql output copy/pasted or with `mysql -t`
starts and ends with the same line

contains with -+ and s
+-----------------------------------------+---------------------+-------------+--------------+---------------+
+ indicates start of column. till end of line
the column lengths are unique

*/
/*
Type 3: tab separated columns.

*/
/*
Type 4: HTML with table <TABLE ...> containing data
*/

const sniffLen = 1024

type sniffSig interface {
	// match returns the MIME type of the data, or "" if unknown.
	match(data []byte, firstNonWS int) PatternType
}

// DetectType ..
func DetectType(data []byte) PatternType {
	if len(data) > sniffLen {
		data = data[:sniffLen]
	}

	firstNonWS := 0
	for ; firstNonWS < len(data) && isWS(data[firstNonWS]); firstNonWS++ {
	}

	for _, sig := range sniffSignatures {
		if ct := sig.match(data, firstNonWS); ct != ptnNoMatch {
			return ct
		}
	}

	return ptnNoMatch // fallback
}

type strStartSig struct {
	sig []byte
	ct  PatternType
}
type strMiscSig struct {
	strType string
	ct      PatternType
}

var sniffSignatures = []sniffSig{
	&strStartSig{[]byte("<TABLE"), ptnHTMLTable},
	&strStartSig{[]byte("<table"), ptnHTMLTable},
	&strStartSig{[]byte("*************************"), ptnMultiLineSQL},
	&strMiscSig{"ascii", ptnASCIITable},
	&strMiscSig{"tabbed", ptnTabbedTable},
}

func (h strStartSig) match(data []byte, _firstNonWS int) PatternType {
	// check starts with
	if bytes.HasPrefix(data, h.sig) {
		return h.ct
	}
	return ptnNoMatch
}

const asciiTableRegexp = `^\+((\-)+\+)+`
const tabbedRegExp = `(\S+\t)*\S+`

func (h strMiscSig) match(data []byte, _firstNonWS int) PatternType {
	var lines = strings.Split(string(data), "\n")
	switch h.strType {
	case "ascii":
		// start with +, have one or more -, ends with +
		if len(lines) > 3 && matchRegexp(asciiTableRegexp, lines[0], lines[2]) {
			return h.ct
		}
	case "tabbed":
		// fmt.Println(matchRegexp(tabbedRegExp, lines[0]))
		if len(lines) > 2 && matchRegexp(tabbedRegExp, lines[0], lines[1]) {
			return h.ct
		}
	}
	return ptnNoMatch
}

func matchRegexp(reg string, data ...string) bool {
	var re = regexp.MustCompile(reg)
	isFalse := true
	for _, str := range data {
		var match = re.FindStringSubmatch(str)
		if len(match) > 0 {
			// fmt.Println(match)
		} else {
			isFalse = false
			break
		}
	}
	return isFalse
}

func isWS(b byte) bool {
	switch b {
	case '\t', '\n', '\x0c', '\r', ' ':
		return true
	}
	return false
}
