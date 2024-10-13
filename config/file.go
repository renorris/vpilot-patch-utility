package config

import (
	"fmt"
	"regexp"
)

// UpdateXML updates the network status and cached server fields.
// It also clears the NetworkLogin and NetworkPassword fields.
func UpdateXML(xmlData []byte, networkStatusURL string, cachedServers []string) ([]byte, error) {
	// Prepare regex patterns to match the relevant XML elements
	networkStatusURLPattern := regexp.MustCompile(`<NetworkStatusURL>[\s\S]*?</NetworkStatusURL>`)
	cachedServersPattern := regexp.MustCompile(`<CachedServers>[\s\S]*?</CachedServers>`)

	// Replace NetworkStatusURL with the new value
	newNetworkStatusURL := fmt.Sprintf("<NetworkStatusURL>%s</NetworkStatusURL>", networkStatusURL)
	xmlData = networkStatusURLPattern.ReplaceAll(xmlData, []byte(newNetworkStatusURL))

	// Construct the new CachedServers XML
	cachedServersXML := "<CachedServers>\n"
	for _, server := range cachedServers {
		cachedServersXML += fmt.Sprintf("    <string>%s</string>\n", server)
	}
	cachedServersXML += "  </CachedServers>"

	// Replace CachedServers with the new value
	xmlData = cachedServersPattern.ReplaceAll(xmlData, []byte(cachedServersXML))

	// Clear CID and password fields
	networkLoginPattern := regexp.MustCompile(`<NetworkLogin>[\s\S]*?</NetworkLogin>`)
	networkPasswordPattern := regexp.MustCompile(`<NetworkPassword>[\s\S]*?</NetworkPassword>`)

	xmlData = networkLoginPattern.ReplaceAll(xmlData, []byte("<NetworkLogin></NetworkLogin>"))
	xmlData = networkPasswordPattern.ReplaceAll(xmlData, []byte("<NetworkPassword></NetworkPassword>"))

	return xmlData, nil
}
