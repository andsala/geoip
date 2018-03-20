package main

import (
	"os"

	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/andsala/geoip/ipdata"
	"github.com/kyokomi/emoji"
	"gopkg.in/urfave/cli.v2"
)

type options struct {
	APIKey    string
	UserAgent string
	IPOnly    bool
	JSON      bool
	NoColor   bool
}

var opt = options{}

func main() {
	app := &cli.App{}

	app.Name = "geoip"
	app.Usage = "Get info about IP geolocation from ipdata.co"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.ArgsUsage = "[IP...]"
	app.CustomAppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}
`

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "api-key",
			Aliases:     []string{"a"},
			Value:       "",
			Usage:       "ipdata.co api key",
			EnvVars:     []string{"GEOIP_API_KEY"},
			Destination: &opt.APIKey,
		},
		&cli.StringFlag{
			Name:        "user-agent",
			Aliases:     []string{"u"},
			Value:       "andsala_" + app.Name + "/" + app.Version,
			Usage:       "HTTP user agent",
			EnvVars:     []string{"GEOIP_USER_AGENT"},
			Destination: &opt.UserAgent,
		},
		&cli.BoolFlag{
			Name:        "ip-only",
			Aliases:     []string{"ip"},
			Usage:       "Print current public IP and exit",
			Value:       false,
			Destination: &opt.IPOnly,
		},
		&cli.BoolFlag{
			Name:        "json",
			Aliases:     []string{"j"},
			Usage:       "Print pure json",
			Value:       false,
			Destination: &opt.JSON,
		},
		&cli.BoolFlag{
			Name:        "no-color",
			Usage:       "Disable color and emoji output",
			Value:       false,
			EnvVars:     []string{"NO_COLOR"},
			Destination: &opt.NoColor,
		},
	}

	app.Action = func(ctx *cli.Context) error {
		client, err := ipdata.NewClient(nil)
		if err != nil {
			return cli.Exit(err, 2)
		}

		client.UserAgent = opt.UserAgent
		if len(opt.APIKey) > 0 {
			client.APIKey = opt.APIKey
		}

		if opt.IPOnly {
			data, err := client.GetMyIPData()
			if err != nil {
				return cli.Exit(err, 2)
			}
			fmt.Println(data.IP)
			return nil
		}

		if ctx.NArg() == 0 {
			data, err := client.GetMyIPData()
			if err != nil {
				return cli.Exit(err, 2)
			}
			printIPData(*data)
		} else {
			for _, ip := range ctx.Args().Slice() {
				data, err := client.GetIPData(ip)
				if err != nil {
					return cli.Exit(err, 2)
				}
				printIPData(*data)
			}
		}

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	app.Run(os.Args)
}

func getFlagRepr(data ipdata.Data) string {
	var flagKey = fmt.Sprintf(":%v:", strings.ToLower(data.CountryCode))
	var flag = emoji.Sprintf(flagKey)
	if opt.NoColor || flag == flagKey {
		flag = data.Flag
	}
	return flag
}

func printIPData(data ipdata.Data) {
	var out = ""

	if opt.JSON {
		out = *data.JSON
	} else {
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
			out += "   Flag:            " + getFlagRepr(data) + "\n"
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
	}

	fmt.Print(out)
}
