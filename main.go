package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
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
	countryCode    string
	countryPattern []byte
)

func main() {

	if len(os.Args) < 2 {
		log.Error("Missing source file")
		log.Info("Usage: ", os.Args[0], " FILE [CONTRY]")
		log.Info("E.g., ", os.Args[0], " ORCID-API-2.0_activities_xml_10_2018.tar.gz NZ")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	countryCode = os.Args[2]
	countryPattern = []byte(">" + countryCode + "<")

	zr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
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
			if strings.HasSuffix(fn, ".xml") {
				recBytes, _ := ioutil.ReadAll(tr)
				wg.Add(1)
				go handleFile(fn, recBytes, &wg)
			}
		default:
			log.Infof("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				h.Typeflag,
				"in file",
				fn,
			)
		}
	}
	wg.Wait()
}

func handleFile(fn string, content []byte, wg *sync.WaitGroup) {
	if bytes.Contains(content, []byte(countryPattern)) {
		var rec record
		xml.Unmarshal(content, &rec)
		if rec.Organization.Address.Country == countryCode || rec.ConveningOrganization.Address.Country == countryCode {
			log.Info(fn)
			err := ioutil.WriteFile(filepath.Base(fn), content, 0644)
			if err != nil {
				log.Error(err)
			}
		}
	}
	wg.Done()
}
