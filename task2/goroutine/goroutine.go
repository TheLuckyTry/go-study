package main

import (
	"fmt"
	"sync"
	"time"
)

// Task 任务结构
type Task struct {
	Name string
	Func func() error
}

// TaskResult 任务结果
type TaskResult struct {
	Name     string
	Duration time.Duration
	Error    error
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
	tasks []Task
}

// NewTaskScheduler 创建任务调度器
func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks: make([]Task, 0),
	}
}

func (ts *TaskScheduler) AddTask(name string, taskFunc func() error) {
	ts.tasks = append(ts.tasks, Task{Name: name, Func: taskFunc})
}

func PrintOdd(n int) {
	for i := 1; i <= n; i++ {
		if i%2 != 0 {
			fmt.Println(i)
		}
	}
}

func PrintEven(n int) {
	for i := 2; i <= n; i += 2 {
		fmt.Println(i)
	}
}

// 示例任务
func sampleTask1() error {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Executing task 1")
	return nil
}

func sampleTask2() error {
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Executing task 2")
	return nil
}

// Execute 执行所有任务
func (ts *TaskScheduler) Execute() {
	var wg sync.WaitGroup

	for _, task := range ts.tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()

			start := time.Now()
			err := t.Func()
			duration := time.Since(start)

			if err != nil {
				fmt.Printf("Task %s failed: %v\n", t.Name, err)
			} else {
				fmt.Printf("Task %s completed in %v\n", t.Name, duration)
			}
		}(task)
	}

	wg.Wait()
}
func main() {

	go func() {
		PrintOdd(10)
	}()
	go func() {
		PrintEven(10)
	}()
	time.Sleep(time.Second)

	scheduler := NewTaskScheduler()

	scheduler.AddTask("task1", sampleTask1)
	scheduler.AddTask("task2", sampleTask2)

	fmt.Println("Starting tasks...")
	scheduler.Execute()
}
