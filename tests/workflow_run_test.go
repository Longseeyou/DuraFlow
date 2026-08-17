package tests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/Longseeyou/DuraFlow/internal/adapters/postgresql"
	"github.com/Longseeyou/DuraFlow/internal/orchestrator"
	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	testUserEmail = "workflow-run-e2e@duraflow.dev"

	defaultDSN     = "host=localhost user=gorm password=gorm dbname=duraflow"
	defaultBrokers = "localhost:9190,localhost:9191,localhost:9192"

	totalTaskRuns = 1600
	runTimeout    = 300 * time.Second
	pollInterval  = time.Second
	expectedDeps  = 1500000
)

// TestWorkflowRunEndToEnd creates a user, a workflow named "test" with 5 layers
// of mock tasks ([1, 2, 4, 2, 1]) fully connected between consecutive layers,
// and runs the workflow.
//
// Prerequisites:
//   - Postgres and the Kafka cluster must be running
//   - The orchestrator (go run ./cmd/orchestrator) and worker (go run ./cmd/worker)
//     binaries must be running
func TestWorkflowRunEndToEnd(t *testing.T) {
	ctx := context.Background()
	testStarted := time.Now()

	t.Log("stage: checking prerequisites (kafka + postgres)")
	requireKafka(t, kafkaBrokers())
	db := openDB(t)

	t.Log("stage: migrating schema and cleaning up previous test data")
	if err := migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	cleanupTestData(db)

	t.Log("stage: creating user")
	userService := user.NewUserService(postgresql.NewPostgresUserRepository(db))
	if _, err := userService.CreateUser(ctx, user.CreateUserRequestDto{
		Name:            "Test User",
		Email:           testUserEmail,
		Password:        "password123",
		RetypedPassword: "password123",
		Role:            user.USER,
	}); err != nil {
		t.Fatalf("create user: %v", err)
	}

	userRepo := postgresql.NewPostgresUserRepository(db)
	u, err := userRepo.GetUserByEmail(ctx, testUserEmail)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}

	taskService := task.NewTaskService(postgresql.NewPostgresTaskRepository(db))
	workflowService := workflow.NewWorkflowService(
		postgresql.NewPostgresWorkflowRepository(db),
		postgresql.NewPostgresWorkflowRepositoryInternal(db),
		&taskService,
	)

	t.Log("stage: creating workflow 'test'")
	w, err := workflowService.CreateWorkflow(ctx, u.ID, workflow.WorkflowRequestDto{
		Name:        new("test"),
		Description: new("workflow run end-to-end test"),
	})
	if err != nil {
		t.Fatalf("create workflow: %v", err)
	}

	t.Log("stage: creating workflow definition")
	wD, err := workflowService.CreateWorkflowDefinition(ctx, u.ID, w.ID)
	if err != nil {
		t.Fatalf("create workflow definition: %v", err)
	}

	layerSizes := []int{
		100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100,
	}
	taskType := task.MOCK_TASK
	timeout := 30 * time.Second

	t.Log("stage: creating task definitions")
	layerTaskIDs := make([][]uuid.UUID, len(layerSizes))
	for layer, size := range layerSizes {
		t.Logf("  layer %d/%d: creating %d mock tasks", layer+1, len(layerSizes), size)
		for i := range size {
			name := fmt.Sprintf("layer-%d-task-%d", layer+1, i+1)
			tD, err := taskService.CreateTaskDefinition(
				ctx,
				u.ID,
				wD.ID,
				task.TaskDefinitionRequestDto{
					Name:     &name,
					TaskType: &taskType,
					Timeout:  &timeout,
				},
			)
			if err != nil {
				t.Fatalf("create task definition %s: %v", name, err)
			}
			layerTaskIDs[layer] = append(layerTaskIDs[layer], tD.ID)
		}
	}

	t.Log("stage: creating fully-connected dependencies between consecutive layers")
	dependencyCount := 0
	for layer := 0; layer < len(layerSizes)-1; layer++ {
		for _, nextTaskID := range layerTaskIDs[layer+1] {
			for _, prevTaskID := range layerTaskIDs[layer] {
				if _, err := taskService.CreateTaskDependency(
					ctx,
					u.ID,
					task.TaskDependencyRequestDto{
						TaskID:         nextTaskID,
						DependOnTaskID: prevTaskID,
					},
				); err != nil {
					t.Fatalf(
						"create task dependency %s -> %s: %v",
						prevTaskID,
						nextTaskID,
						err,
					)
				}
				dependencyCount++
			}
		}
	}
	t.Logf("  created %d dependencies", dependencyCount)

	t.Log("stage: activating workflow definition (DAG check)")
	status := workflow.WORKFLOW_DEFINITION_ACTIVATED
	if _, err := workflowService.UpdateWorkflowDefinitionByUserAndID(
		ctx,
		u.ID,
		wD.ID,
		workflow.WorkflowDefinitionRequestDto{Status: &status},
	); err != nil {
		t.Fatalf("activate workflow definition: %v", err)
	}

	t.Log("stage: triggering workflow run")
	orch := orchestrator.NewOrchestrator(postgresql.NewPostgresOrchestratorRepository(db))
	if err := orch.RunWorkflow(ctx, u.ID, wD.ID); err != nil {
		t.Fatalf("run workflow: %v", err)
	}

	var wR workflow.WorkflowRun
	if err := db.Where("workflow_definition_id = ?", wD.ID).
		Order("created_at DESC").
		First(&wR).Error; err != nil {
		t.Fatalf("get workflow run: %v", err)
	}
	t.Logf("  workflow run %s started (%d task runs)", wR.ID, totalTaskRuns)

	t.Log("stage: waiting for all task runs to complete")
	var runs []task.TaskRun
	deadline := time.Now().Add(runTimeout)
	for time.Now().Before(deadline) {
		runs = nil
		if err := db.Where("workflow_run_id = ?", wR.ID).Find(&runs).Error; err != nil {
			t.Fatalf("get task runs: %v", err)
		}

		completed := 0
		for _, r := range runs {
			if r.Status == task.TASK_RUN_COMPLETED {
				completed++
			}
		}

		t.Logf("  progress: %d/%d task runs completed", completed, totalTaskRuns)
		if len(runs) == totalTaskRuns && completed == totalTaskRuns {
			break
		}

		time.Sleep(pollInterval)
	}

	statusCount := map[task.TaskRunStatus]int{}
	for _, r := range runs {
		statusCount[r.Status]++
	}
	if len(runs) != totalTaskRuns {
		t.Fatalf("expected %d task runs, got %d (%v)", totalTaskRuns, len(runs), statusCount)
	}
	if statusCount[task.TASK_RUN_COMPLETED] != totalTaskRuns {
		t.Fatalf(
			"workflow run %s did not complete in %s, task run statuses: %v "+
				"(is the orchestrator and worker running?)",
			wR.ID,
			runTimeout,
			statusCount,
		)
	}
	t.Logf("stage: all %d task runs completed in %s", totalTaskRuns, time.Since(testStarted))

	t.Log(
		"stage: verifying dependency ordering (dependents started after their dependencies ended)",
	)
	var dependencies []task.TaskDependency
	if err := db.Joins("JOIN task_definitions td_task ON td_task.id = task_dependencies.task_id").
		Where("td_task.workflow_definition_id = ?", wD.ID).
		Find(&dependencies).Error; err != nil {
		t.Fatalf("get dependencies: %v", err)
	}
	if len(dependencies) != expectedDeps {
		t.Errorf("expected %d dependencies, got %d", expectedDeps, len(dependencies))
	}

	for _, dep := range dependencies {
		var dependent, dependedOn task.TaskRun
		if err := db.Where("workflow_run_id = ? AND task_definition_id = ?", wR.ID, dep.TaskID).
			First(&dependent).Error; err != nil {
			t.Fatalf("get dependent run: %v", err)
		}
		if err := db.Where("workflow_run_id = ? AND task_definition_id = ?", wR.ID, dep.DependOnTaskID).
			First(&dependedOn).
			Error; err != nil {
			t.Fatalf("get dependency run: %v", err)
		}

		var dependentAttempt task.TaskAttempt
		if err := db.Where("task_run_id = ?", dependent.ID).
			Order("attempt_number ASC").
			First(&dependentAttempt).Error; err != nil {
			t.Fatalf("get dependent attempt: %v", err)
		}

		if dependentAttempt.StartedAt.Before(*dependedOn.EndedAt) {
			t.Errorf(
				"task run %s (first attempt started %s) started before dependency task run %s ended (%s)",
				dependent.ID,
				dependentAttempt.StartedAt,
				dependedOn.ID,
				dependedOn.EndedAt,
			)
		}
	}

	t.Logf(
		"stage: done — workflow run %s completed (%d task runs, %d dependencies verified, %s elapsed)",
		wR.ID,
		totalTaskRuns,
		expectedDeps,
		time.Since(testStarted),
	)
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&user.User{},
		&workflow.Workflow{},
		&workflow.WorkflowDefinition{},
		&workflow.WorkflowRun{},
		&task.TaskDefinition{},
		&task.TaskDependency{},
		&task.TaskRun{},
		&task.TaskAttempt{},
		// &task.TaskEvent{},
	)
}

func cleanupTestData(db *gorm.DB) {
	var u user.User
	if err := db.Unscoped().Where("email = ?", testUserEmail).First(&u).Error; err != nil {
		return
	}

	var workflows []workflow.Workflow
	db.Unscoped().Where("user_id = ?", u.ID).Find(&workflows)
	for _, w := range workflows {
		var wDs []workflow.WorkflowDefinition
		db.Unscoped().Where("workflow_id = ?", w.ID).Find(&wDs)
		for _, wD := range wDs {
			var tDs []task.TaskDefinition
			db.Unscoped().Where("workflow_definition_id = ?", wD.ID).Find(&tDs)

			tDIDs := make([]uuid.UUID, 0, len(tDs))
			for _, tD := range tDs {
				tDIDs = append(tDIDs, tD.ID)
			}

			if len(tDIDs) > 0 {
				var runs []task.TaskRun
				db.Unscoped().Where("task_definition_id IN ?", tDIDs).Find(&runs)

				runIDs := make([]uuid.UUID, 0, len(runs))
				for _, r := range runs {
					runIDs = append(runIDs, r.ID)
				}

				if len(runIDs) > 0 {
					db.Unscoped().Where("task_run_id IN ?", runIDs).Delete(&task.TaskAttempt{})
				}
				db.Unscoped().Where("task_definition_id IN ?", tDIDs).Delete(&task.TaskRun{})
				db.Unscoped().
					Where("task_id IN ? OR depend_on_task_id IN ?", tDIDs, tDIDs).
					Delete(&task.TaskDependency{})
				db.Unscoped().Where("id IN ?", tDIDs).Delete(&task.TaskDefinition{})
			}

			db.Unscoped().Where("workflow_definition_id = ?", wD.ID).Delete(&workflow.WorkflowRun{})
			db.Unscoped().Where("id = ?", wD.ID).Delete(&workflow.WorkflowDefinition{})
		}
		db.Unscoped().Where("id = ?", w.ID).Delete(&workflow.Workflow{})
	}

	db.Unscoped().Where("id = ?", u.ID).Delete(&user.User{})
}

func openDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = defaultDSN
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf(
			"connect database %q: %v (is Postgres running?)",
			dsn,
			err,
		)
	}
	return db
}

func requireKafka(t *testing.T, brokers []string) {
	t.Helper()

	config := sarama.NewConfig()
	config.Net.DialTimeout = 3 * time.Second

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		t.Fatalf(
			"kafka not reachable at %v: %v (run: docker compose -f docker-compose-kafka-cluster.yaml up -d)",
			brokers,
			err,
		)
	}
	defer client.Close()
}

func kafkaBrokers() []string {
	raw := os.Getenv("KAFKA_BROKERS")
	if raw == "" {
		raw = defaultBrokers
	}
	return strings.Split(raw, ",")
}
