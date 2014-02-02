package log_reader

import "testing"
import "reflect"

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

	header_fields := [...]string{"Month", "Date", "Time", "Uuid", "Username"}

	for _, header_field := range header_fields {
		expected := reflect.ValueOf(*expected_header).FieldByName(header_field).String()
		result := reflect.ValueOf(*header).FieldByName(header_field).String()
		if expected != result {
			t.Fatalf("Expected %s, got %s in %s", expected, result, header_field)
		}
	}

	if body != " Processing by V4::SubtitlesController#show as HTML" {
		t.Fatalf("Expected Processing by V4::SubtitlesController#show as HTML, got %s in body", body)
	}
}
