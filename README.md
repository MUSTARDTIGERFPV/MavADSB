# MavADSB

> [!IMPORTANT]
**This functionality is now built into Mission Planner via [this PR](https://github.com/ArduPilot/MissionPlanner/pull/3251). You likely don't need this tool anymore. :)**

## Description

Loads ADSB data into Mission Planner or QGC without a need for an ADSB receiver using public HTTPS APIs.
This data can be forwarded to the UAS (via telemetry link) for collision avoidance without the need for an onboard ADSB-in receiver or any additional hardware.
Works natively on Mac, Linux or Windows without runtimes or additional dependencies.

This is a port of https://github.com/MohammadAdib/MavADSB to Golang with some additional fixes

### Basic usage
Download a release from https://github.com/MUSTARDTIGERFPV/MavADSB/releases/latest for your platform. Builds are provided for Windows, macOS, Linux, FreeBSD, OpenBSD, and select ARM variants.

The application will automatically detect your location by your IP address, and pull 250nm radius from there. If you wish to override that, see the configuration section.

### Configuration
Configuration is handled as runtime flags. As an example, to set the SBS serving port to 1234, you'd add `-sbs.port 1234` to the end of your command line execution.
```
  -http.port string
        HTTP Serving port (default "3000")
  -location.lat float
        User latitude
  -location.lng float
        User longitude
  -location.radius uint
        Radius to request data for in nautical miles (default 250)
  -sbs.port string
        SBS Serving port (default "30003")
  -upstream.api_base string
        ADSB.one API base URL (default "https://api.adsb.one/v2")
  -upstream.refresh_interval uint
        Interval in seconds between API calls (default 5)
```

### Mission Planner integration
After running the application, please open Mission Planner and go to config tab -> Planner and look for the "Adsb" checkbox. Enter the details for your MavADSB server (the defaults will work if you're running it on the same machine with default settings)
Enable this and restart Mission Planner

### QGroundControl integration
After running the application, please open QGC and click the Q on the top-left. Go into application settings and scroll to the bottom. Enable ADSB and restart QGC

### API & Docs
ADSB-One API: https://github.com/ADSB-One/api/blob/main/README.md
SBS-1 Info: http://woodair.net/SBS/Article/Barebones42_Socket_Data.htm

The SBS-1 Server listens on 0.0.0.0:30003
The application serves metrics on HTTP port 3000 in Prometheus format

### Traditional methods
GCS connected USB dongle: https://uavionix.com/products/pingusb/
UAS connected ADSB rx: https://uavionix.com/products/pingrx-pro/
