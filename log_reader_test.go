package log_reader

import "testing"
import "reflect"

func CompareHeaderStruct(t *testing.T, expected_result *Header, actual_result *Header) {
	fields := [...]string{"Month", "Date", "Time", "Uuid", "Username"}
	for _, field := range fields {
		expected := reflect.ValueOf(*expected_result).FieldByName(field).String()
		result := reflect.ValueOf(*actual_result).FieldByName(field).String()
		if expected != result {
			t.Fatalf("Expected %s, got %s in %s", expected, result, field)
		}
	}
}

func CompareRequestStartStruct(t *testing.T, expected_result *RequestStartBody, actual_result *RequestStartBody) {
	fields := [...]string{"Action", "Url", "Ip", "Date", "Time"}
	for _, field := range fields {
		expected := reflect.ValueOf(*expected_result).FieldByName(field).String()
		result := reflect.ValueOf(*actual_result).FieldByName(field).String()
		if expected != result {
			t.Fatalf("Expected %s, got %s in %s", expected, result, field)
		}
	}
}

func TestParseHeader(t *testing.T) {
	input := "Jan 31 07:24:33 184.173.146.35 production.log: [1c0bd418-61fb-4b1d-9501-5fbe69d4d50f] [user1] Processing by V4::SubtitlesController#show as HTML"
	header, body := ParseHeader(input)

	expected_header := &Header{
		Month:    "Jan",
		Date:     "31",
		Time:     "07:24:33",
		Ip:       "184.173.146.35",
		Uuid:     "1c0bd418-61fb-4b1d-9501-5fbe69d4d50f",
		Username: "user1"}

	CompareHeaderStruct(t, expected_header, header)
	if body != " Processing by V4::SubtitlesController#show as HTML" {
		t.Fatalf("Expected Processing by V4::SubtitlesController#show as HTML, got %s in body", body)
	}
}

func TestParseHeaderWithoutUsername(t *testing.T) {
	input := "Jan 31 07:24:33 184.173.146.35 production.log: [1c0bd418-61fb-4b1d-9501-5fbe69d4d50f] Processing by V4::SubtitlesController#show as HTML"
	header, body := ParseHeader(input)

	expected_header := &Header{
		Month:    "Jan",
		Date:     "31",
		Time:     "07:24:33",
		Ip:       "184.173.146.35",
		Uuid:     "1c0bd418-61fb-4b1d-9501-5fbe69d4d50f" }

	CompareHeaderStruct(t, expected_header, header)
	if body != "Processing by V4::SubtitlesController#show as HTML" {
		t.Fatalf("Expected Processing by V4::SubtitlesController#show as HTML, got %s in body", body)
	}
}

func TestRequestStartBody(t *testing.T) {
	input := "Started GET \"/v4/videos/1024918v/subtitles/en.srt?app=65535a&t=1391153071&site=www.viki.com&token=mfQNWY6DNFsJWmEC8NBLqe9X_02&sig=29a39a663074d1f443b48fb439afb22e2cb47650\" for 192.69.221.178 at 2014-01-31 07:24:33 +0000"
	body := ParseRequestStartBody(input)

	expected_body := &RequestStartBody {
		Action: "GET",
		Url: "/v4/videos/1024918v/subtitles/en.srt?app=65535a&t=1391153071&site=www.viki.com&token=mfQNWY6DNFsJWmEC8NBLqe9X_02&sig=29a39a663074d1f443b48fb439afb22e2cb47650",
		Ip: "192.69.221.178",
		Date: "2014-01-31",
		Time: "07:24:33 +0000" }

	CompareRequestStartStruct(t, expected_body, body)
}
