package riot

import "testing"

func TestResolveServer(t *testing.T) {
	tests := []struct {
		input            string
		expectedRegion   string
		expectedPlatform string
	}{
		{"br", "americas", "br1"},
		{"br1", "americas", "br1"},
		{"na", "americas", "na1"},
		{"na1", "americas", "na1"},
		{"lan", "americas", "la1"},
		{"las", "americas", "la2"},
		{"oce", "americas", "oc1"},
		{"euw", "europe", "euw1"},
		{"euw1", "europe", "euw1"},
		{"eune", "europe", "eun1"},
		{"tr", "europe", "tr1"},
		{"ru", "europe", "ru"},
		{"kr", "asia", "kr"},
		{"jp", "asia", "jp1"},
		{"sg", "sea", "sg2"},
		{"ph", "sea", "ph2"},
		{"vn", "sea", "vn2"},
		// Broad region names
		{"americas", "americas", "na1"},
		{"europe", "europe", "euw1"},
		{"asia", "asia", "kr"},
		{"sea", "sea", "sg2"},
		// Case insensitivity
		{"BR", "americas", "br1"},
		{"EUW", "europe", "euw1"},
		// Whitespace trimming
		{"  br  ", "americas", "br1"},
		// Unknown fallback
		{"unknown", "americas", "na1"},
		{"", "americas", "na1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			info := ResolveServer(tt.input)
			if info.Region != tt.expectedRegion {
				t.Errorf("ResolveServer(%q).Region = %q, want %q", tt.input, info.Region, tt.expectedRegion)
			}
			if info.Platform != tt.expectedPlatform {
				t.Errorf("ResolveServer(%q).Platform = %q, want %q", tt.input, info.Platform, tt.expectedPlatform)
			}
		})
	}
}

func TestRegionToBaseURL(t *testing.T) {
	tests := []struct {
		region   string
		expected string
	}{
		{"americas", "https://americas.api.riotgames.com"},
		{"europe", "https://europe.api.riotgames.com"},
		{"asia", "https://asia.api.riotgames.com"},
		{"sea", "https://sea.api.riotgames.com"},
		{"AMERICAS", "https://americas.api.riotgames.com"},
		{"unknown", "https://americas.api.riotgames.com"},
	}

	for _, tt := range tests {
		t.Run(tt.region, func(t *testing.T) {
			got := RegionToBaseURL(tt.region)
			if got != tt.expected {
				t.Errorf("RegionToBaseURL(%q) = %q, want %q", tt.region, got, tt.expected)
			}
		})
	}
}

func TestPlatformToBaseURL(t *testing.T) {
	tests := []struct {
		platform string
		expected string
	}{
		{"br1", "https://br1.api.riotgames.com"},
		{"na1", "https://na1.api.riotgames.com"},
		{"KR", "https://kr.api.riotgames.com"},
	}

	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			got := PlatformToBaseURL(tt.platform)
			if got != tt.expected {
				t.Errorf("PlatformToBaseURL(%q) = %q, want %q", tt.platform, got, tt.expected)
			}
		})
	}
}
