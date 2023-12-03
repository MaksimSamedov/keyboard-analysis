package flow

import (
	"errors"
	"keyboard-analysis/internal/models"
)

var ErrInvalidEventFlow = errors.New("invalid event flow")

type KeyboardAnalysis struct {
	Time      []int
	Intervals []int
}

func (ka *KeyboardAnalysis) Inline() []int {
	l := len(ka.Time) + len(ka.Intervals)
	var res = make([]int, l)
	for i, val := range ka.Time {
		res[2*i] = val
	}
	for i, val := range ka.Intervals {
		res[2*i+1] = val
	}
	return res
}

func AnalyseEvents(flow *models.KeyboardFlow) (*KeyboardAnalysis, error) {
	var events []*models.KeyboardEvent

	// соберём последовательность клавиш по моменту их нажатия
	var pattern []*Pattern
	var remain []*models.KeyboardEvent
	for _, ev := range flow.Flow {
		if !ev.Up {
			pattern = append(pattern, &Pattern{down: ev})
			events = append(events, ev)
		} else {
			remain = append(remain, ev)
		}
	}

	// проверим, что слово было набрано
	pattern, err := OnlyCorrectPatterns(flow, pattern)
	if err != nil {
		return nil, err
	}

	// запишем в паттерн события отпускания клавиш
	for _, pat := range pattern {
		var keyFound bool
		for j, ev := range remain {
			if ev.Up && ev.Key == pat.down.Key {
				keyFound = true
				pat.up = ev
				remain = append(remain[:j], remain[j+1:]...)
				break
			}
		}
		if !keyFound {
			return nil, ErrInvalidEventFlow
		}
	}

	var res KeyboardAnalysis
	var last *Pattern
	for _, pat := range pattern {
		if last != nil {
			res.Intervals = append(res.Intervals, int(pat.down.Time)-int(last.up.Time))
		}
		res.Time = append(res.Time, int(pat.up.Time)-int(pat.down.Time))
		last = pat
	}

	return &res, nil
}

type Pattern struct {
	up   *models.KeyboardEvent
	down *models.KeyboardEvent
}

var ErrPhraseIsNotComplete = errors.New("phrase from patterns is not complete")

func OnlyCorrectPatterns(flow *models.KeyboardFlow, patterns []*Pattern) ([]*Pattern, error) {
	var res []*Pattern

	phrase := flow.Password.Password
	phraseLength := len(phrase)
	letter := phrase[0]
	letterIndex := 0
	for _, pattern := range patterns {
		if pattern.down.Key == string(letter) {
			res = append(res, pattern)
			letterIndex++
			if letterIndex == phraseLength {
				break
			}
			letter = phrase[letterIndex]
		}
	}
	if phraseLength != letterIndex {
		return nil, ErrPhraseIsNotComplete
	}

	return res, nil
}

//type KeyboardVectorElement struct {
//	Y          int
//	X          int
//	Time       uint
//	IsInterval bool
//}

//type KeyboardClick struct {
//	Duration   int
//	IsInterval bool
//}

//func EventsToVector(events []*models.KeyboardEvent) ([]KeyboardVectorElement, error) {
//	var flow []*models.KeyboardEvent
//	flow = append(flow, events...)
//
//	i := 1
//	l := len(flow)
//	for i < l {
//		if flow[i-1].Time > flow[i].Time {
//			flow[i-1], flow[i] = flow[i], flow[i-1]
//			i--
//		} else {
//			i++
//		}
//	}
//
//	var res []KeyboardVectorElement
//	var lastTime uint = 0
//	var stack []string
//	for i, ev := range flow {
//		delta := ev.Time - lastTime
//		lastTime = ev.Time
//
//		var track *KeyboardVectorElement
//		if ev.Up {
//			// отпускание клавиши
//			sl := len(stack)
//			if sl == 0 {
//				fmt.Println("err empty stack", i, ev.Key, strings.Join(stack, ", "))
//				return nil, ErrInvalidEventFlow
//			}
//
//			removeFromStack := -1
//			for j, pressed := range stack {
//				if ev.Key == pressed {
//					removeFromStack = j
//					break
//				}
//			}
//			if removeFromStack < 0 {
//				fmt.Println("err key not found in stack", ev.Key, strings.Join(stack, ", "))
//				return nil, ErrInvalidEventFlow
//			}
//
//			track = &KeyboardVectorElement{
//				Y:          len(stack),
//				X:          int(delta),
//				Time:       ev.Time,
//				IsInterval: false,
//			}
//			stack = append(stack[:removeFromStack], stack[removeFromStack+1:]...)
//		} else {
//			// нажатие клавиши
//			if len(stack) == 0 {
//				track = &KeyboardVectorElement{
//					Y:          len(stack),
//					X:          int(delta),
//					Time:       ev.Time,
//					IsInterval: true,
//				}
//			}
//			stack = append(stack, ev.Key)
//		}
//		// фиксируем промежуток
//		if i != 0 && track != nil {
//			res = append(res, *track)
//		}
//		fmt.Println("stack: ", strings.Join(stack, ", "))
//	}
//	if len(stack) != 0 {
//		fmt.Println("err unused elements", strings.Join(stack, ", "))
//		return nil, ErrInvalidEventFlow
//	}
//
//	return res, nil
//}
