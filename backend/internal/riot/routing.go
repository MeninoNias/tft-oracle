package riot

import "strings"

// RegionToBaseURL maps a routing region to its base API URL.
// Regional endpoints handle Account V1 and TFT Match V1.
func RegionToBaseURL(region string) string {
	switch strings.ToLower(region) {
	case "americas":
		return "https://americas.api.riotgames.com"
	case "europe":
		return "https://europe.api.riotgames.com"
	case "asia":
		return "https://asia.api.riotgames.com"
	case "sea":
		return "https://sea.api.riotgames.com"
	default:
		return "https://americas.api.riotgames.com"
	}
}

// PlatformToBaseURL maps a platform to its base API URL.
// Platform endpoints handle TFT Summoner V1 and TFT League V1.
func PlatformToBaseURL(platform string) string {
	return "https://" + strings.ToLower(platform) + ".api.riotgames.com"
}

// ServerInfo holds routing data for a player's server.
type ServerInfo struct {
	Region   string // Regional routing: "americas", "europe", "asia", "sea"
	Platform string // Platform routing: "br1", "na1", "euw1", etc.
}

// serverMap maps user-facing server names to routing info.
var serverMap = map[string]ServerInfo{
	// Americas
	"br":   {Region: "americas", Platform: "br1"},
	"br1":  {Region: "americas", Platform: "br1"},
	"na":   {Region: "americas", Platform: "na1"},
	"na1":  {Region: "americas", Platform: "na1"},
	"lan":  {Region: "americas", Platform: "la1"},
	"la1":  {Region: "americas", Platform: "la1"},
	"las":  {Region: "americas", Platform: "la2"},
	"la2":  {Region: "americas", Platform: "la2"},
	"oce":  {Region: "americas", Platform: "oc1"},
	"oc1":  {Region: "americas", Platform: "oc1"},
	// Europe
	"euw":  {Region: "europe", Platform: "euw1"},
	"euw1": {Region: "europe", Platform: "euw1"},
	"eune": {Region: "europe", Platform: "eun1"},
	"eun1": {Region: "europe", Platform: "eun1"},
	"tr":   {Region: "europe", Platform: "tr1"},
	"tr1":  {Region: "europe", Platform: "tr1"},
	"ru":   {Region: "europe", Platform: "ru"},
	// Asia
	"kr":  {Region: "asia", Platform: "kr"},
	"jp":  {Region: "asia", Platform: "jp1"},
	"jp1": {Region: "asia", Platform: "jp1"},
	// SEA
	"ph":  {Region: "sea", Platform: "ph2"},
	"ph2": {Region: "sea", Platform: "ph2"},
	"sg":  {Region: "sea", Platform: "sg2"},
	"sg2": {Region: "sea", Platform: "sg2"},
	"th":  {Region: "sea", Platform: "th2"},
	"th2": {Region: "sea", Platform: "th2"},
	"tw":  {Region: "sea", Platform: "tw2"},
	"tw2": {Region: "sea", Platform: "tw2"},
	"vn":  {Region: "sea", Platform: "vn2"},
	"vn2": {Region: "sea", Platform: "vn2"},
}

// ResolveServer looks up routing info from a server/region string.
// Accepts platform codes ("br1"), short names ("br"), or region names ("americas").
// For broad region names without a platform, defaults to the most common platform.
func ResolveServer(input string) ServerInfo {
	lower := strings.ToLower(strings.TrimSpace(input))

	if info, ok := serverMap[lower]; ok {
		return info
	}

	// Fallback: treat as broad region name
	switch lower {
	case "americas":
		return ServerInfo{Region: "americas", Platform: "na1"}
	case "europe":
		return ServerInfo{Region: "europe", Platform: "euw1"}
	case "asia":
		return ServerInfo{Region: "asia", Platform: "kr"}
	case "sea":
		return ServerInfo{Region: "sea", Platform: "sg2"}
	default:
		return ServerInfo{Region: "americas", Platform: "na1"}
	}
}
