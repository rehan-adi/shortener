package utils

import "github.com/mssola/user_agent"

func ParseUserAgent(uaString string) (device, browser, os string) {

	ua := user_agent.New(uaString)
	name, version := ua.Browser()

	if ua.Mobile() {
		device = "Mobile"
	} else if ua.Bot() {
		device = "Bot"
	} else {
		device = "Desktop"
	}

	browser = name + " " + version
	os = ua.OS()
	return
}
