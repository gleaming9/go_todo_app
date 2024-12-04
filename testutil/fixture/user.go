package fixture

import (
	"go_todo_app/entity"
	"math/rand"
	"strconv"
	"time"
)

//테스트용 데이터 생성, 테스트 헬퍼로 더미 데이터 생성 함수를 만들어 코드 기반이 늘어나도 관리가 수월해진다.

func User(u *entity.User) *entity.User {
	//기본 값을 가지는 User 객체 생성
	result := &entity.User{
		ID:       entity.UserID(rand.Int()),                  // 랜덤 ID
		Name:     "gleaming9" + strconv.Itoa(rand.Int())[:5], // 랜덤 숫자 포함한 이름
		Password: "password",
		Role:     "admin",
		Created:  time.Now(),
		Modified: time.Now(),
	}

	if u == nil {
		return result
	}
	if u.ID != 0 {
		result.ID = u.ID
	}
	if u.Name != "" {
		result.Name = u.Name
	}
	if u.Password != "" {
		result.Password = u.Password
	}
	if u.Role != "" {
		result.Role = u.Role
	}
	if !u.Created.IsZero() {
		result.Created = u.Created
	}
	if !u.Modified.IsZero() {
		result.Modified = u.Modified
	}

	return result
}
