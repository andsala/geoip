# geoip - IP geolocation with ipdata.co

[![Build Status](https://travis-ci.org/andsala/geoip.svg?branch=master)](https://travis-ci.org/andsala/geoip)
[![Go Report Card](https://goreportcard.com/badge/github.com/andsala/geoip)](https://goreportcard.com/report/github.com/andsala/geoip)
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

   Flag:            https://ipdata.co/flags/us.png
   Currency:        USD ($)
   Calling code:    +1

   Organization:    Google LLC
   AS number:       AS15169

```

## Options
```
--api-key, -a value     ipdata.co api key [GEOIP_API_KEY]
--ip-only, --ip         Print current public IP and exit
--json, -j              Print pure json
--user-agent, -u value  HTTP user agent [GEOIP_USER_AGENT]
```

# License
This project is distributed under [MIT license](https://opensource.org/licenses/MIT).
