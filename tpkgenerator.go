package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	tasks := []*tpkTask{
		&tpkTask{drName: "DR_ALSACE"},
		&tpkTask{drName: "DR_NORD"},
		&tpkTask{drName: "DR_RHONE"},
		&tpkTask{drName: "DR_PACA"},
		&tpkTask{drName: "DR_IDF"},
		&tpkTask{drName: "DR_BRETAGNE"},
	}

	initInput := make(chan *tpkTask, len(tasks))
	dirInput := make(chan *tpkTask)
	genInput := make(chan *tpkTask)
	dlInput := make(chan *tpkTask)
	endInput := make(chan *tpkTask)
	wg := sync.WaitGroup{}

	go initiator(initInput, dirInput)
	go dirCreator(dirInput, genInput, initInput)
	go tpkGenerator(genInput, dlInput, initInput)
	go tpkDownloador(dlInput, endInput, initInput)
	go end(endInput, &wg)

	start(tasks, initInput, &wg)
	wg.Wait()

	fmt.Printf("STATUS: ")
	for _, task := range tasks {
		fmt.Printf("%s: dir=%t, tpk=%t, dld:%t\n", task.drName, task.dirCreated, task.tpkGenerated, task.tpkDownloaded)
	}
}

func start(tasks []*tpkTask, next chan *tpkTask, wg *sync.WaitGroup) {
	for _, task := range tasks {
		if !task.tpkDownloaded {
			wg.Add(1)
			next <- task
		}
	}
}

func initiator(input chan *tpkTask, next chan *tpkTask) {
	for task := range input {
		next <- task
	}
}

func dirCreator(input chan *tpkTask, next chan *tpkTask, fail chan *tpkTask) {
	for task := range input {
		if task.dirCreated {
			fmt.Printf("dirCreator: %s: dir already created\n", task.drName)
			next <- task
		} else {
			task.dirCreated = createDir(task.drName)
			if task.dirCreated {
				fmt.Printf("dirCreator: %s: dir created successfully\n", task.drName)
				next <- task
			} else {
				fmt.Printf("dirCreator: %s: dir not created :(\n", task.drName)
				fail <- task
			}
		}
	}
}

func tpkGenerator(input chan *tpkTask, next chan *tpkTask, fail chan *tpkTask) {
	generator := rand.New(rand.NewSource(42))
	for task := range input {
		if task.tpkGenerated {
			fmt.Printf("tpkGenerator: %s: dir already generated\n", task.drName)
			next <- task
		} else {
			task.tpkGenerated = generateTpk(task.drName, generator)
			if task.tpkGenerated {
				fmt.Printf("tpkGenerator: %s: tpk created successfully\n", task.drName)
				next <- task
			} else {
				fmt.Printf("tpkGenerator: %s: tpk not created :(\n", task.drName)
				fail <- task
			}
		}
	}
}

func tpkDownloador(input chan *tpkTask, next chan *tpkTask, fail chan *tpkTask) {
	generator := rand.New(rand.NewSource(51))
	for task := range input {
		if task.tpkDownloaded {
			fmt.Printf("tpkDownloador: %s: dir already downloaded\n", task.drName)
			next <- task
		} else {
			task.tpkDownloaded = downloadTpk(task.drName, generator)
			if task.tpkDownloaded {
				fmt.Printf("tpkDownloador: %s: tpk downloaded successfully\n", task.drName)
				next <- task
			} else {
				fmt.Printf("tpkDownloador: %s: tpk not downloaded :(\n", task.drName)
				fail <- task
			}
		}
	}
}

func end(input chan *tpkTask, wg *sync.WaitGroup) {
	for task := range input {
		fmt.Printf("end: %s: finished: dir: %t, gen: %t, dl: %t\n", task.drName,
			task.dirCreated, task.tpkGenerated, task.tpkDownloaded)
		wg.Done()
	}
}

type tpkTask struct {
	drName        string
	dirCreated    bool
	tpkGenerated  bool
	tpkDownloaded bool
}

func createDir(dirName string) bool {
	fmt.Printf("createDir: %s\n", dirName)

	time.Sleep(1 * time.Second)

	fmt.Printf("createDir: %s: done\n", dirName)
	return true
}

func generateTpk(name string, generator *rand.Rand) bool {
	fmt.Printf("generateTpk: %s\n", name)

	time.Sleep(1 * time.Second)

	val := generator.Intn(100)
	done := val > 80
	fmt.Printf("generateTpk: %s: done: %t (%d)\n", name, done, val)
	return done
}

func downloadTpk(name string, generator *rand.Rand) bool {
	fmt.Printf("downloadTpk: %s\n", name)

	time.Sleep(5 * time.Second)

	fmt.Printf("downloadTpk: %s: done\n", name)
	val := generator.Intn(100)
	done := val > 80
	fmt.Printf("downloadTpk: %s: done: %t (%d)\n", name, done, val)
	return done
}
