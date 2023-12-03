package flow

import (
	"errors"
	"keyboard-analysis/internal/models"
)

type Analyser struct {
	tasks []*Task
	props AnalyserProps
}

func NewAnalyser(props AnalyserProps) *Analyser {
	return &Analyser{
		props: props,
	}
}

func (analyser *Analyser) AddTask(flow *models.KeyboardFlow, samples []*models.KeyboardFlow) {
	analyser.tasks = append(analyser.tasks, &Task{
		Flow:    flow,
		Samples: samples,
		props:   analyser.props,
	})
}

var ErrNoTasksToAnalyze = errors.New("no tasks to analyze")
var ErrInvalidDeviation = errors.New("'max deviation' must be between 0 and 100")
var ErrInvalidErrorsInTask = errors.New("'max errors in task' must be between 0 and 100")

func (analyser *Analyser) Analyse() (bool, error) {
	if err := analyser.props.Validate(); err != nil {
		return false, err
	}

	l := len(analyser.tasks)
	if l == 0 {
		return false, ErrNoTasksToAnalyze
	}

	//wg := sync.WaitGroup{}
	//wg.Add(l)
	for _, task := range analyser.tasks {
		//go func(task *Task) {
		task.Analyse()
		//	wg.Done()
		//}(task)
	}
	//wg.Wait()

	for _, task := range analyser.tasks {
		if !task.Success() {
			return false, task.Error()
		}
	}
	return true, nil
}

type AnalyserProps struct {
	MaxDeviation             float64 // максимальное допустимое среднее отклонение промежутка времени (в процентах)
	MinSuccessfulComparisons float64 // минимальная доля успешных сравнений почерка (в процентах)
	MaxErrorsInTask          float64 // максимальная доля ошибок при анализе одного таска (в процентах)
}

func (props AnalyserProps) Validate() error {
	if props.MaxDeviation <= 0 || props.MaxDeviation >= 100 {
		return ErrInvalidDeviation
	}
	if props.MaxErrorsInTask <= 0 || props.MaxErrorsInTask >= 100 {
		return ErrInvalidErrorsInTask
	}
	return nil
}
