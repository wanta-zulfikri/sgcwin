package pkg

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"sync"
)

type Validation struct {
	Badwords map[string]struct{}
}

func (v *Validation) Validate(sentence string) bool {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
	}
	sentence = reg.ReplaceAllString(sentence, " ")
	words := strings.Split(sentence, " ")
	wordschan := make(chan string)
	wg := &sync.WaitGroup{}
	errchan := make(chan error, len(words))
	wg.Add(len(words) + 1)
	go func(wg *sync.WaitGroup, wordschan <-chan string) {
		wg2 := &sync.WaitGroup{}
		for i := 0; i < 3; i++ {
			wg2.Add(1)
			var err error
			go func(wg *sync.WaitGroup, wordschan <-chan string, wg2 *sync.WaitGroup, err error) {
				for word := range wordschan {
					if _, ok := v.Badwords[strings.ToLower(word)]; ok {
						errchan <- errors.New("badword")
						err = errors.New("badword")
					}
					wg.Done()
				}
				defer wg2.Done()
			}(wg, wordschan, wg2, err)
		}
		wg2.Wait()
		if err == nil {
			errchan <- nil
		}
		wg.Done()
	}(wg, wordschan)
	for _, word := range words {
		wordschan <- word
	}
	close(wordschan)
	wg.Wait()
	select {
	case err := <-errchan:
		if err != nil {
			close(errchan)
			return false
		}
	}
	close(errchan)
	return true
}
