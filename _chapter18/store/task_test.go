package store

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"go_todo_app/clock"
	"go_todo_app/entity"
	"go_todo_app/testutil"
	"testing"
)

/*
prepareTasks				테스트 환경을 초기화하고, task 테이블에 데이터를 미리 삽입.
TestRepository_ListTasks	ListTasks 메서드가 테이블의 데이터를 올바르게 조회하는지 테스트.
TestRepository_AddTask		AddTask 메서드가 데이터를 올바르게 추가하는지 테스트하며, sqlmock을 사용해 데이터베이스 동작을 모킹.
*/

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	// entity.Task를 작성하는 다른 테스트 케이스와 섞이면 테스트가 실패한다.
	// 이를 위해 트랜잭션을 적용해서 테스트 케이스 내로 한정된 테이블 상태를 만든다.
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	// 이 테스트 케이스가 끝나면 원래 상태로 되돌린다.
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	// task 테이블에 테스트 데이터를 삽입
	wants := prepareTasks(ctx, t, tx)

	// 데이터베이스에서 데이터를 가져온다.
	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 가져온 결과와 기대값을 비교
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()
	// 테이블 초기화
	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}
	c := clock.FixedClocker{}
	//미리 정의된 wants 데이터 (3개의 task)
	wants := entity.Tasks{
		{
			Title: "want task 1", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 2", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 3", Status: "done",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	// SQL INSERT 데이터 삽입
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, status, created, modified)
			VALUES
			    (?, ?, ?, ?),
			    (?, ?, ?, ?),
			    (?, ?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified,
		wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified,
		wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	// 데이터베이스에서 생성된 ID를 가져와 wants의 ID 필드에 저장
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	// 추가할 테스트 데이터를 정의
	okTask := &entity.Task{
		Title:    "ok task",
		Status:   "todo",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	// sqlmock을 사용해서 가짜 데이터베이스 생성
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	mock.ExpectExec(
		// 이스케이프 필요
		`INSERT INTO task \(title, status, created, modified\) VALUES \(\?, \?, \?, \?\)`,
	).WithArgs(okTask.Title, okTask.Status, okTask.Created, okTask.Modified).
		// WithArgs를 통해 입력으로 전달되는 값을 검증
		WillReturnResult(sqlmock.NewResult(wantID, 1)) // 삽입된 데이터의 ID를 반환

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	// 데이터가 추가되는지 확인
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}
