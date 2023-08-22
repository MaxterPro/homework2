package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	numWorkers := 5

	jobs := make(chan int, 100)
	result := make(chan User, 100)

	for i := 0; i < numWorkers; i++ {
		go workerUser(jobs, result, i)
	}

	for i := 0; i < 100; i++ {
		jobs <- i
		fmt.Printf("generated user %d\n", i+1)
	}
	close(jobs)

	for i := 0; i < 100; i++ {
		saveUserInfo(result)
	}

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(result chan User) {
	r := <-result
	fmt.Printf("WRITING FILE FOR UID %d\n", r.id)

	filename := fmt.Sprintf("users/uid%d.txt", r.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(r.getActivityInfo())
	time.Sleep(time.Second)
}

func workerUser(jobs chan int, result chan User, id int) {
	for i := range jobs {
		result <- User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@company.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
	}
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
