package timetable

import (
	"encoding/json"
	//"fmt"

	"io"

	"github.com/goferwplynie/librusApi/client"
	"github.com/goferwplynie/librusApi/internal/models/timetable"
)

func GetTimeTable(c *client.Client) (*timetable.TimetableResponse, error) {
	res, err := c.Get("TimetableEntries")
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	//fmt.Println(string(body))

	var timeTable timetable.TimetableResponse
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &timeTable)
	if err != nil {
		return nil, err
	}

	return &timeTable, nil
}
