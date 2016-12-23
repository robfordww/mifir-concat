// mifirconcat_test.go
package main

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestErrors(t *testing.T) {
	// Illegal country code, dates, empty string,
	concatTester(t, "NX", "19801224", "Mr Jon", "Anderson", "")
	concatTester(t, "NO", "19801224", "", "", "")
	concatTester(t, "NO", "198012241", "Jon", "Snow", "")
	concatTester(t, "NO", "1980122", "Jon", "Snow", "")
	concatTester(t, "UK", "19800122", "Sir Jon", "Snow", "")
	concatTester(t, "NOK", "19800122", "Sir Jon", "Snow", "")
}

func TestExpected(t *testing.T) {
	concatTester(t, "NO", "19801224", "Jon", "Snow", "NO19801224JON##SNOW#")
	concatTester(t, "GB", "19800122", "Sir Jon", "Snow", "GB19800122JON##SNOW#")
	concatTester(t, "US", "19800502", "Dr. Robert", "Ford", "US19800502ROBERFORD#")
	concatTester(t, "NO", "19801224", "Jon", "Snow", "NO19801224JON##SNOW#")
	concatTester(t, "NO", "19801224", "Jon", "Snow", "NO19801224JON##SNOW#")
}

// Verify ESMA Guideline examples
func TestGuidelineExamples(t *testing.T) {
	/* John O'Brian  */
	concatTester(t, "IE", "19800113", "John", "O'Brian", "IE19800113JOHN#OBRIA")
	// Ludwig Van der Rohe
	concatTester(t, "HU", "19810214", "Ludwig", "Van der Rohe", "HU19810214LUDWIROHE#")
	// Victor Vandenberg US19730322VICTOVANDE
	concatTester(t, "US", "19730322", "Victor", "Vandenberg", "US19730322VICTOVANDE")
	// Eli Ødegård
	concatTester(t, "NO", "19760315", "Eli", "Ødegård", "NO19760315ELI##ODEGA")
	// Willeke de Bruijn
	concatTester(t, "LU", "19660416", "Willeke", "de Bruijn", "LU19660416WILLEBRUIJ")
	// Jon Ian Dewitt
	concatTester(t, "US", "19650417", "Jon Ian", "Dewitt", "US19650417JON##DEWIT")
	// Amy-Ally Garção de Magalhães
	concatTester(t, "PT", "19900517", "Amy-Ally", "Garção de Magalhães", "PT19900517AMYALGARCA")
	// Giovani dos Santos
	concatTester(t, "FR", "19900618", "Giovani", "dos Santos", "FR19900618GIOVASANTO")
	// Günter Voẞ
	concatTester(t, "DE", "19800715", "Günter", "Voẞ", "DE19800715GUNTEVOS##")
}

// Verify guidlien inline examples. (Derived from other examples in the text)
func TestGuidelineExamplesExtended(t *testing.T) {
	// Sean Murphy
	concatTester(t, "IE", "19760227", "SEAN", "MURPHY", "IE19760227SEAN#MURPH")
	// Thomas Maccormack
	concatTester(t, "IE", "19511212", "THOMAS", "MACCORMACK", "IE19511212THOMAMACCO")
	// Pierre Marie Dupont
	concatTester(t, "FR", "19760227", "PIERRE MARIE", "DUPONT", "FR19760227PIERRDUPON")
}

func concatTester(t *testing.T, country, date, firstname, lastname string, expected string) {
	c, err := createConcat(country, date, firstname, lastname)

	if expected == "" {
		// Expect an error
		if err == nil || c != "" {
			t.Error(strings.Join([]string{country, date, firstname, lastname}, " "))
			t.Fail()
		}
		return
	}

	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if c != expected {
		t.Error("Expected", expected, " got ", c)
		t.Fail()
	}
}

func init() {

}

func BenchmarkCONCAT(b *testing.B) {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ."
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	randletters := func(n int) string {
		b := make([]byte, n)
		// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
		for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
			if remain == 0 {
				cache, remain = rand.Int63(), letterIdxMax
			}
			if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
				b[i] = letterBytes[idx]
				i--
			}
			cache >>= letterIdxBits
			remain--
		}
		return string(b)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			datestr := time.Now().Add(-time.Hour * 24 * time.Duration(rand.Int()%1000)).Format("20060102")
			firstname := randletters(8)
			lastname := randletters(8)
			createConcat(ccalpha2[int(rand.Int63())%len(ccalpha2)], datestr,
				firstname, lastname)
		}
	})
}
