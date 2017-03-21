package main

import "fmt"
import "time"

func main() {
	tasks := []*tpkTask{
		&tpkTask{drName: "DR_ALSACE"},
		&tpkTask{drName: "DR_NORD"},
		&tpkTask{drName: "DR_RHONE"},
		&tpkTask{drName: "DR_PACA"},
		&tpkTask{drName: "DR_IDF"},
		&tpkTask{drName: "DR_BRETAGNE"},
	}

	dirInput := make(chan *tpkTask)
	genInput := make(chan *tpkTask)
	dlInput := make(chan *tpkTask)
	endInput := make(chan *tpkTask)

	go dirCreator(dirInput, genInput)
	go tpkGenerator(genInput, dlInput)
	go tpkDownloador(dlInput, endInput)
	go end(endInput)

	for {
		for _, task := range tasks {
			dirInput <- task
		}
	}
}

func dirCreator(input chan *tpkTask, next chan *tpkTask) {
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
			}
		}
	}
}

func tpkGenerator(input chan *tpkTask, next chan *tpkTask) {
	for task := range input {
		if task.tpkGenerated {
			fmt.Printf("tpkGenerator: %s: dir already generated\n", task.drName)
			next <- task
		} else {
			task.tpkGenerated = generateTpk(task.drName)
			if task.tpkGenerated {
				fmt.Printf("tpkGenerator: %s: tpk created successfully\n", task.drName)
				next <- task
			} else {
				fmt.Printf("tpkGenerator: %s: tpk not created :(\n", task.drName)
			}
		}
	}
}

func tpkDownloador(input chan *tpkTask, next chan *tpkTask) {
	for task := range input {
		if task.tpkDownloaded {
			fmt.Printf("tpkDownloador: %s: dir already downloaded\n", task.drName)
			next <- task
		} else {
			task.tpkDownloaded = downloadTpk(task.drName)
			if task.tpkDownloaded {
				fmt.Printf("tpkDownloador: %s: tpk downloaded successfully\n", task.drName)
				next <- task
			} else {
				fmt.Printf("tpkDownloador: %s: tpk not downloaded :(\n", task.drName)
			}
		}
	}
}

func end(input chan *tpkTask) {
	for task := range input {
		fmt.Printf("end: %s: finished: dir: %t, gen: %t, dl: %t\n", task.drName,
			task.dirCreated, task.tpkGenerated, task.tpkDownloaded)
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

func generateTpk(name string) bool {
	fmt.Printf("generateTpk: %s\n", name)

	time.Sleep(5 * time.Second)

	fmt.Printf("generateTpk: %s: done\n", name)
	return true
}

func downloadTpk(name string) bool {
	fmt.Printf("downloadTpk: %s\n", name)

	time.Sleep(5 * time.Second)

	fmt.Printf("downloadTpk: %s: done\n", name)
	return true
}
