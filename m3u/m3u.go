package m3u

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

/*
	Auxiliary Functions
*/
func readFile(file string) string {
	var lines []string
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error on readFile: %v", err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()

	return strings.Join(lines[:], "\n")
}

func parseM3U(fileContent string) ChannelGroups {
	var parsedData ChannelGroups
	re := regexp.MustCompile(`(?m)group-title="(.*?)",(.*?)\n(.*?)\n`)

	matches := re.FindAllStringSubmatch(fileContent, -1)
	for _, match := range matches {
		channel := Channel{}

		groupName := string(match[1])
		channel.Name = string(match[2])
		channel.StreamLink = string(match[3])

		indexOfGroup, err := parsedData.SearchGroupByName(groupName)
		if err != nil {
			group := ChannelGroup{}
			group.Name = groupName
			group.ChannelList = append(group.ChannelList, channel)
			parsedData = append(parsedData, group)
		} else {
			group := &parsedData[indexOfGroup]
			group.ChannelList = append(group.ChannelList, channel)
		}

	}
	return parsedData
}

/*
	Parsers Functions
*/

// ParseFile parses an m3u file
func ParseFile(fileName string) ChannelGroups {
	fileContent := readFile(fileName)
	return parseM3U(fileContent)
}

// ParseURL parses an url of m3u file
func ParseURL(url string) ChannelGroups {
	var result ChannelGroups
	response, err := http.Get(url)

	if err != nil {
		return result
	}
	defer response.Body.Close()

	fileContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result
	}

	result = parseM3U(string(fileContent))

	return result
}

/*
	Search Function
*/
// SearchGroupByName searchs an group of channels by name
func (groups ChannelGroups) SearchGroupByName(groupName string) (int, error) {
	for i, group := range groups {
		if group.Name == groupName {
			return i, nil
		}
	}

	return -1, errors.New("Group not found")
}

// SearchChannelsByName searchs an channel in an ChannelList by name
func (groups ChannelGroups) SearchChannelsByName(channelName string) Channels {
	var channels Channels
	commChannel := make(chan Channels, 100)
	foundChannel := make(chan Channel, 100)

	go func() {
		defer close(commChannel)
		for _, group := range groups {
			commChannel <- group.ChannelList
		}
	}()

	go func() {
		defer close(foundChannel)
		for channelList := range commChannel {
			for _, channel := range channelList {
				if strings.Contains(strings.ToLower(channel.Name), strings.ToLower(channelName)) {
					foundChannel <- channel
				}
			}
		}
	}()

	for channel := range foundChannel {
		channels = append(channels, channel)
	}

	return channels
}

// SearchChannelsByName search an channel in ChannelList by link
func (groups ChannelGroups) SearchChannelsByLink(streamLink string) Channel {
	channelResult := Channel{}
	commChannel := make(chan Channels, 100)

	go func() {
		defer close(commChannel)
		for _, group := range groups {
			commChannel <- group.ChannelList
		}
	}()

	for channelList := range commChannel {
		for _, channel := range channelList {
			if channel.StreamLink == streamLink {
				channelResult = channel
				break
			}
		}
	}

	return channelResult

}
