# geoip - IP geolocation with ipdata.co

[![Build Status](https://travis-ci.org/andsala/geoip.svg?branch=master)](https://travis-ci.org/andsala/geoip)
[![Go Report Card](https://goreportcard.com/badge/github.com/andsala/geoip)](https://goreportcard.com/report/github.com/andsala/geoip)
[![codecov](https://codecov.io/gh/andsala/geoip/branch/master/graph/badge.svg)](https://codecov.io/gh/andsala/geoip)
[![GoDoc](https://godoc.org/github.com/andsala/geoip/ipdata?status.svg)](https://godoc.org/github.com/andsala/geoip/ipdata)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fandsala%2Fgeoip.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fandsala%2Fgeoip?ref=badge_shield)

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
--no-color              Disable color and emoji output [NO_COLOR] (http://no-color.org)
--user-agent, -u value  HTTP user agent [GEOIP_USER_AGENT]
```

# License
This project is distributed under the [MIT License](https://github.com/andsala/geoip/blob/master/LICENSE).


[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fandsala%2Fgeoip.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fandsala%2Fgeoip?ref=badge_large)