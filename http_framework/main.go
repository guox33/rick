package main

import (
	"code.byted.org/gopkg/retry"
	"errors"
	"fmt"
	"time"
)

type fileRespWrite struct {
}

/*func (f *fileRespWrite) Write() {

}

func (f *fileRespWrite) Header() http.Header {

}

func (f *fileRespWrite) WriteHeader(statusCode int) {

}*/

func main() {
	err := retry.Do("what", 2, time.Second, func() error {
		fmt.Println("hello")
		// time.Sleep(time.Minute)
		return errors.New("world")
	})
	fmt.Println(err.Error())
}
