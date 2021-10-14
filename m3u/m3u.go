package m3u

import (
	"regexp"
	"strings"
	"os"
	"errors"
	"log"
	"bufio"
	"net/http"
	"io/ioutil"
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

	return strings.Join(lines[:],"\n")
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

func searchChannelsByName(channelName string, groups ChannelGroups) (Channels) {
	var channels Channels
	commChannel := make(chan Channels, 100)
	
	go func() {
		defer close(commChannel)
		for _, group := range groups {
			commChannel <- group.ChannelList
		}
	}()

	for channelList := range commChannel {
		for _, channel := range channelList {
			if strings.Contains(strings.ToLower(channel.Name), strings.ToLower(channelName)) {
				channels = append(channels, channel)
			}

		}
	}

	return channels
}

/*
	Parsers Functions
*/
func ParseFile(fileName string) ChannelGroups {
	fileContent := readFile(fileName)
	return parseM3U(fileContent)
}

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

func (groups ChannelGroups) SearchChannelsByName(channelName string) (Channels) {
	var channels Channels
	halfSize := len(groups)/2
	
	chans1 := searchChannelsByName(channelName, groups[:halfSize])
	chans2 := searchChannelsByName(channelName, groups[halfSize:])

	channels = append(channels, chans1...)
	channels = append(channels, chans2...)

	return channels
}

func (groups ChannelGroups) SearchGroupByName(groupName string) (int, error) {
	for i, group := range groups {
		if group.Name == groupName {
			return i, nil
		}
	}

	return -1, errors.New("Group not found")
}

