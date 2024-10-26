package conductorworker

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"

	"github.com/conductor-sdk/conductor-go/sdk/model"
)

func Number(task *model.Task) (interface{}, error) {
	var err error
	var number float64

	_, ok := task.InputData["number"].(string)
	if ok {
		number, err = strconv.ParseFloat(task.InputData["number"].(string), 64)
		if err != nil {
			log.Errorf("Couldn't convert number to int: %v", err)
			return nil, err
		}
	} else {
		number = task.InputData["number"].(float64)
	}

	return map[string]interface{}{
		"number": number,
	}, nil
}

func Square(task *model.Task) (interface{}, error) {

	number := task.InputData["number"].(float64)
	squareResult := number * number

	return map[string]interface{}{
		"square": squareResult,
	}, nil
}

func Sleepms(task *model.Task) (interface{}, error) {
	sleep := int(task.InputData["square"].(float64))
	time.Sleep(time.Duration(sleep) * time.Millisecond)

	return map[string]interface{}{
		"sleep": sleep,
	}, nil
}
