package m3u

type Channels 		[]Channel
type ChannelGroups 	[]ChannelGroup

type Channel struct {
	Name 		string
	StreamLink 	string
}

type ChannelGroup struct {
	Name			string
	ChannelList	 	Channels
}