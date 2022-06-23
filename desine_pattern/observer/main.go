package observer

import (
	"fmt"
	"github.com/thoas/go-funk"
	"sync"
)

type IObserver interface {
	Update(msg string)
}

type ISubject interface {
	Subscribe(observer IObserver)
	Remove(observer IObserver)
	Notify(msg string)
}

type Media struct {
	users []IObserver

	locker sync.RWMutex
}

func (m *Media) Subscribe(user IObserver) {
	m.locker.Lock()
	defer m.locker.Unlock()

	m.users = append(m.users, user)
}

func (m *Media) Remove(user IObserver) {
	m.locker.Lock()
	defer m.locker.Unlock()

	idx := funk.IndexOf(m.users, user)
	if idx == -1 {
		return
	}

	m.users = append(m.users[:idx], m.users[idx+1:]...)
}

func (m *Media) Notify(msg string) {
	for _, user := range m.users {
		user.Update(msg)
	}
}

type User struct {
}

func (u *User) Update(msg string) {
	fmt.Println("I got news: ", msg)
}

func Main() {
	var media ISubject = &Media{}

	var userA IObserver = &User{}
	var userB IObserver = &User{}

	media.Subscribe(userA)
	media.Subscribe(userB)

	media.Notify("this is a public news: hello world")

	media.Remove(userB)

	media.Notify("this is a private news: whatever")
}
