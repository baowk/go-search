package ua

import "fmt"

// GenerateFirefoxUA generates a Firefox user agent string based on the provided parameters
func GenerateFirefoxUA(platform string, geckoVersion string, firefoxVersion string, isMobile bool, isTablet bool) string {
	// For mobile Firefox 10+, gecko-trail is the same as firefox-version
	geckoTrail := "20100101"
	if isMobile && firefoxVersion >= "10.0" {
		geckoTrail = firefoxVersion
	}

	// Handle mobile and tablet indicators
	if isMobile {
		if platform != "" {
			platform += "; Mobile"
		} else {
			platform = "Mobile"
		}
	} else if isTablet {
		if platform != "" {
			platform += "; Tablet"
		} else {
			platform = "Tablet"
		}
	}

	// Format the user agent string
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%s) Gecko/%s Firefox/%s",
		platform, geckoVersion, geckoTrail, firefoxVersion)
}

// GenerateWindowsFirefoxUA generates a Firefox user agent for Windows platforms
func GenerateWindowsFirefoxUA(ntVersion string, is64bit bool, geckoVersion string, firefoxVersion string) string {
	if ntVersion == "" {
		ntVersion = "10.0"
	}
	platform := fmt.Sprintf("Windows NT %s", ntVersion)
	if is64bit {
		platform += "; Win64; x64"
	}
	return GenerateFirefoxUA(platform, geckoVersion, firefoxVersion, false, false)
}

// GenerateMacFirefoxUA generates a Firefox user agent for macOS platforms
func GenerateMacFirefoxUA(macVersion string, isPPC bool, geckoVersion string, firefoxVersion string) string {
	architecture := "Intel"
	if isPPC {
		architecture = "PPC"
	}
	platform := fmt.Sprintf("Macintosh; %s Mac OS X %s", architecture, macVersion)
	return GenerateFirefoxUA(platform, geckoVersion, firefoxVersion, false, false)
}

// GenerateLinuxFirefoxUA generates a Firefox user agent for Linux platforms
func GenerateLinuxFirefoxUA(cpuArch string, geckoVersion string, firefoxVersion string) string {
	platform := fmt.Sprintf("X11; Linux %s", cpuArch)
	return GenerateFirefoxUA(platform, geckoVersion, firefoxVersion, false, false)
}

// GenerateAndroidFirefoxUA generates a Firefox user agent for Android platforms
func GenerateAndroidFirefoxUA(androidVersion string, isMobile bool, isTablet bool, geckoVersion string, firefoxVersion string) string {
	formFactor := ""
	if isMobile {
		formFactor = "Mobile"
	} else if isTablet {
		formFactor = "Tablet"
	}

	platform := fmt.Sprintf("Android %s; %s", androidVersion, formFactor)
	return GenerateFirefoxUA(platform, geckoVersion, firefoxVersion, isMobile, isTablet)
}

// GenerateFirefoxOSUA generates a Firefox user agent for Firefox OS
func GenerateFirefoxOSUA(isMobile bool, isTablet bool, isTV bool, deviceID string, geckoVersion string, firefoxVersion string) string {
	var platform string
	switch {
	case isTV:
		platform = "TV"
	case isTablet:
		platform = "Tablet"
	case deviceID != "":
		platform = fmt.Sprintf("Mobile; %s", deviceID)
	default:
		platform = "Mobile"
	}

	// For Firefox OS, we don't include the platform details in the same way
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%s) Gecko/%s Firefox/%s",
		platform, geckoVersion, geckoVersion, firefoxVersion)
}
