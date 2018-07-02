package subzero

import "testing"
import "fmt"
import "strings"
import "errors"

func TestAggregateSuccessfulResults(t *testing.T) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Success: 0},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	counter := 0

	for _ = range AggregateSuccessfulResults(fakeResultsChan) {
		counter++
	}

	if counter != 3 {
		t.Fatalf("expected '%v' successful results, got '%v'", 3, counter)
	}
}

func TestAggregateFailedResults(t *testing.T) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Success: 0},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	counter := 0

	for _ = range AggregateFailuedResults(fakeResultsChan) {
		counter++
	}

	if counter != 2 {
		t.Fatalf("expected '%v' failed results, got '%v'", 2, counter)
	}
}

func TestAggregateCustomResults(t *testing.T) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Success: false},
		&Result{Success: 0},
		&Result{Success: "wiggle"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	counter := 0

	for _ = range AggregateCustomResults(fakeResultsChan, func(r *Result) bool {
		_, ok := r.Success.(bool)
		return ok
	}) {
		counter++
	}

	if counter != 2 {
		t.Fatalf("expected '%v' successful results, got '%v'", 2, counter)
	}
}

func TestAggregateCustomResultsMore(t *testing.T) {
	fakeResults := []*Result{
		&Result{Success: true},
		&Result{Success: false},
		&Result{Success: 0},
		&Result{Success: "picat"},
		&Result{Success: "was"},
		&Result{Success: "here"},
		&Result{Failure: errors.New("example1")},
		&Result{Failure: errors.New("example2")},
	}

	fakeResultsChan := make(chan *Result)
	go func(fakeResults []*Result, fakeResultsChan chan *Result) {
		defer close(fakeResultsChan)
		for _, result := range fakeResults {
			fakeResultsChan <- result
		}
	}(fakeResults, fakeResultsChan)

	puzzle := []string{}

	for result := range AggregateCustomResults(fakeResultsChan, func(r *Result) bool {
		_, ok := r.Success.(string)
		return ok
	}) {
		puzzle = append(puzzle, result.Success.(string))
	}

	if strings.Join(puzzle, " ") != "picat was here" {
		t.Fatalf("expected '%v', got '%v'", "picat was here", strings.Join(puzzle, " "))
	}
}
