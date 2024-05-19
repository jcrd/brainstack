package main

import (
	"context"

	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Stack struct {
	gorm.Model
	Name string `json:"name"`
	Order uint `json:"order"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	gorm.Model
	StackID uint `json:"stack_id"`
	Text string `json:"text"`
	Done bool `json:"done"`
	Order uint `json:"order"`
}

// App struct
type App struct {
	ctx context.Context
	db *gorm.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func initDB(db *gorm.DB) {
	db.AutoMigrate(&Stack{}, &Task{})

	var stacks []Stack
	result := db.Find(&stacks)
	if result.Error != nil {
		panic("Failed to get stacks after DB migration")
	}
	if len(stacks) == 0 {
		db.Create(&Stack{Name: "inbox"})
	}
}

func (s *Stack) BeforeDelete(db *gorm.DB) error {
	result := db.Clauses(clause.Returning{}).Where("stack_id = ?", s.ID).Delete(&Task{})
	if result.Error != nil {
		log.Println("Failed to delete tasks:", result.Error)
		return result.Error
	}
	return nil
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	initDB(db)
	a.db = db
}

func (a *App) CreateStack(name string) (uint, error) {
	stack := Stack{Name: name}
	result := a.db.Create(&stack)
	if result.Error != nil {
		log.Println("Failed to create stack:", result.Error)
		return 0, result.Error
	}
	return stack.ID, nil
}

func (a *App) DeleteStack(stackID uint) (uint, error) {
	stack := &Stack{}
	stack.ID = stackID
	result := a.db.Delete(stack)
	if result.Error != nil {
		log.Println("Failed to delete stack:", result.Error)
		return 0, result.Error
	}
	return stackID, nil
}

func (a *App) EditStack(stackID uint, name string) (uint, error) {
	result := a.db.Model(&Stack{}).Where("ID = ?", stackID).Update("name", name)
	if result.Error != nil {
		log.Println("Failed to update stack:", result.Error)
		return 0, result.Error
	}
	return stackID, nil
}

func (a *App) GetStacks() ([]Stack, error) {
	var stacks []Stack
	result := a.db.Find(&stacks)
	if result.Error != nil {
		log.Println("Failed to get stacks:", result.Error)
		return stacks, result.Error
	}
	return stacks, nil
}

func (a *App) GetStack(id uint) (Stack, error) {
	var stack Stack
	result := a.db.Preload("Tasks").First(&stack, id)
	if result.Error != nil {
		log.Println("Failed to get stack:", result.Error)
		return stack, result.Error
	}
	return stack, nil
}

func (a *App) AddTask(stackID uint, text string, order uint) (uint, error) {
	task := Task{StackID: stackID, Text: text, Order: order}
	result := a.db.Create(&task)
	if result.Error != nil {
		log.Println("Failed to create task:", result.Error)
		return 0, result.Error
	}
	return task.ID, nil
}

func (a *App) ReorderTasks(tasks []Task) ([]Task, error) {
	for _, t := range tasks {
		result := a.db.Model(&t).Update("order", t.Order)
		if result.Error != nil {
			log.Println("Failed to reorder tasks:", result.Error)
			return tasks, result.Error
		}
	}
	return tasks, nil
}

func (a *App) UpdateTaskDone(taskID uint, state bool) (uint, error) {
	result := a.db.Model(&Task{}).Where("ID = ?", taskID).Update("done", state)
	if result.Error != nil {
		log.Println("Failed to update task:", result.Error)
		return 0, result.Error
	}
	return taskID, nil
}

func (a *App) EditTask(taskID uint, text string) (uint, error) {
	result := a.db.Model(&Task{}).Where("ID = ?", taskID).Update("text", text)
	if result.Error != nil {
		log.Println("Failed to update task:", result.Error)
		return 0, result.Error
	}
	return taskID, nil
}

func (a *App) DeleteTask(taskID uint) (uint, error) {
	result := a.db.Delete(&Task{}, taskID)
	if result.Error != nil {
		log.Println("Failed to delete task:", result.Error)
		return 0, result.Error
	}
	return taskID, nil
}
