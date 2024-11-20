package store

import (
	"context"
	_ "github.com/go-sql-driver/mysql" // 추가 해야 오류 안남
	"go_todo_app/entity"
)

// ListTasks 함수는 모든 할 일 목록을 반환한다.
func (r *Repository) ListTasks(
	ctx context.Context, db Queryer,
) (entity.Tasks, error) {
	tasks := entity.Tasks{} // 결과를 저장할 슬라이스 초기화
	sql := `SELECT
			id, title,
			status, created, modified
			FROM task;`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}

// AddTask는 새로운 Task를 데이터베이스에 추가
func (r *Repository) AddTask(
	ctx context.Context, db Execer, t *entity.Task,
) error {
	// 현재 시간을 Created와 Modified에 저장
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()
	sql := `INSERT INTO task
			(title, status, created, modified)
			VALUES (?, ?, ?, ?)` // task 테이블에 데이터 추가
	result, err := db.ExecContext(
		ctx, sql, t.Title, t.Status,
		t.Created, t.Modified,
	)
	if err != nil {
		return err
	}
	// 삽입된 데이터의 ID를 가져와 Task의 ID 필드에 저장
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TaskID(id)
	return nil
}
