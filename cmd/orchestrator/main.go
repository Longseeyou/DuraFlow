package main

import (
	"context"
	"fmt"

	"github.com/Longseeyou/DuraFlow/internal/adapters/postgresql"
	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=gorm password=gorm dbname=duraflow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	ctx := context.Background()

	db.AutoMigrate(&workflow.WorkflowDefinition{})
	db.AutoMigrate(&workflow.WorkflowRun{})
	db.AutoMigrate(&workflow.Workflow{})
	db.AutoMigrate(&user.User{})

	// err = gorm.G[user.User](db).Create(ctx, &user.User{Email: "hi"})
	user, err := gorm.G[user.User](db).Where("Email like ?", "h%").First(ctx)

	workflowService := workflow.WorkflowService{Repository: postgresql.NewWorkflowRepositoryPostgres(db)}

	// workflowService.CreateWorkflow(ctx, user.ID, "hi", "hello")

	w := workflow.Workflow{}
	w.ID, _ = uuid.Parse("a0315a2c-5564-4f60-b9c7-94996dce4ae1")
	fmt.Println(workflowService.GetWorkflowByUserAndId(ctx, user.ID, w))

	// w := workflowService.GetWorkflowByUser(ctx, user.ID)
	// for i, v := range w {
	// 	fmt.Println("index:", i, ", value:", v)
	// }
}
