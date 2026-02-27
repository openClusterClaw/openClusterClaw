严格按照AGENTS.md 中的规范进行开发
后端使用air进行热加载
开发使用 make dev 进行调试。make dev 会启动前端、后端。

# 数据库迁移规范
- 必须使用 GORM 的 AutoMigrate 方法进行数据库迁移
- 所有表结构定义必须在 Model 层使用 GORM struct tags
- 严禁使用 migrations/*.sql 文件进行数据库迁移
- 迁移代码必须写在 Go 代码中（通常在 main.go 的 runMigrations 函数）
- 示例：
  ```go
  // runMigrations 使用 GORM AutoMigrate 创建表
  func runMigrations(db *gorm.DB) error {
      return db.AutoMigrate(
          &model.User{},
          &model.Tenant{},
          &model.Project{},
          // ... 其他模型
      )
  }
  ```