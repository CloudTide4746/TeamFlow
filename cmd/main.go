package main

import (
	"database/sql"
	"fmt"
	"os"
	"teamflow/config"
	"teamflow/internal/controller"
	"teamflow/internal/database"
	"teamflow/internal/repository"
	"teamflow/internal/router"
	"teamflow/internal/service"
	"teamflow/pkg/jwt"
	"teamflow/pkg/logger"
	"teamflow/storage"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// 0. 加载环境变量（可选）
	if err := godotenv.Load(); err != nil {
		fmt.Println("未找到 .env 文件，使用配置文件中的默认值")
	}

	// 1. 加载配置
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}
	if err := jwt.Configure(cfg.JWT.Secret); err != nil {
		panic(fmt.Sprintf("初始化 JWT 失败: %v", err))
	}

	// 2. 初始化日志
	if err := logger.InitLogger(&cfg.Log); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer logger.Sync()

	logger.Info("应用启动",
		zap.String("app_name", cfg.App.Name),
		zap.String("env", cfg.App.Env),
	)

	fmt.Println("✅ 配置和日志系统初始化成功！")

	// 3. 初始化数据库
	storage.InitDB(cfg.Database.GetDSN())

	if err := database.AutoMigrate(); err != nil {
		panic(fmt.Sprintf("数据库迁移失败: %v", err))
	}
	userService := service.NewUserService()
	userController := controller.NewUserController(userService)
	teamService := service.NewTeamService()
	teamController := controller.NewTeamController(teamService)
	projectService := service.NewProjectService()
	projectController := controller.NewProjectController(projectService)
	taskRepo := repository.NewTaskRepository(storage.DB)
	taskService := service.NewTaskService(taskRepo, &service.NoopNotificationService{})
	taskController := controller.NewTaskController(taskService)
	r := router.SetupRouter(userController, teamController, projectController, taskController)

	if err := r.Run(":8080"); err != nil {
		panic(fmt.Sprintf("启动服务器失败: %v", err))
	}
	fmt.Println("✅ 服务器启动成功，监听端口: 8080")
}

// TestAllComponents 测试所有已完成的组件
// 运行方式: go run cmd/main.go -test
// ⚠️ 注意: 此测试函数会真正连接数据库，请确保数据库配置正确
func TestAllComponents() {
	fmt.Println("========== 开始测试所有组件 ==========")

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ 未找到 .env 文件，使用配置文件中的默认值")
	}

	// 1. 测试配置加载
	fmt.Println("\n📋 测试 1: 配置加载")
	testConfig, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("   ❌ 配置加载失败: %v\n", err)
		return
	}
	fmt.Printf("   ✅ 配置加载成功\n")
	fmt.Printf("      - 应用名称: %s\n", testConfig.App.Name)
	fmt.Printf("      - 环境: %s\n", testConfig.App.Env)
	fmt.Printf("      - 数据库: %s:%d/%s\n", testConfig.Database.Host, testConfig.Database.Port, testConfig.Database.DBName)

	// 2. 测试日志初始化
	fmt.Println("\n📋 测试 2: 日志初始化")
	if err := logger.InitLogger(&testConfig.Log); err != nil {
		fmt.Printf("   ❌ 日志初始化失败: %v\n", err)
		return
	}
	logger.Info("测试日志 - Info 级别")
	logger.Debug("测试日志 - Debug 级别")
	logger.Warn("测试日志 - Warn 级别")
	logger.Error("测试日志 - Error 级别")
	fmt.Println("   ✅ 日志初始化成功")

	// 3. 测试数据库连接
	fmt.Println("\n📋 测试 3: 数据库连接")
	storage.InitDB(testConfig.Database.GetDSN())
	if storage.DB != nil {
		sqlDB, err := storage.DB.DB()
		if err != nil {
			fmt.Printf("   ❌ 获取数据库实例失败: %v\n", err)
			return
		}
		if err := sqlDB.Ping(); err != nil {
			fmt.Printf("   ❌ 数据库 Ping 失败: %v\n", err)
			return
		}
		fmt.Printf("   ✅ 数据库连接成功\n")
		stats := sqlDB.Stats()
		fmt.Printf("      - 最大打开连接数: %d\n", stats.MaxOpenConnections)
		fmt.Printf("      - 打开连接数: %d\n", stats.OpenConnections)
	} else {
		fmt.Println("   ❌ 数据库实例为空")
		return
	}

	// 4. 测试数据库自动迁移
	//fmt.Println("\n📋 测试 4: 数据库自动迁移")
	//if err := database.AutoMigrate(); err != nil {
	//	fmt.Printf("   ❌ 自动迁移失败: %v\n", err)
	//	return
	//}
	//fmt.Println("   ✅ 自动迁移成功")

	// 5. 测试配置的各项 Getter 方法
	fmt.Println("\n📋 测试 5: 配置 Getter 方法")
	fmt.Printf("   ✅ DSN: %s\n", testConfig.Database.GetDSN())
	fmt.Printf("   ✅ Redis Addr: %s\n", testConfig.Redis.GetAddr())
	fmt.Printf("   ✅ JWT Expire: %v\n", testConfig.JWT.GetExpireDuration())
	fmt.Printf("   ✅ Is Development: %v\n", testConfig.App.IsDevelopment())
	fmt.Printf("   ✅ Is Production: %v\n", testConfig.App.IsProduction())

	fmt.Println("\n========== 所有组件测试完成 ==========")
}

// RunTests 当传入 -test 参数时运行测试
func GetDBStatus() sql.DBStats {
	sqlDB, err := storage.DB.DB()
	if err != nil {
		panic(fmt.Sprintf("获取数据库实例失败: %v", err))
	}
	return sqlDB.Stats()
}
