package main

import (
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const configURL = "https://www.nhk.or.jp/radio/config/config_web.xml"

type stream struct {
	XMLName xml.Name `xml:"data"`
	AreaJP  string   `xml:"areajp"`
	Area    string   `xml:"area"`
	APIKey  string   `xml:"apikey"`
	AreaKey string   `xml:"areakey"`
	R1HLS   string   `xml:"r1hls"`
	R2HLS   string   `xml:"r2hls"`
	FMHLS   string   `xml:"fmhls"`
}

type streams []stream

func (s streams) Find(area, stream string) (string, error) {
	for _, i := range s {
		if i.Area == area || i.AreaJP == area {
			switch stream {
			case "r1":
				return i.R1HLS, nil
			case "r2":
				return i.R2HLS, nil
			case "fm":
				return i.FMHLS, nil
			default:
				return "", errors.New("invalid stream requested")
			}
		}
	}
	return "", errors.New("invalid area requested")
}

type config struct {
	XMLName xml.Name `xml:"radiru_config"`

	// お知らせ
	Info string `xml:"info"`

	// 各地域のストリームURL
	StreamURL streams `xml:"stream_url>data"`

	// noa api
	URLProgramNOA string `xml:"url_program_noa"`

	// program detail api
	URLProgramDay string `xml:"url_program_day"`

	// program info api
	URLProgramDetail string `xml:"url_program_detail"`

	// tweet cgi @radiru
	RadiruTwitterTimeline string `xml:"radiru_twitter_timeline"`
}

func main() {
	area, stream := "tokyo", "r1"

	if len(os.Args) > 1 {
		area = strings.ToLower(os.Args[1])
	}

	if len(os.Args) > 2 {
		stream = strings.ToLower(os.Args[2])
	}

	resp, err := http.Get(configURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var data config
	if err := xml.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	url, err := data.StreamURL.Find(area, stream)
	if err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("mpv", "--really-quiet", url).Run(); err != nil {
		log.Fatal(err)
	}
}
