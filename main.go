package main

import (
	"os"

	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/andsala/geoip/ipdata"
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
	app.Version = "0.2.0"
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

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func getFlagRepr(data ipdata.Data) string {
	var flag = data.EmojiFlag
	if opt.NoColor || len(flag) < 1 {
		flag = data.Flag
	}
	return flag
}

func timeZoneString(timeZone ipdata.TimeZone) string {
	name := len(timeZone.Name) > 0
	abbr := len(timeZone.Abbr) > 0
	offset := len(timeZone.Offset) > 0
	//current := len(timeZone.CurrentTime) > 0
	out := ""

	if name {
		out += timeZone.Name
		if timeZone.IsDST {
			out += " DST"
		}
		if offset {
			out += ","
		}
		out += " "
	}
	if offset {
		out += "GMT" + timeZone.Offset
		if !name && bool(timeZone.IsDST) {
			out += " DST"
		}
	}
	if abbr {
		if len(out) > 0 {
			out += fmt.Sprintf(" (%s)", timeZone.Abbr)
		} else {
			out += timeZone.Abbr
		}
	}
	return out
}

func currencyString(currency ipdata.Currency) string {
	name := len(currency.Name) > 0
	code := len(currency.Code) > 0
	symbol := len(currency.Symbol) > 0
	out := ""

	if name {
		out += currency.Name
		if code && symbol {
			out += fmt.Sprintf(" (%s, %s)", currency.Code, currency.Symbol)
		} else if code {
			out += fmt.Sprintf(" (%s)", currency.Code)
		} else if symbol {
			out += fmt.Sprintf(" (%s)", currency.Symbol)
		}
	} else {
		if code && symbol {
			out += fmt.Sprintf("%s (%s)", currency.Code, currency.Symbol)
		} else if code {
			out += currency.Code
		} else if symbol {
			out += currency.Symbol
		}
	}
	return out
}

func languagesString(languages []ipdata.Language) string {
	out := ""

	for _, language := range languages {
		out += language.Name + ", "
	}
	out = strings.TrimSuffix(out, ", ")

	return out
}

func threatString(threat ipdata.Threat) string {
	out := ""
	var threats []string

	if threat.IsTor {
		threats = append(threats, "Tor")
	}
	if threat.IsProxy {
		threats = append(threats, "Proxy")
	}
	if threat.IsAnonymous {
		threats = append(threats, "Anonymous")
	}
	if threat.IsKnownAttacker {
		threats = append(threats, "Known attacker")
	}
	if threat.IsKnownAbuser {
		threats = append(threats, "Known abuser")
	}
	if threat.IsThreat {
		threats = append(threats, "Threat")
	}
	if threat.IsBogon {
		threats = append(threats, "Bogon")
	}

	if len(threats) > 0 {
		for _, threat := range threats {
			out += threat + ", "
		}
		out = strings.TrimSuffix(out, ", ")
	} else {
		out = "None"
	}

	return out
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

		timezone := timeZoneString(data.TimeZone)
		if len(timezone) > 0 {
			out += "   Time zone:       " + timezone + "\n"
		}

		currency := currencyString(data.Currency)
		if len(currency) > 0 {
			out += "   Currency:        " + currency + "\n"
		}

		languages := languagesString(data.Languages)
		if len(languages) > 0 {
			out += "   Languages:       " + languages + "\n"
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

		out += "   Threat:          " + threatString(data.Threat) + "\n"

		if !strings.HasSuffix(out, "\n\n") {
			out += "\n"
		}
	}

	fmt.Print(out)
}
