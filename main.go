package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type record struct {
	XMLName            xml.Name // `xml:"employment"`
	Text               string   `xml:",chardata"`
	Internal           string   `xml:"internal,attr"`
	Address            string   `xml:"address,attr"`
	Email              string   `xml:"email,attr"`
	History            string   `xml:"history,attr"`
	Employment         string   `xml:"employment,attr"`
	Person             string   `xml:"person,attr"`
	Education          string   `xml:"education,attr"`
	OtherName          string   `xml:"other-name,attr"`
	PersonalDetails    string   `xml:"personal-details,attr"`
	Bulk               string   `xml:"bulk,attr"`
	Common             string   `xml:"common,attr"`
	Record             string   `xml:"record,attr"`
	Keyword            string   `xml:"keyword,attr"`
	Activities         string   `xml:"activities,attr"`
	Deprecated         string   `xml:"deprecated,attr"`
	ExternalIdentifier string   `xml:"external-identifier,attr"`
	Funding            string   `xml:"funding,attr"`
	Error              string   `xml:"error,attr"`
	Preferences        string   `xml:"preferences,attr"`
	Work               string   `xml:"work,attr"`
	ResearcherURL      string   `xml:"researcher-url,attr"`
	PeerReview         string   `xml:"peer-review,attr"`
	PutCode            string   `xml:"put-code,attr"`
	Path               string   `xml:"path,attr"`
	Visibility         string   `xml:"visibility,attr"`
	CreatedDate        string   `xml:"created-date"`
	LastModifiedDate   string   `xml:"last-modified-date"`
	RoleTitle          string   `xml:"role-title"`
	Organization       struct {
		Text    string `xml:",chardata"`
		Name    string `xml:"name"`
		Address struct {
			Text    string `xml:",chardata"`
			City    string `xml:"city"`
			Region  string `xml:"region"`
			Country string `xml:"country"`
		} `xml:"address"`
	} `xml:"organization"`
	ConveningOrganization struct {
		Text    string `xml:",chardata"`
		Name    string `xml:"name"`
		Address struct {
			Text    string `xml:",chardata"`
			City    string `xml:"city"`
			Country string `xml:"country"`
		} `xml:"address"`
	} `xml:"convening-organization"`
}

var (
	countryCode                                      string
	countryPattern                                   []byte
	recordType, dest, searchString, searchExpression string
	searchPattern                                    []byte
	searchRE                                         *regexp.Regexp
)

func main() {

	app := cli.NewApp()
	app.Name = "extract-orcid"
	app.Usage = `extract filtered data from ORCID profile activity public data`
	app.Version = "1.1.0"
	app.ArgsUsage = "FILE"
	defaultDest, _ := os.Getwd()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "country,c",
			Value:       "NZ",
			Usage:       "the country the record is related to",
			Destination: &countryCode,
		},
		cli.StringFlag{
			Name:        "type,t",
			Usage:       "the record type: emp[ployment], edu[cation], work, fund[ing], peer[-review] ...",
			Destination: &recordType,
		},
		cli.StringFlag{
			Name:        "output,o",
			Usage:       "the output destination directory",
			Value:       defaultDest,
			Destination: &dest,
		},
		cli.StringFlag{
			Name:        "search,s",
			Usage:       "the search string",
			Destination: &searchString,
		},
		cli.StringFlag{
			Name:        "regex,r",
			Usage:       "the search regular expression (https://github.com/google/re2/wiki/Syntax)",
			Destination: &searchExpression,
		},
	}
	app.Action = func(c *cli.Context) error {
		return extract(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func extract(c *cli.Context) error {

	if c.NArg() < 1 {
		log.Info("Usage: ", os.Args[0], " FILE")
		log.Info("E.g., ", os.Args[0], " ORCID-API-2.0_activities_xml_10_2018.tar.gz NZ")
		return errors.New("missing source file")
	}

	f, err := os.Open(c.Args().Get(0))
	if err != nil {
		return err
	}

	if searchString != "" {
		searchPattern = []byte(searchString)
	}

	if searchExpression != "" {
		if searchRE, err = regexp.Compile(searchExpression); err != nil {
			log.WithError(err).Errorf("Failed to compile regular expression %q", searchExpression)
		}
	}

	// Used for pre-filter content to reduce time on xml parsing
	countryPattern = []byte(">" + countryCode + "<")

	zr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	tr := tar.NewReader(zr)

	var wg sync.WaitGroup

	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}

		fn := h.Name
		switch h.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if strings.HasSuffix(fn, ".xml") &&
				(recordType == "" || strings.Contains(fn, recordType)) {
				recBytes, _ := ioutil.ReadAll(tr)
				wg.Add(1)
				go handleFile(fn, recBytes, &wg)
			}
		default:
			log.Infof("%s : %c %s %s\n",
				"yikes! Unable to figure out type",
				h.Typeflag,
				"in file",
				fn,
			)
		}
	}
	wg.Wait()
	return nil
}

func handleFile(fn string, content []byte, wg *sync.WaitGroup) {
	if bytes.Contains(content, countryPattern) &&
		(searchPattern == nil || bytes.Contains(content, searchPattern)) &&
		(searchRE == nil || searchRE.Match(content)) {
		var rec record
		xml.Unmarshal(content, &rec)
		if rec.Organization.Address.Country == countryCode || rec.ConveningOrganization.Address.Country == countryCode {
			log.Info(fn)
			destFn := filepath.Join(dest, fn)
			err := os.MkdirAll(filepath.Dir(destFn), os.ModePerm)
			if err != nil {
				log.Error(err)
			}
			err = ioutil.WriteFile(destFn, content, 0644)
			if err != nil {
				log.Error(err)
			}
		}
	}
	wg.Done()
}
