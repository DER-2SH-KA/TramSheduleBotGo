package parser

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Parser struct {
	Client *http.Client
}

const (
	URL_BASIC                 = "https://nskgortrans.ru/components/com_planrasp/helpers/grasp.php?tv=mr&m=%s&t=3&r=%s&sch=%s&s=0&v=0"
	ROUTE_2                   = "0002"
	ROUTE_8                   = "0008"
	DIRECTION_FROM_CH_SLOBODA = "A"
	DIRECTION_TO_CH_SLOBODA   = "B"
	WORK_DAYS                 = "11"
	WEEKENDS                  = "5"
)

func (p *Parser) Parse(route, days, direction string) (string, error) {
	schedule, err := p.parse(route, days, direction)
	if err != nil {
		return "", err
	}

	return schedule, nil
}

func (p *Parser) parse(route, days, direction string) (string, error) {
	resp, err := p.Client.Get(fmt.Sprintf(URL_BASIC, route, days, direction))
	if err != nil {
		return "", fmt.Errorf("failed to fetch url: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch url: response status code %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get document from response body: %w", err)
	}

	schedule := doc.Find("td[align='center']")
	if schedule == nil {
		return "", fmt.Errorf("failed to find schedule in response body")
	}

	station := schedule.Find("h2").First()
	if station == nil {
		return "", fmt.Errorf("failed to find station in schedule")
	}

	scheduleAsString := ""

	hours, minutes := make([]string, 0, 24), make([]string, 0, 24)

	schedule.Find("table[class='tbl_plan_rasp']").First().Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(i int, s *goquery.Selection) {

			if s.HasClass("td_plan_n") || strings.TrimSpace(s.Text()) == "" {
				return
			}

			if s.HasClass("td_plan_h") {
				hours = append(hours, strings.TrimSpace(s.Text()))
				return
			}

			if s.HasClass("td_plan_m") {
				minutesForCurrentHour := ""

				s.Find("div").Each(func(i int, s *goquery.Selection) {
					minutesForCurrentHour += strings.TrimSpace(s.Text()) + " "
				})

				minutes = append(minutes, minutesForCurrentHour)
				return
			}
		})
	})

	scheduleAsString = station.Text() + "\n"
	for i, _ := range hours {
		scheduleAsString += fmt.Sprintf("Час: %s. Минуты: %s\n", hours[i], minutes[i])
	}

	return scheduleAsString, nil
}

func htmlToText(s *goquery.Selection) string {
	html, err := goquery.OuterHtml(s)
	if err != nil {
		log.Println(err)
	}
	return html
}
