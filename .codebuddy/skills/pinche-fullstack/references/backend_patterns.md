# 后端代码模板

## 命名规范

- Handler: `xxx_handler.go` → `XxxHandler` struct
- Service: `xxx_service.go` → `XxxService` struct
- Repository: `xxx_repo.go` → `XxxRepository` struct
- Model: `xxx.go` → `Xxx` struct + `XxxReq`/`XxxResp`

## 完整示例：新增功能

### 1. Model 定义 (model/example.go)

```go
package model

type Example struct {
    ID        uint64    `json:"id"`
    Name      string    `json:"name"`
    Status    int8      `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ExampleCreateReq struct {
    Name   string `json:"name" binding:"required"`
    Status int8   `json:"status"`
}

type ExampleUpdateReq struct {
    Name   string `json:"name"`
    Status *int8  `json:"status"`
}

type ExampleListReq struct {
    Status   *int8  `form:"status"`
    Search   string `form:"search"`
    Page     int    `form:"page,default=1"`
    PageSize int    `form:"page_size,default=20"`
}
```

### 2. Repository 数据访问 (repository/example_repo.go)

```go
package repository

import (
    "pinche/internal/database"
    "pinche/internal/model"
)

type ExampleRepository struct{}

func NewExampleRepository() *ExampleRepository {
    return &ExampleRepository{}
}

func (r *ExampleRepository) GetByID(id uint64) (*model.Example, error) {
    query := `SELECT id, name, status, created_at, updated_at 
        FROM examples WHERE id = ?`
    
    row := database.DB.QueryRow(query, id)
    
    ex := &model.Example{}
    err := row.Scan(&ex.ID, &ex.Name, &ex.Status, &ex.CreatedAt, &ex.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return ex, nil
}

func (r *ExampleRepository) Create(ex *model.Example) error {
    query := `INSERT INTO examples (name, status, created_at, updated_at) 
        VALUES (?, ?, NOW(), NOW())`
    
    result, err := database.DB.Exec(query, ex.Name, ex.Status)
    if err != nil {
        return err
    }
    
    id, _ := result.LastInsertId()
    ex.ID = uint64(id)
    return nil
}

func (r *ExampleRepository) Update(ex *model.Example) error {
    query := `UPDATE examples SET name = ?, status = ?, updated_at = NOW() 
        WHERE id = ?`
    
    _, err := database.DB.Exec(query, ex.Name, ex.Status, ex.ID)
    return err
}

func (r *ExampleRepository) Delete(id uint64) error {
    query := `DELETE FROM examples WHERE id = ?`
    _, err := database.DB.Exec(query, id)
    return err
}

func (r *ExampleRepository) List(req *model.ExampleListReq) ([]*model.Example, int, error) {
    var conditions []string
    var args []interface{}
    
    if req.Status != nil {
        conditions = append(conditions, "status = ?")
        args = append(args, *req.Status)
    }
    
    if req.Search != "" {
        conditions = append(conditions, "name LIKE ?")
        args = append(args, "%"+req.Search+"%")
    }
    
    whereClause := ""
    if len(conditions) > 0 {
        whereClause = "WHERE " + strings.Join(conditions, " AND ")
    }
    
    // 查询总数
    countQuery := "SELECT COUNT(*) FROM examples " + whereClause
    var total int
    database.DB.QueryRow(countQuery, args...).Scan(&total)
    
    // 分页查询
    offset := (req.Page - 1) * req.PageSize
    listQuery := fmt.Sprintf(`SELECT id, name, status, created_at, updated_at 
        FROM examples %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereClause)
    
    args = append(args, req.PageSize, offset)
    rows, err := database.DB.Query(listQuery, args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    var list []*model.Example
    for rows.Next() {
        ex := &model.Example{}
        rows.Scan(&ex.ID, &ex.Name, &ex.Status, &ex.CreatedAt, &ex.UpdatedAt)
        list = append(list, ex)
    }
    
    return list, total, nil
}
```

### 3. Service 业务逻辑 (service/example_service.go)

```go
package service

import (
    "pinche/internal/model"
    "pinche/internal/repository"
)

type ExampleService struct {
    repo *repository.ExampleRepository
}

func NewExampleService() *ExampleService {
    return &ExampleService{
        repo: repository.NewExampleRepository(),
    }
}

func (s *ExampleService) GetByID(id uint64) (*model.Example, error) {
    return s.repo.GetByID(id)
}

func (s *ExampleService) Create(req *model.ExampleCreateReq) (*model.Example, error) {
    ex := &model.Example{
        Name:   req.Name,
        Status: req.Status,
    }
    
    if err := s.repo.Create(ex); err != nil {
        return nil, err
    }
    return ex, nil
}

func (s *ExampleService) Update(id uint64, req *model.ExampleUpdateReq) (*model.Example, error) {
    ex, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    if req.Name != "" {
        ex.Name = req.Name
    }
    if req.Status != nil {
        ex.Status = *req.Status
    }
    
    if err := s.repo.Update(ex); err != nil {
        return nil, err
    }
    return ex, nil
}

func (s *ExampleService) Delete(id uint64) error {
    return s.repo.Delete(id)
}

func (s *ExampleService) List(req *model.ExampleListReq) (*model.ListResp, error) {
    list, total, err := s.repo.List(req)
    if err != nil {
        return nil, err
    }
    
    return &model.ListResp{
        List:     list,
        Total:    total,
        Page:     req.Page,
        PageSize: req.PageSize,
    }, nil
}
```

### 4. Handler HTTP处理 (handler/example_handler.go)

```go
package handler

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "pinche/internal/model"
    "pinche/internal/service"
)

type ExampleHandler struct {
    service *service.ExampleService
}

func NewExampleHandler() *ExampleHandler {
    return &ExampleHandler{
        service: service.NewExampleService(),
    }
}

func (h *ExampleHandler) Get(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "Invalid ID"))
        return
    }
    
    result, err := h.service.GetByID(id)
    if err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeNotFound, "Not found"))
        return
    }
    
    c.JSON(http.StatusOK, model.Success(result))
}

func (h *ExampleHandler) Create(c *gin.Context) {
    var req model.ExampleCreateReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, err.Error()))
        return
    }
    
    result, err := h.service.Create(&req)
    if err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "Failed to create"))
        return
    }
    
    c.JSON(http.StatusOK, model.Success(result))
}

func (h *ExampleHandler) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "Invalid ID"))
        return
    }
    
    var req model.ExampleUpdateReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, err.Error()))
        return
    }
    
    result, err := h.service.Update(id, &req)
    if err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "Failed to update"))
        return
    }
    
    c.JSON(http.StatusOK, model.Success(result))
}

func (h *ExampleHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "Invalid ID"))
        return
    }
    
    if err := h.service.Delete(id); err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "Failed to delete"))
        return
    }
    
    c.JSON(http.StatusOK, model.Success(nil))
}

func (h *ExampleHandler) List(c *gin.Context) {
    var req model.ExampleListReq
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, err.Error()))
        return
    }
    
    result, err := h.service.List(&req)
    if err != nil {
        c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "Failed to list"))
        return
    }
    
    c.JSON(http.StatusOK, model.Success(result))
}
```

### 5. 路由注册 (router/router.go)

```go
// 在 SetupRouter 函数中添加
exampleHandler := handler.NewExampleHandler()
api.GET("/examples", exampleHandler.List)
api.GET("/examples/:id", exampleHandler.Get)
api.POST("/examples", exampleHandler.Create)
api.PUT("/examples/:id", exampleHandler.Update)
api.DELETE("/examples/:id", exampleHandler.Delete)
```

## 安全要点

1. **必须使用参数化查询**
```go
// 正确
query := "SELECT * FROM users WHERE id = ?"
db.Query(query, id)

// 错误 - SQL 注入风险
query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", id)
```

2. **权限验证**
```go
// 从 JWT 获取当前用户
userID := c.GetUint64("user_id")

// 验证操作权限
if trip.UserID != userID {
    c.JSON(http.StatusOK, model.Error(model.ErrCodeForbidden, "无权操作"))
    return
}
```
