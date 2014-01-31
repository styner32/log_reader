package main

import "os"
import "log"
import "fmt"
import "bufio"
import "regexp"
import "encoding/json"

type Header struct {
	Month string
	Date string
	Time string
	Ip string
	Uuid string
	Username string
}

type Body struct {
	Length int
	Contents []string
}

func main() {
	file, err := os.Open("example.log")
	var logs [128]string
	headers := make(map[string]*Header)
	bodies := make(map[string]*Body)
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

		if len(header.Uuid) == 0 {
			continue
		}

		if value,ok := headers[header.Uuid]; ok {
			bodies[value.Uuid].AddContent(body)
		} else {
			headers[header.Uuid] = header
			bodies[header.Uuid] = &Body {Length: 0, Contents: make([]string, 10)}
			bodies[header.Uuid].AddContent(body)
		}
	}

	resultFile, err := os.Create("result.json")
 	if err != nil {
		log.Fatal(err)
 	}
 	defer resultFile.Close()

	for uuid, header := range headers {
   		header_in_json, err := json.Marshal(header)
     	if err != nil {
			log.Fatal(err)
     	}

     	headerLength, err := resultFile.Write(header_in_json)
		if err != nil {
			log.Fatal(err)
     	}
     	fmt.Println(headerLength)

    	body_in_json, err := json.Marshal(bodies[uuid])
     	if err != nil {
			log.Fatal(err)
     	}

		bodyLength, err := resultFile.Write(body_in_json)
		if err != nil {
			log.Fatal(err)
     	}
     	fmt.Println(bodyLength)
    }

    fmt.Println("Done!")
}

func ParseLog(str string) (*Header, string) {
	expr := `(\w{3}) (\d+) (\d{2}:\d{2}:\d{2}) ([\d\.]+) production.log: \[([\w\d-]+)\] (\[([\w\d]+)\])? (.*)`
	log, _ := regexp.Compile(expr)
	if len(log.FindString(str)) > 0 {
		matches := log.FindStringSubmatch(str)
		header := &Header{
			Month: matches[1],
			Date: matches[2],
			Time: matches[3],
			Ip: matches[4],
			Uuid: matches[5],
			Username: matches[7] }
		return header, matches[8]
	}
	return &Header{}, ""
}

func (body *Body) AddContent(content string) {
	body.Contents[body.Length] = content
	body.Length = body.Length + 1
}
