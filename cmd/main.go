package main

import (
	"os"

	"fmt"
	"github.com/andsala/geoip/ipdata"
	"github.com/urfave/cli"
	"strings"
)

type Options struct {
	ApiKey    string
	UserAgent string
}

var opt = Options{}

func main() {
	app := cli.NewApp()

	app.Name = "geoip"
	app.Usage = "Get info about IP geolocation"
	app.Version = "0.0.1"
	app.Author = "andsala"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "api-key, a",
			Value:       "",
			Usage:       "ipdata.co api key",
			EnvVar:      "GEOIP_API_KEY",
			Destination: &opt.ApiKey,
		},
		cli.StringFlag{
			Name:        "user-agent, u",
			Value:       "",
			Usage:       "HTTP user agent",
			EnvVar:      "GEOIP_USER_AGENT",
			Destination: &opt.UserAgent,
		},
	}

	app.Action = func(ctx *cli.Context) error {
		client, err := ipdata.NewClient(nil)
		if err != nil {
			return cli.NewExitError(err, 2)
		}

		client.UserAgent = opt.UserAgent
		if len(opt.ApiKey) > 0 {
			client.ApiKey = opt.ApiKey
		}

		if ctx.NArg() == 0 {
			data, err := client.GetMyIpData()
			if err != nil {
				return cli.NewExitError(err, 2)
			}
			printIPData(*data)
		} else {
			// Other IP data
		}

		return nil
	}

	app.Run(os.Args)
}

func printIPData(data ipdata.Data) {
	var out string = ""

	out += "IP: " + data.IP + "\n"

	if len(data.Postal) > 0 {
		if len(data.Region) > 0 {
			out += fmt.Sprintf("   %v %v\n", data.Postal, data.City)
		} else {
			out += fmt.Sprintf("   %v\n", data.Postal)
		}
	} else {
		if len(data.Region) > 0 {
			out += fmt.Sprintf("   %v\n", data.City)
		}
	}

	if len(data.Region) > 0 {
		out += fmt.Sprintf("   %v\n", data.Region)
	}

	if len(data.CountryName) > 0 {
		out += "   " + data.CountryName
		if len(data.CountryCode) > 0 {
			out += fmt.Sprintf(" (%v)", data.CountryCode)
		}
		out += "\n"
	}

	if len(data.ContinentName) > 0 {
		out += "   " + data.ContinentName
		if len(data.ContinentCode) > 0 {
			out += fmt.Sprintf(" (%v)", data.ContinentCode)
		}
		out += "\n"
	}

	out += fmt.Sprintf("   Coordinates:     %g, %g\n", data.Latitude, data.Longitude)
	out += "\n"

	if len(data.Flag) > 0 {
		out += "   Flag:            " + data.Flag + "\n"
	}

	if len(data.TimeZone) > 0 {
		out += "   Timezone:        " + data.TimeZone + "\n"
	}

	if len(data.Currency) > 0 {
		out += "   Currency:        " + data.Currency
		if len(data.CurrencySymbol) > 0 {
			out += " (" + data.CurrencySymbol + ")"
		}
		out += "\n"
	}

	if len(data.CallingCode) > 0 {
		out += "   Calling code:    +" + data.CallingCode + "\n"
	}
	out += "\n"

	if len(data.Organisation) > 0 {
		out += "   Organization:    " + data.Organisation + "\n"
	}

	if len(data.ASN) > 0 {
		out += "   AS number:       " + data.ASN + "\n"
	}

	if !strings.HasSuffix(out, "\n\n") {
		out += "\n"
	}

	fmt.Print(out)
}
