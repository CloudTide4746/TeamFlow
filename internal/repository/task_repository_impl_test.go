package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"teamflow/config"
	"teamflow/internal/dto"
	"teamflow/storage"

	"github.com/joho/godotenv"
)

// TestMain runs once before the tests in this package.
func TestMain(m *testing.M) {
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		fmt.Fprintln(os.Stderr, "resolve project root:", err)
		os.Exit(1)
	}

	if err := godotenv.Load(filepath.Join(projectRoot, ".env")); err != nil && !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "load test environment:", err)
		os.Exit(1)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = filepath.Join(projectRoot, "config", "config.yaml")
	} else if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(projectRoot, configPath)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "加载测试配置失败:", err)
		os.Exit(1)
	}
	storage.InitDB(cfg.Database.GetDSN())

	exitCode := m.Run()
	if storage.DB != nil {
		if sqlDB, err := storage.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
	os.Exit(exitCode)
}

func Test_taskRepository_ListTasks(t *testing.T) {
	if storage.DB == nil {
		t.Fatal("test database was not initialized")
	}

	r := NewTaskRepository(storage.DB).(*taskRepository)
	pageSize := 10
	tasks, total, err := r.ListTasks(dto.TaskQuery{
		SortBy:   "created_at",
		SortDir:  "desc",
		Page:     1,
		PageSize: pageSize,
	})
	if err != nil {
		t.Fatalf("ListTasks() returned an error: %v", err)
	}
	if len(tasks) > pageSize {
		t.Errorf("ListTasks() returned %d tasks; page size is %d", len(tasks), pageSize)
	}
	if total < int64(len(tasks)) {
		t.Errorf("ListTasks() total = %d, smaller than returned task count %d", total, len(tasks))
	}
}
