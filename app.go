package main

import (
	"context"

	"log"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Stack struct {
	gorm.Model
	Name string `json:"name"`
	Order uint `json:"order"`
	Tasks []Task `json:"tasks"`
	Tags []Tag `json:"tags"`
}

type Tag struct {
	gorm.Model
	StackID uint `json:"stack_id"`
	Name string `json:"name"`
	Tasks []*Task `json:"tasks" gorm:"many2many:task_tags;"`
}

type Task struct {
	gorm.Model
	StackID uint `json:"stack_id"`
	Text string `json:"text"`
	ParsedText string `json:"parsed_text"`
	Tags []*Tag `json:"tags" gorm:"many2many:task_tags;"`
	Done bool `json:"done"`
	Order uint `json:"order"`
}

type TaskDelta struct {
	Text string `json:"text"`
	ParsedText string `json:"parsed_text"`
	Tags []string `json:"tags"`
}

// App struct
type App struct {
	ctx context.Context
	dbPath string
	db *gorm.DB
}

// NewApp creates a new App application struct
func NewApp(dbPath string) *App {
	return &App{dbPath: dbPath}
}

func (a *App) onSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
    runtime.WindowUnminimise(a.ctx)
    runtime.Show(a.ctx)
    go runtime.EventsEmit(a.ctx, "launchArgs", secondInstanceData.Args)
}

func initDB(db *gorm.DB) {
	db.AutoMigrate(&Stack{}, &Task{}, &Tag{})

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
	result = db.Clauses(clause.Returning{}).Where("stack_id = ?", s.ID).Delete(&Tag{})
	if result.Error != nil {
		log.Println("Failed to delete tags:", result.Error)
		return result.Error
	}
	return nil
}

func (a *App) cleanupTags(stackID uint) error {
	tags, err := a.GetTags(stackID)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		if len(tag.Tasks) == 0 {
			result := a.db.Delete(&tag)
			if result.Error != nil {
				log.Println("Failed to delete tag:", result.Error)
				return result.Error
			}
		}
	}
	return nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	db, err := gorm.Open(sqlite.Open(a.dbPath), &gorm.Config{})
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
	result := a.db.Preload("Tasks.Tags").Preload("Tags").First(&stack, id)
	if result.Error != nil {
		log.Println("Failed to get stack:", result.Error)
		return stack, result.Error
	}
	return stack, nil
}

func (a *App) GetTags(stackID uint) (tags []*Tag, err error) {
	result := a.db.Where("stack_id = ?", stackID).Preload("Tasks").Find(&tags)
	if result.Error != nil {
		log.Println("Failed to get tags:", result.Error)
		return tags, result.Error
	}
	return tags, nil
}

func (a *App) makeTags(stackID uint, tagNames []string) ([]*Tag, error) {
	var tags []*Tag
	var existing *Tag

	for _, name := range tagNames {
		tag := &Tag{StackID: stackID, Name: name}
		result := a.db.Model(&Tag{}).Where("name = ?", name).Limit(1).Find(&existing)
		if existing != nil && result.Error == nil {
			tag.ID = existing.ID
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (a *App) AddTask(stackID uint, delta TaskDelta, order uint) (Task, error) {
	var task Task
	tags, err := a.makeTags(stackID, delta.Tags)
	if err != nil {
		return task, err
	}
	task = Task{StackID: stackID, Text: delta.Text, ParsedText: delta.ParsedText, Tags: tags, Order: order}
	result := a.db.Create(&task)
	if result.Error != nil {
		log.Println("Failed to create task:", result.Error)
		return task, result.Error
	}
	return task, nil
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

func (a *App) EditTask(taskID uint, delta TaskDelta) (Task, error) {
	var task Task
	result := a.db.Model(&Task{}).Where("ID = ?", taskID).First(&task)
	if result.Error != nil {
		log.Println("Failed to get task:", result.Error)
		return task, result.Error
	}

	tags, err := a.makeTags(task.StackID, delta.Tags)
	if err != nil {
		return task, err
	}
	task.Text = delta.Text
	task.ParsedText = delta.ParsedText

	result = a.db.Model(&task).Updates(&task)
	if result.Error != nil {
		log.Println("Failed to save task:", result.Error)
		return task, result.Error
	}

	err = a.db.Model(&task).Association("Tags").Replace(tags)
	if err != nil {
		return task, err
	}
	err = a.cleanupTags(task.StackID)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (a *App) DeleteTask(taskID uint) (uint, error) {
	var task Task
	result := a.db.Model(&Task{}).Where("ID = ?", taskID).First(&task)
	if result.Error != nil {
		log.Println("Failed to get task to delete:", result.Error)
		return 0, result.Error
	}
	result = a.db.Delete(&task)
	if result.Error != nil {
		log.Println("Failed to delete task:", result.Error)
		return 0, result.Error
	}
	err := a.cleanupTags(task.StackID)
	if err != nil {
		return 0, err
	}
	return taskID, nil
}
