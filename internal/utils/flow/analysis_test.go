package flow

import (
	"fmt"
	"keyboard-analysis/internal/models"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestAnalyseEvents(t *testing.T) {
	type args struct {
		flow *models.KeyboardFlow
	}
	tests := []struct {
		name    string
		args    args
		want    *KeyboardAnalysis
		wantErr bool
	}{
		{
			name:    "Test 1 [Straight - abcd]",
			args:    args{flow: &models.KeyboardFlow{Flow: testEventsStraight("abcd"), Password: models.Password{Password: "abcd"}}},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Test 2 [Concurrent - abcd]",
			args:    args{flow: &models.KeyboardFlow{Flow: testEventsConcurrent("abcd"), Password: models.Password{Password: "abcd"}}},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, x := range tt.args.flow.Flow {
				fmt.Println(x.Key, x.Up, x.Time)
			}
			got, err := AnalyseEvents(tt.args.flow)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventsToVector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventsToVector() got = %v, want %v", got, tt.want)
				fmt.Println("Explanation: ")
				for _, x := range got.Inline() {
					fmt.Println(x)
				}
			}
		})
	}
}

func testEventsStraight(phrase string) []*models.KeyboardEvent {
	var res []*models.KeyboardEvent
	for i, key := range []rune(phrase) {
		res = append(res, &models.KeyboardEvent{
			Key:  string(key),
			Up:   false,
			Time: uint(i * 100),
		}, &models.KeyboardEvent{
			Key:  string(key),
			Up:   true,
			Time: uint(i*100) + 1,
		})
	}
	return res
}

func testEventsConcurrent(phrase string) []*models.KeyboardEvent {
	var res []*models.KeyboardEvent
	start := time.Now()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i, key := range []rune(phrase) {
		wg.Add(2)

		k := time.Duration(i * 1000)
		go func(key rune) {
			time.Sleep(k * time.Millisecond)
			mu.Lock()
			res = append(res, &models.KeyboardEvent{
				Key:  string(key),
				Up:   false,
				Time: uint(time.Since(start)),
			})
			mu.Unlock()
			wg.Done()
		}(key)

		go func(key rune) {
			k2 := k + time.Duration(rand.Intn(2000))
			time.Sleep(k2 * time.Millisecond)
			mu.Lock()
			res = append(res, &models.KeyboardEvent{
				Key:  string(key),
				Up:   true,
				Time: uint(time.Since(start)),
			})
			mu.Unlock()
			wg.Done()
		}(key)
	}
	wg.Wait()
	return res
}
