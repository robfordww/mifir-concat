// MiFIR CONCAT library
// MIT License
// Copyright (c) 2016 robfordww
// "Doesn't look like anything to me"

package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	/* These, and only these are to be removed */
	prefix = []string{"AM", "AUF", "AUF DEM", "AUS DER", "D", "DA", "DE",
		"DE L’", "DEL", "DE LA", "DE LE", "DI", "DO", "DOS", "DU", "IM", "LA",
		"LE", "MAC", "MC", "MHAC", "MHÍC", "MHIC GIOLLA", "MIC", "NI", "NÍ", "NÍC",
		"O", "Ó", "UA", "UI", "UÍ", "VAN", "VAN DE", "VAN DEN", "VAN DER", "VOM",
		"VON", "VON DEM", "VON", "DEN", "VON DER"}

	ccalpha2 = []string{"AD", "AE", "AF", "AG", "AI", "AL", "AM", "AO", "AQ", "AR",
		"AS", "AT", "AU", "AW", "AX", "AZ", "BA", "BB", "BD", "BE", "BF", "BG", "BH",
		"BI", "BJ", "BL", "BM", "BN", "BO", "BQ", "BQ", "BR", "BS", "BT", "BV", "BW",
		"BY", "BZ", "CA", "CC", "CD", "CF", "CG", "CH", "CI", "CK", "CL", "CM", "CN",
		"CO", "CR", "CU", "CV", "CW", "CX", "CY", "CZ", "DE", "DJ", "DK", "DM", "DO",
		"DZ", "EC", "EE", "EG", "EH", "ER", "ES", "ET", "FI", "FJ", "FK", "FM", "FO",
		"FR", "GA", "GB", "GD", "GE", "GF", "GG", "GH", "GI", "GL", "GM", "GN", "GP",
		"GQ", "GR", "GS", "GT", "GU", "GW", "GY", "HK", "HM", "HN", "HR", "HT", "HU",
		"ID", "IE", "IL", "IM", "IN", "IO", "IQ", "IR", "IS", "IT", "JE", "JM", "JO",
		"JP", "KE", "KG", "KH", "KI", "KM", "KN", "KP", "KR", "KW", "KY", "KZ", "LA",
		"LB", "LC", "LI", "LK", "LR", "LS", "LT", "LU", "LV", "LY", "MA", "MC", "MD",
		"ME", "MF", "MG", "MH", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS",
		"MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA", "NC", "NE", "NF", "NG", "NI",
		"NL", "NO", "NP", "NR", "NU", "NZ", "OM", "PA", "PE", "PF", "PG", "PH", "PK",
		"PL", "PM", "PN", "PR", "PS", "PT", "PW", "PY", "QA", "RE", "RO", "RS", "RU",
		"RW", "SA", "SB", "SC", "SD", "SE", "SG", "SH", "SI", "SJ", "SK", "SL", "SM",
		"SN", "SO", "SR", "SS", "ST", "SV", "SX", "SY", "SZ", "TC", "TD", "TF", "TG",
		"TH", "TJ", "TK", "TL", "TM", "TN", "TO", "TR", "TT", "TV", "TW", "TZ", "UA",
		"UG", "UM", "US", "UY", "UZ", "VA", "VC", "VE", "VG", "VI", "VN", "VU", "WF",
		"WS", "YE", "YT", "ZA", "ZM", "ZW"}

	/* Including but not limited to the following titles */
	titles = []string{"ATTY", "COACH", "DAME", "DR", "FR", "GOV", "HONORABLE",
		"MADAM", "MADAME", "MAID", "MASTER", "MISS", "MONSIEUR", "MR", "MRS", "MS",
		"MX", "OFC", "PH.D", "PRES", "PROF", "REV", "SIR"}

	concatCharMap = map[rune]rune{
		0x00C4: 'A', 0x00E4: 'A', 0x00C0: 'A', 0x00E0: 'A', 0x00C1: 'A', 0x00E1: 'A',
		0x00C2: 'A', 0x00E2: 'A', 0x00C3: 'A', 0x00E3: 'A', 0x00C5: 'A', 0x00E5: 'A',
		0x01CD: 'A', 0x01CE: 'A', 0x0104: 'A', 0x0105: 'A', 0x0102: 'A', 0x0103: 'A',
		0x00C6: 'A', 0x00E6: 'A',

		0x00C7: 'C', 0x00E7: 'C', 0x0106: 'C', 0x0107: 'C', 0x0108: 'C', 0x0109: 'C',
		0x010C: 'C', 0x010D: 'C',

		0x010E: 'D', 0x0111: 'D', 0x0110: 'D', 0x010F: 'D', 0x00F0: 'D',

		0x00C8: 'E', 0x00E8: 'E', 0x00C9: 'E', 0x00E9: 'E', 0x00CA: 'E', 0x00EA: 'E',
		0x00CB: 'E', 0x00EB: 'E', 0x011A: 'E', 0x011B: 'E', 0x0118: 'E', 0x0119: 'E',

		0x011C: 'G', 0x011D: 'G', 0x0122: 'G', 0x0123: 'G', 0x011E: 'G', 0x011F: 'G',

		0x0124: 'H', 0x0125: 'H',

		0x00CC: 'I', 0x00EC: 'I', 0x00CD: 'I', 0x00ED: 'I', 0x00CE: 'I', 0x00EE: 'I',
		0x00CF: 'I', 0x00EF: 'I', 0x0131: 'I',

		0x0134: 'J', 0x0135: 'J',

		0x0136: 'K', 0x0137: 'K',

		0x0139: 'L', 0x013A: 'L', 0x013B: 'L', 0x013C: 'L', 0x0141: 'L', 0x0142: 'L', 0x013D: 'L',
		0x013E: 'L',

		0x00D1: 'N', 0x00F1: 'N', 0x0143: 'N', 0x0144: 'N', 0x0147: 'N', 0x0148: 'N',

		0x00D6: 'O', 0x00F6: 'O', 0x00D2: 'O', 0x00F2: 'O', 0x00D3: 'O', 0x00F3: 'O',
		0x00D4: 'O', 0x00F4: 'O', 0x00D5: 'O', 0x00F5: 'O', 0x0150: 'O', 0x0151: 'O',
		0x00D8: 'O', 0x00F8: 'O', 0x0152: 'O', 0x0153: 'O',

		0x0154: 'R', 0x0155: 'R', 0x0158: 'R', 0x0159: 'R',

		0x1E9E: 'S', 0x00DF: 'S', 0x015A: 'S', 0x015B: 'S', 0x015C: 'S', 0x015D: 'S',
		0x015E: 'S', 0x015F: 'S', 0x0160: 'S', 0x0161: 'S', 0x0218: 'S', 0x0219: 'S',

		0x0164: 'T', 0x0165: 'T', 0x0162: 'T', 0x0163: 'T', 0x00DE: 'T', 0x00FE: 'T', 0x021A: 'T',
		0x021B: 'T',

		0x00DC: 'U', 0x00FC: 'U', 0x00D9: 'U', 0x00F9: 'U', 0x00DA: 'U', 0x00FA: 'U',
		0x00DB: 'U', 0x00FB: 'U', 0x0170: 'U', 0x0171: 'U', 0x0168: 'U', 0x0169: 'U', 0x0172: 'U',
		0x0173: 'U', 0x016E: 'U', 0x016F: 'U',

		0x0174: 'W', 0x0175: 'W',

		0x00DD: 'Y', 0x00FD: 'Y', 0x0178: 'Y', 0x00FF: 'Y', 0x0176: 'Y', 0x0177: 'Y',

		0x0179: 'Z', 0x017A: 'Z', 0x017D: 'Z', 0x017E: 'Z', 0x017B: 'Z', 0x017C: 'Z',
	}
	prefixset, ccalpha2set, titlesset map[string]bool
)

func init() {
	// Create fast lookup set dynamically to preserve readability of the
	// literals above
	prefixset = createSet(prefix)
	ccalpha2set = createSet(ccalpha2)
	titlesset = createSet(titles)
}

func createSet(sa []string) map[string]bool {
	set := make(map[string]bool)
	for i := range sa {
		set[sa[i]] = true
	}
	return set
}

func characterRewrite(input string) string {
	// Apply the A-Z untouched, apply charactermap, any other char is deleted
	var b bytes.Buffer
	for _, r := range input {
		if r >= 'A' && r <= 'z' {
			b.WriteRune(r)
		} else {
			c, found := concatCharMap[r]
			if found {
				b.WriteRune(c)
			}
		}
	}
	return b.String()
}

func validateCountryCode(countrycode string) bool {
	return ccalpha2set[strings.ToUpper(countrycode)]
}

func removePrefix(target string, set map[string]bool) string {
	target = strings.ToUpper(target)
	for li := strings.LastIndex(target, " "); li > 0; {
		if set[target[0:li]] {
			return target[li+1:]
		}
		li = strings.LastIndex(target[0:li], " ")
	}
	return strings.TrimSpace(target)
}

func capAndPad(s string) string {
	if len(s) < 5 {
		return s + "#####"[:5-len(s)]
	}
	return s[:5]
}

func removeChars(s string, rs []rune) string {
	return strings.Map(func(sr rune) rune {
		for i := range rs {
			if sr == rs[i] {
				return -1
			}
		}
		return sr
	}, s)
}

// CreateConcat returns a CONCAT string from the input parameters. An error
// returns an empty string and error message
func CreateConcat(countrycode, birthdate, firstname, lastname string) (string, error) {
	if len(firstname) < 1 || len(lastname) < 1 {
		return "", errors.New("Zero length names")
	}
	if !validateCountryCode(countrycode) {
		return "", errors.New("Invalid Country code")
	}
	if _, err := time.Parse("20060102", birthdate); len(birthdate) != 8 || err != nil {
		return "", errors.New("Invalid birthdate format" + err.Error())
	}

	firstname, lastname = strings.TrimSpace(firstname), strings.TrimSpace(lastname)
	// Not specified, but remove first .,; so that titles like "dr.",  "mr."
	// are are removed as well

	punctlist := []rune{'.', ',', ';'}
	firstname = removeChars(firstname, punctlist)
	lastname = removeChars(lastname, punctlist)

	firstname = removePrefix(firstname, titlesset)
	firstname = removePrefix(firstname, prefixset)
	lastname = removePrefix(lastname, titlesset)
	lastname = removePrefix(lastname, prefixset)

	// if firstanme contains more names like "Erwin Rudolf Josef Alexander", split
	// on " " to separate first names. This seems in alignment with specifications
	// that emphasis "first name".
	firstname = strings.Split(firstname, " ")[0]
	firstnamepart := capAndPad(characterRewrite(firstname))
	lastnamepart := capAndPad(characterRewrite(lastname))
	concat := strings.ToUpper(countrycode + birthdate + firstnamepart + lastnamepart)
	return concat, nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of: %s\n", os.Args[0])
		fmt.Fprint(os.Stderr, "\nMifirCONCAT accepts input data on stdin in the form:\n"+
			"ISOCOUNTRYCODE|BIRTHDATE|FIRSTNAME|LASTNAME<LF>\n"+
			"and returns one CONCAT on success, or \"***FAIL***\" on fail\n")
		fmt.Fprintf(os.Stderr, "\nStart %s with --parse to enable stdin reading", os.Args[0])
	}

	flag.Parse()

	stdinreader := bufio.NewReader(os.Stdin)
	for {
		line, err := stdinreader.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		} else if err != nil {
			os.Exit(1)
		}
		if line == "\n" || line == "\r\n" {
			continue
		}
		params := strings.Split(line, "|")
		if len(params) != 4 {
			fmt.Fprintln(os.Stdout, "***FAIL***")
			continue
		}
		concat, err := CreateConcat(params[0], params[1], params[2], params[3])
		if err != nil {
			fmt.Fprintln(os.Stdout, "***FAIL***")
		} else {
			fmt.Fprintln(os.Stdout, concat)
		}
	}
}

/*

A
0x00C4 0x00E4 0x00C0 0x00E0 0x00C1 0x00E1
0x00C2 0x00E2 0x00C3 0x00E3 0x00C5 0x00E5
0x01CD 0x01CE 0x0104 0x0105 0x0102 0x0103
0x00C6 0x00E6

C
0x00C7 0x00E7 0x0106 0x0107 0x0108 0x0109 0x010C
0x010D

D
0x010E 0x0111 0x0110 0x010F 0x00F0

E
0x00C8 0x00E8 0x00C9 0x00E9 0x00CA 0x00EA
0x00CB 0x00EB 0x011A 0x011B 0x0118 0x0119

G
0x011C 0x011D 0x0122 0x0123 0x011E 0x011F

H
0x0124 0x0125

I
0x00CC 0x00EC 0x00CD 0x00ED 0x00CE 0x00EE
0x00CF 0x00EF 0x0131

J
0x0134 0x0135

K
0x0136 0x0137

L
0x0139 0x013A 0x013B 0x013C 0x0141 0x0142 0x013D
0x013E

N
0x00D1 0x00F1 0x0143 0x0144 0x0147 0x0148

O
0x00D6 0x00F6 0x00D2 0x00F2 0x00D3 0x00F3
0x00D4 0x00F4 0x00D5 0x00F5 0x0150 0x0151
0x00D8 0x00F8 0x0152 0x0153

R
0x0154 0x0155 0x0158 0x0159

S
0x1E9E 0x00DF 0x015A 0x015B 0x015C 0x015D
0x015E 0x015F 0x0160 0x0161 0x0218 0x0219

T
0x0164 0x0165 0x0162 0x0163 0x00DE 0x00FE 0x021A
0x021B

U
0x00DC 0x00FC 0x00D9 0x00F9 0x00DA 0x00FA
0x00DB 0x00FB 0x0170 0x0171 0x0168 0x0169 0x0172
0x0173 0x016E 0x016F

W
0x0174 0x0175

Y
0x00DD 0x00FD 0x0178 0x00FF 0x0176 0x0177

Z
0x0179 0x017A 0x017D 0x017E 0x017B 0x017C
*/
