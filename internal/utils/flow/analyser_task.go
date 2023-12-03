package flow

import (
	"errors"
	"fmt"
	"keyboard-analysis/internal/models"
	"math"
)

type Task struct {
	Flow    *models.KeyboardFlow
	Samples []*models.KeyboardFlow
	props   AnalyserProps
	success bool
	err     error
}

var ErrNoSamples = errors.New("no samples to compare with")

func (task *Task) Analyse() {
	// дополнительно провалидируем конфиг для анализа
	if err := task.props.Validate(); err != nil {
		task.err = err
		return
	}

	// обязательно проверим, что есть с чем сравнивать
	total := float64(len(task.Samples))
	if total == 0 {
		task.err = ErrNoSamples
		return
	}

	// проанализируем полученный поток
	flow, err := AnalyseEvents(task.Flow)
	if err != nil {
		task.err = err
		return
	}

	// проанализируем потоки из базы и сравним
	errCnt := float64(0)
	successCnt := float64(0)
	fmt.Println("Process: ", task.Flow.Password.Password)
	for _, sampleFlow := range task.Samples {
		sample, err := AnalyseEvents(sampleFlow)
		if err != nil {
			errCnt++
			continue
		}
		//fmt.Println("\ncompare", task.Flow.Password.Password, sampleFlow.Password.Password)
		if task.compare(*flow, *sample) {
			successCnt++
		}
	}
	fmt.Printf("successCnt/len %.1f/%.1f\n", successCnt, total)

	// получим процент успеха/ошибок
	successPercentage := (successCnt / total) * 100
	errorPercentage := (errCnt / total) * 100

	// определим успешность сравнения
	if errorPercentage > task.props.MaxErrorsInTask {
		task.success = false
	} else if successPercentage < task.props.MinSuccessfulComparisons {
		task.success = false
	} else {
		task.success = true
	}
}

func (task *Task) Success() bool {
	return task.success
}

func (task *Task) Error() error {
	return task.err
}

func (task *Task) compare(flowInfo, sampleInfo KeyboardAnalysis) bool {
	// получим линейные значения
	flow := flowInfo.Inline()
	sample := sampleInfo.Inline()
	//fmt.Println("flow:")
	//for _, x := range flow {
	//	fmt.Print(x, " ")
	//}
	//fmt.Println("\nsample:")
	//for _, x := range sample {
	//	fmt.Print(x, " ")
	//}
	//fmt.Println()

	// определим длину
	lFlow := len(flow)
	lSample := len(sample)
	//fmt.Println("len(Flow):", lFlow)
	//fmt.Println("len(Sample):", lSample)
	var l, errCnt int
	if lFlow < lSample {
		l = lFlow
		errCnt = lSample - lFlow
	} else {
		l = lSample
		errCnt = lFlow - lSample
	}

	// если длина слишком сильно не совпадает - не засчитываем
	errPercentage := float64(errCnt) / float64(l+errCnt)
	if errPercentage*100 > task.props.MaxErrorsInTask {
		//fmt.Println("DIFFERENT, reason=errPercentage", errPercentage)
		//fmt.Println("errCnt: ", errCnt)
		//fmt.Println("l+errCnt: ", l+errCnt)
		return false
	}
	matchPercentage := float64(l-errCnt) - errPercentage
	if matchPercentage*100 < task.props.MinSuccessfulComparisons {
		//fmt.Println("DIFFERENT, reason=matchPercentage", matchPercentage)
		//fmt.Println("l-errCnt", l-errCnt)
		//fmt.Println("errPercentage", errPercentage)
		return false
	}

	// считаем отклонения
	totalDeviation := float64(errCnt * 100)
	firstGot := float64(flow[0])
	firstWant := float64(sample[0])
	for i := 0; i < l; i++ {
		got := float64(flow[i])
		want := float64(sample[i])

		// пробуем перейти к коэффициентам
		got = got / firstGot
		want = want / firstWant

		// коэффициенты
		if want == 0 {
			want = 1
		}
		diff := (got - want) / want

		// разница сигмоид
		//diff := sig(got) - sig(want)

		// разница логарифмов
		//if got < 0 {
		//	got, want = got-2*got, want-2*got
		//}
		//if want < 0 {
		//	got, want = got-2*want, want-2*want
		//}
		//if got == 0 {
		//	got++
		//}
		//if want == 0 {
		//	want++
		//}
		//diff := math.Log10(got) - math.Log10(want)

		// сигмоид разницы
		//diff := got - want
		//diff = sig(diff / 100)

		// гип. тангенс
		//diff := gipTan(got - want)

		// Логарифм разницы
		//diff := got - want
		//if diff != 0 {
		//	diff = math.Log2(math.Abs(diff))
		//}

		// Логарифм множителя
		//add := math.Min(0, math.Min(got, want))
		//got, want = got+add, want+add
		//diff := got/want - 1
		//if diff != 0 {
		//	diff = math.Abs(diff)
		//}

		// разница логарифмов
		//if got < 0 {
		//	got, want = -got, want-2*got
		//}
		//if want < 0 {
		//	got, want = got-2*want, -want
		//}
		//if got == 0 {
		//	got++
		//}
		//if want == 0 {
		//	want++
		//}
		//diff := math.Log2(got) - math.Log2(want)

		//
		//min := math.Min(got, want)
		//max := math.Max(got, want)
		//if max < 0 {
		//	min, max = -max, -min
		//}
		//diff := math.Abs(min - max)
		//diff = diff / max

		//fmt.Printf("diff: %.6f\n", diff)
		if diff >= 0 {
			totalDeviation += diff
		} else {
			totalDeviation -= diff
		}
	}

	// проверяем среднее отклонение
	avgDeviation := totalDeviation / float64(l) * 100
	//fmt.Println("avgDeviation", avgDeviation)
	//fmt.Println("MaxDeviation", task.props.MaxDeviation)
	if avgDeviation <= task.props.MaxDeviation {
		return true
	}

	return false
}

// сигмоид
func sig(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// гиперболический тангенс
func gipTan(x float64) float64 {
	return (math.Exp(2*x) - 1) / (math.Exp(2*x) + 1)
}
