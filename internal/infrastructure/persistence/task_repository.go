package persistence

import (
	"ai-translate/internal/domain/task"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type taskRepository struct {
	db gdb.DB
}

// NewTaskRepository 创建任务仓储实例
func NewTaskRepository() task.TaskRepository {
	return &taskRepository{
		db: g.DB(),
	}
}

func (r *taskRepository) FindByID(id uint64) (*task.Task, error) {
	var t task.Task
	err := r.db.Model("tasks").Where("id", id).Scan(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *taskRepository) FindByType(taskType int) ([]*task.Task, error) {
	var tasks []*task.Task
	err := r.db.Model("tasks").Where("type", taskType).Scan(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) FindByStatus(status int) ([]*task.Task, error) {
	var tasks []*task.Task
	err := r.db.Model("tasks").Where("status", status).Scan(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) Save(task *task.Task) error {
	_, err := r.db.Model("tasks").Insert(task)
	return err
}

func (r *taskRepository) Update(task *task.Task) error {
	_, err := r.db.Model("tasks").Where("id", task.ID).Update(task)
	return err
}

func (r *taskRepository) Delete(id uint64) error {
	_, err := r.db.Model("tasks").Where("id", id).Delete()
	return err
}

type taskQueue struct {
	redis *g.Redis
}

// NewTaskQueue 创建任务队列实例
func NewTaskQueue() task.TaskQueue {
	return &taskQueue{
		redis: g.Redis(),
	}
}

func (q *taskQueue) Push(t *task.Task) error {
	// 使用Redis的List数据结构实现队列
	_, err := q.redis.LPush("task_queue", t)
	return err
}

func (q *taskQueue) Pop() (*task.Task, error) {
	// 从队列右侧弹出任务
	var t task.Task
	err := q.redis.RPop("task_queue").Scan(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (q *taskQueue) Remove(id uint64) error {
	// 从队列中移除指定ID的任务
	// 这里需要遍历队列找到对应ID的任务并删除
	tasks, err := q.redis.LRange("task_queue", 0, -1).Result()
	if err != nil {
		return err
	}
	
	for i, taskStr := range tasks {
		var t task.Task
		if err := g.Json().DecodeTo([]byte(taskStr), &t); err != nil {
			continue
		}
		if t.ID == id {
			_, err := q.redis.LRem("task_queue", i, taskStr)
			return err
		}
	}
	return nil
}

func (q *taskQueue) GetLength() (int64, error) {
	return q.redis.LLen("task_queue").Result()
}

func (q *taskQueue) Clear() error {
	return q.redis.Del("task_queue").Err()
} 