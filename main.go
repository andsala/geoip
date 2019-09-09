package main

import (
	"fmt"
	"os"
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
	var out strings.Builder

	if name {
		out.WriteString(timeZone.Name)
		if timeZone.IsDST {
			out.WriteString(" DST")
		}
		if offset {
			out.WriteString(",")
		}
		out.WriteString(" ")
	}
	if offset {
		out.WriteString("GMT")
		out.WriteString(timeZone.Offset)
		if !name && bool(timeZone.IsDST) {
			out.WriteString(" DST")
		}
	}
	if abbr {
		if out.Len() > 0 {
			out.WriteString(fmt.Sprintf(" (%s)", timeZone.Abbr))
		} else {
			out.WriteString(timeZone.Abbr)
		}
	}
	return out.String()
}

func currencyString(currency ipdata.Currency) string {
	name := len(currency.Name) > 0
	code := len(currency.Code) > 0
	symbol := len(currency.Symbol) > 0
	var out strings.Builder

	if name {
		out.WriteString(currency.Name)
		if code && symbol {
			out.WriteString(fmt.Sprintf(" (%s, %s)", currency.Code, currency.Symbol))
		} else if code {
			out.WriteString(fmt.Sprintf(" (%s)", currency.Code))
		} else if symbol {
			out.WriteString(fmt.Sprintf(" (%s)", currency.Symbol))
		}
	} else {
		if code && symbol {
			out.WriteString(fmt.Sprintf("%s (%s)", currency.Code, currency.Symbol))
		} else if code {
			out.WriteString(currency.Code)
		} else if symbol {
			out.WriteString(currency.Symbol)
		}
	}
	return out.String()
}

func languagesString(languages []ipdata.Language) string {
	var out strings.Builder

	for _, language := range languages {
		out.WriteString(language.Name)
		out.WriteString(", ")
	}

	return strings.TrimSuffix(out.String(), ", ")
}

func threatString(threat ipdata.Threat) string {
	var out strings.Builder
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
			out.WriteString(threat)
			out.WriteString(", ")
		}
		return strings.TrimSuffix(out.String(), ", ")
	} else {
		return "None"
	}
}

func printIPData(data ipdata.Data) {
	var out strings.Builder

	if opt.JSON {
		out.WriteString(*data.JSON)
	} else {
		out.WriteString(fmt.Sprintf("IP: %v\n", data.IP))

		if len(data.Postal) > 0 {
			if len(data.Region) > 0 {
				out.WriteString(fmt.Sprintf("   %v %v\n", data.Postal, data.City))
			} else {
				out.WriteString(fmt.Sprintf("   %v\n", data.Postal))
			}
		} else {
			if len(data.Region) > 0 {
				out.WriteString(fmt.Sprintf("   %v\n", data.City))
			}
		}

		if len(data.Region) > 0 {
			out.WriteString(fmt.Sprintf("   %v\n", data.Region))
		}

		if len(data.CountryName) > 0 {
			out.WriteString("   " + data.CountryName)
			if len(data.CountryCode) > 0 {
				out.WriteString(fmt.Sprintf(" (%v)", data.CountryCode))
			}
			out.WriteString("\n")
		}

		if len(data.ContinentName) > 0 {
			out.WriteString("   " + data.ContinentName)
			if len(data.ContinentCode) > 0 {
				out.WriteString(fmt.Sprintf(" (%v)", data.ContinentCode))
			}
			out.WriteString("\n")
		}

		out.WriteString(fmt.Sprintf("   Coordinates:     %g, %g\n\n", data.Latitude, data.Longitude))

		if len(data.Flag) > 0 {
			out.WriteString(fmt.Sprintf("   Flag:            %v\n", getFlagRepr(data)))
		}

		timezone := timeZoneString(data.TimeZone)
		if len(timezone) > 0 {
			out.WriteString(fmt.Sprintf("   Time zone:       %v\n", timezone))
		}

		currency := currencyString(data.Currency)
		if len(currency) > 0 {
			out.WriteString(fmt.Sprintf("   Currency:        %v\n", currency))
		}

		languages := languagesString(data.Languages)
		if len(languages) > 0 {
			out.WriteString(fmt.Sprintf("   Languages:       %v\n", languages))
		}

		if len(data.CallingCode) > 0 {
			out.WriteString(fmt.Sprintf("   Calling code:    +%v\n", data.CallingCode))
		}
		out.WriteString("\n")

		if len(data.ASN.Name) > 0 {
			out.WriteString("   Organization:    " + data.ASN.Name)
			if len(data.ASN.Domain) > 0 {
				out.WriteString(" (" + data.ASN.Domain + ")")
			}
			out.WriteString("\n")
		}

		if len(data.ASN.ASN) > 0 {
			out.WriteString("   AS number:       " + data.ASN.ASN + "\n")
		}

		out.WriteString("   Threat:          " + threatString(data.Threat) + "\n\n")
	}

	fmt.Print(out.String())
}
