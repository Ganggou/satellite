package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
)

type Satellite struct {
	SatelliteName   string
	ApogeeAltitude  string
	PerigeeAltitude string
	Inclination     string
	ArgOfPerigee    string
	Meannomaly      string
	RAAN            string
	OrbitEpoch      string
}

const filename = "output"
const limitNum = 50

func Fetch() {
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"

	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Println(err)
		}
		var satellite = new(Satellite)
		nameNode := htmlquery.Find(doc, `//*[@id="ctl00_lblTitle"]`)
		satellite.SatelliteName = strings.Split(htmlquery.InnerText(nameNode[0]), " ")[0]
		if !strings.Contains(satellite.SatelliteName, "STARLINK") {
			return
		}

		tableNodes := htmlquery.Find(doc, `//table[3]//tr//td[2]`)
		satellite.OrbitEpoch = htmlquery.InnerText(tableNodes[0])
		satellite.Inclination = strings.ReplaceAll(htmlquery.InnerText(tableNodes[2]), "°", "")
		satellite.PerigeeAltitude = strings.Split(htmlquery.InnerText(tableNodes[3]), " ")[0]
		satellite.ApogeeAltitude = strings.Split(htmlquery.InnerText(tableNodes[4]), " ")[0]
		satellite.RAAN = strings.ReplaceAll(htmlquery.InnerText(tableNodes[5]), "°", "")
		satellite.ArgOfPerigee = strings.ReplaceAll(htmlquery.InnerText(tableNodes[6]), "°", "")
		satellite.Meannomaly = strings.ReplaceAll(htmlquery.InnerText(tableNodes[8]), "°", "")
		satellite.Print()
	})

	var limit = make(chan int, limitNum)
	for i := 0; i < limitNum; i++ {
		limit <- 1
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString("<SatelliteStore>\n"); err != nil {
		panic(err)
	}
	for satelliteID := 44238; satelliteID < 47182; satelliteID++ {
		go func(id int) {
			<-limit
			c.Visit(fmt.Sprintf("https://www.heavens-above.com/orbit.aspx?satid=%v", id))
			limit <- 1
		}(satelliteID)
	}
	for i := 0; i < limitNum; i++ {
		<-limit
	}
	if _, err = f.WriteString("</SatelliteStore>"); err != nil {
		panic(err)
	}
}

func (s *Satellite) Print() {
	content := fmt.Sprintf(`	<Satellite>
		<SatelliteName>%v</SatelliteName>
		<ApogeeAltitude>%v</ApogeeAltitude>
		<PerigeeAltitude>%v</PerigeeAltitude>
		<Inclination>%v</Inclination>
		<ArgOfPerigee>%v</ArgOfPerigee>
		<Meannomaly>%v</Meannomaly>
		<RAAN>%v</RAAN>
		<OrbitEpoch>%v</OrbitEpoch>
	</Satellite>
`, s.SatelliteName, s.ApogeeAltitude, s.PerigeeAltitude, s.Inclination, s.ArgOfPerigee, s.Meannomaly, s.RAAN, s.OrbitEpoch)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		panic(err)
	}
}

func main() {
	Fetch()
}
