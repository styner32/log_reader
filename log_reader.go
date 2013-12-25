package main

import "os"
import "log"
import "fmt"
import "bufio"
import "regexp"

type Header struct {
	month string
	date string
	time string
	ip string
	uuid string
	username string
}

type Body struct {
	length int
	contents []string
}

func main() {
	file, err := os.Open("example.log")
	var logs [128]string
	headers := make(map[string]Header)
	bodies := make(map[string]Body)
	var log_count int

	if err != nil {
		log.Fatal(err)
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i := 0; i < len(lines); i++ {
		logs[log_count] = string(lines[i])
		log_count += 1
	}

	for i := 0; i < len(logs); i++ {
		header, body := ParseLog(logs[i])
		if value,ok := headers[header.uuid]; ok {
			bodies[value.uuid].AddContent(body)
		} else {
			headers[header.uuid] = header
			bodies[header.uuid] = Body {
				length: 0,
				contents: make([]string, 5) }
			bodies[header.uuid].AddContent(body)
		}
	}

	fmt.Println(headers)
	fmt.Println(bodies)
}

func ParseLog(str string) (Header, string) {
	expr := `(\w{3}) (\d+) (\d{2}:\d{2}:\d{2}) ([\d\.]+) production.log: \[([\w\d-]+)\] (\[([\w\d]+)\])? (.*)`
	log, _ := regexp.Compile(expr)
	if len(log.FindString(str)) > 0 {
		matches := log.FindStringSubmatch(str)
		header := Header{
			month: matches[1],
			date: matches[2],
			time: matches[3],
			ip: matches[4],
			uuid: matches[5],
			username: matches[7] }
		return header, matches[8]
	}
	return Header{}, ""
}

func (body Body) AddContent(content string) {
	body.contents[body.length] = content
	body.length = body.length + 1

}
