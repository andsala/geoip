# geoip - IP geolocation with ipdata.co

[![Build Status](https://travis-ci.org/andsala/geoip.svg?branch=master)](https://travis-ci.org/andsala/geoip)
[![Go Report Card](https://goreportcard.com/badge/github.com/andsala/geoip)](https://goreportcard.com/report/github.com/andsala/geoip)
[![codecov](https://codecov.io/gh/andsala/geoip/branch/master/graph/badge.svg)](https://codecov.io/gh/andsala/geoip)
[![GoDoc](https://godoc.org/github.com/andsala/geoip/ipdata?status.svg)](https://godoc.org/github.com/andsala/geoip/ipdata)

# Installation
```sh
go get -u github.com/andsala/geoip
```

# Usage
```sh
$ geoip 8.8.8.8
IP: 8.8.8.8
   United States (US)
   North America (NA)
   Coordinates:     37.751, -97.822

   Flag:            🇺🇸
   Time zone:       America/Chicago, GMT-0500 (CDT)
   Currency:        US Dollar (USD, $)
   Languages:       English
   Calling code:    +1

   Organization:    Google LLC (google.com)
   AS number:       AS15169
   Threat:          None

```

## Options
```
--api-key value, -a value     ipdata.co api key [$GEOIP_API_KEY]
--ip-only, --ip               Print current public IP and exit
--json, -j                    Print pure json
--no-color                    Disable color and emoji output [$NO_COLOR] (http://no-color.org)
--user-agent value, -u value  HTTP user agent [$GEOIP_USER_AGENT]
--help, -h                    show help
--version, -v                 print the version
```

# License
This project is distributed under the [MIT License](https://github.com/andsala/geoip/blob/master/LICENSE).
