package station

import (
	"encoding/json"
	"errors"
	"mrt-schedule/common/client"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	GetAllStation() (response []StationResponse, err error)
	CheckScheduleByStation(id string) (response []ScheduleResponse, err error)
}

type service struct {
	client *http.Client
}

func newService() Service {
	return &service{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *service) GetAllStation() (response []StationResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var station []Station
	err = json.Unmarshal(byteResponse, &station)

	for _, item := range station {
		response = append(response, StationResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}
	return
}

func (s *service) CheckScheduleByStation(id string) (response []ScheduleResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var schedule []Schedule
	err = json.Unmarshal(byteResponse, &schedule)
	if err != nil {
		return
	}

	var scheduleSelected Schedule
	for _, item := range schedule {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}
	if scheduleSelected.StationId == "" {
		err = errors.New("Station not found")
		return
	}

	response, err = ConvertDataToResponse(scheduleSelected)
	if err != nil {
		return
	}
	return
}

func ConvertDataToResponse(schedule Schedule) (response []ScheduleResponse, err error) {
	var LebakBulusTripName = "Stasiun Lebak Bulus Grab"
	var BundaranHITripName = "Stasiun Bundaran HI Bank DKI"

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHi := schedule.ScheduleBundaranHi

	scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return
	}

	scheduleBundaranHiParsed, err := ConvertScheduleToTimeFormat(scheduleBundaranHi)
	if err != nil {
		return
	}

	for _, item := range scheduleLebakBulusParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	for _, item := range scheduleBundaranHiParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHITripName,
				Time:        item.Format("15:04"),
			})
		}
	}
	return
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var parsedTime time.Time
	var schedules = strings.Split(schedule, ",")

	for _, item := range schedules {
		trimmedTime := strings.TrimSpace(item)
		if trimmedTime == "" {
			continue
		}
		parsedTime, err = time.Parse("15:04", trimmedTime)

		if err != nil {
			err = errors.New("Invalid time format" + trimmedTime)
			return
		}

		response = append(response, parsedTime)
	}

	return
}
