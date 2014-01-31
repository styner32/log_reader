package main

import "os"
import "log"
import "fmt"
import "bufio"
import "regexp"
import "encoding/json"
import "strconv"

type Header struct {
	Month string
	Date string
	Time string
	Ip string
	Uuid string
	Username string
}

type Body struct {
	ReqeustStart *RequestStartBody
	Processor *ProcessorBody
	Parameters *ParametersBody
	Complete *CompleteBody
	Length int
	Contents []string
}

type RequestStartBody struct {
	Action string
	Url string
	Ip string
	Date string
	Time string
}

type ProcessorBody struct {
	ControllerName string
	Action string
	MimeType string
}

type ParametersBody struct {
	Parameters string
}

type CompleteBody struct {
	StatusCode string
	TotalTime float64
	DatabaseTime float64
	ViewRenderTime float64
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

		if header == nil {
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
     	resultFile.WriteString("\n")
     	fmt.Println(headerLength)

    	body_in_json, err := json.Marshal(bodies[uuid])
     	if err != nil {
			log.Fatal(err)
     	}

		bodyLength, err := resultFile.Write(body_in_json)
		if err != nil {
			log.Fatal(err)
     	}
     	resultFile.WriteString("\n")
     	fmt.Println(bodyLength)
    }
    resultFile.Sync()

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
	return nil, ""
}

func ParseRequestStartBody(str string) (*RequestStartBody) {
	expr := `Started (\w+) "([^"]+)" for ([\d\.]+) at (\d{4}-\d{2}-\d{2}) (\d{2}:\d{2}:\d{2} [+-]\d{4})`
	log, _ := regexp.Compile(expr)
	if len(log.FindString(str)) > 0 {
		matches := log.FindStringSubmatch(str)
		request_start_body := &RequestStartBody{
			Action: matches[1],
			Url: matches[2],
			Ip: matches[3],
			Date: matches[4],
			Time: matches[5] }
		return request_start_body
	}
	return nil
}

func ParseProcessorBody(str string) (*ProcessorBody) {
	expr := `Processing by ([:\w]+)#(\w+) as (.*)`
	log, _ := regexp.Compile(expr)
	if len(log.FindString(str)) > 0 {
		matches := log.FindStringSubmatch(str)
		processor_body := &ProcessorBody{
			ControllerName: matches[1],
			Action: matches[2],
			MimeType: matches[3] }
		return processor_body
	}
	return nil
}

func ParseParametersBody(str string) (*ParametersBody) {
	expr := `Parameters: (.*)`
	log, _ := regexp.Compile(expr)
	if len(log.FindString(str)) > 0 {
		matches := log.FindStringSubmatch(str)
		parameter_body := &ParametersBody{ Parameters: matches[1] }
		return parameter_body
	}
	return nil
}

func ParseCompleteBody(str string) (*CompleteBody) {
	expr := `Completed (\d+) [\w\s]+ in ([\d\.]+)m?s \(Views: ([\d\.]+)m?s \| ActiveRecord: ([\d\.]+)m?s\)`
	log, _ := regexp.Compile(expr)
	if len(log.FindString(str)) > 0 {
		matches := log.FindStringSubmatch(str)
		total_time, _ := strconv.ParseFloat(matches[2], 64)
		database_time, _ := strconv.ParseFloat(matches[3], 64)
		view_render_time, _ := strconv.ParseFloat(matches[4], 64)
		complete_body := &CompleteBody{
			StatusCode: matches[1],
			TotalTime: total_time,
			DatabaseTime: database_time,
			ViewRenderTime: view_render_time }
		return complete_body
	}
	return nil
}

func (body *Body) AddContent(content string) {
	body.Contents[body.Length] = content
	body.Length = body.Length + 1
	request_start_body := ParseRequestStartBody(content)
	if request_start_body != nil {
		body.ReqeustStart = request_start_body
		return
	}
	processor_body := ParseProcessorBody(content)
	if processor_body != nil {
		body.Processor = processor_body
		return
	}
	parameters_body := ParseParametersBody(content)
	if parameters_body != nil {
		body.Parameters = parameters_body
		return
	}
	complete_body := ParseCompleteBody(content)
	if complete_body != nil {
		body.Complete = complete_body
		return
	}
}
