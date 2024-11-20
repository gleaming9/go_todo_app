package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go_todo_app/entity"
)

// RegisterUser 사용자 데이터를 데이터베이스에 저장
func (r *Repository) RegisterUser(
	ctx context.Context, db Execer, u *entity.User,
) error {
	// 1. 생성 및 수정 시간을 현재 시간으로 설정
	u.Created = r.Clocker.Now()  // 생성 시간
	u.Modified = r.Clocker.Now() // 수정 시간

	// 2. 사용자 데이터를 삽입하는 SQL 쿼리
	sql := `INSERT INTO user(
        	name, password, role, created, modified
			) VALUES (?,?,?,?,?)`

	// 3. 데이터베이스에 사용자 데이터를 삽입
	result, err := db.ExecContext(ctx, sql, u.Name, u.Password, u.Role, u.Created, u.Modified)
	if err != nil { // MySQL 오류 체크
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			// 중복된 이름이 있을 경우 사용자 생성 불가 에러 반환
			return fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry)
		}
		// 기타 데이터베이스 오류 반환
		return err
	}

	// 4. 삽입된 데이터의 ID를 가져와 사용자 엔티티에 저장
	id, err := result.LastInsertId()
	if err != nil {
		return nil
	}
	u.ID = entity.UserID(id) // ID 설정

	// 5. 성공적으로 데이터베이스에 저장 완료
	return nil
}
