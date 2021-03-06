package m3u

import (
	"testing"
)

var cg ChannelGroups

func TestSearchGroupByName(t *testing.T) {
	var cg ChannelGroups
	testString := []string{"a", "b", "c", "d", "e", "f", "g"}

	for _, element := range testString {
		ch := ChannelGroup{}
		ch.Name = element
		cg = append(cg, ch)
	}

	index, err := cg.SearchGroupByName("c")

	if err != nil {
		t.Errorf("Not found an existent group")
	}

	if index != 2 {
		t.Errorf("Expected index %d and founded index %d", 2, index)
	}

	index, err = cg.SearchGroupByName("A")

	if err == nil {
		t.Errorf("Found an non-existent group")
	}

}

func TestParseFile(t *testing.T) {
	m3uFile := "iptv.m3u"
	cg = ParseFile(m3uFile)

	if cg == nil || len(cg) <= 0 {
		t.Errorf("ParsedFile is nil or results <= 0")
	}
}

func TestSearchChannelsByName(t *testing.T) {
	foundChannels := cg.SearchChannelsByName("Globo RJ")
	expectedResult := 5
	foundResult := len(foundChannels)

	if foundResult != expectedResult {
		t.Errorf("Expected %d channels, found %d", expectedResult, foundResult)
	}
}

func TestSearchChannelsByLink(t *testing.T) {
	foundResult := cg.SearchChannelsByLink("http://pfsv.io:80/8019390/3415/392")

	if foundResult.StreamLink != "http://pfsv.io:80/8019390/3415/392" {
		t.Errorf("Expected http://pfsv.io:80/8019390/3415/392 channel, found %v+\n", foundResult)
	}
}

func TestParseURL(t *testing.T) {
	url := "http://ccsv.me/3206464/8372"
	cg := ParseURL(url)

	if cg == nil {
		t.Errorf("Failed to get page, expected != nil")
	}

	if len(cg) <= 0 {
		t.Errorf("The response of url is 0 length")
	}
}

func TestReadFile(t *testing.T) {
	result := readFile("iptv.m3u")
	if result == "" {
		t.Errorf("FileRead isn't acting like expected")
	}
}
