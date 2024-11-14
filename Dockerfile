#최종 배포할 실행 파일 만드는 과정
#multi-stage build에서 이 단계의 이름을 'deploy-builder'로 지정
FROM golang:1.23.1-bullseye AS deploy-builder

#app 디렉토리 생성 및 작업 디렉토리로 설정
WORKDIR /app

#종속성 파일(go.mod go.sum)만 복사해서 캐시 활용
#소스코드가 변경되어도 종속성이 변경되지 않았다면 도커의 캐시 활용
COPY go.mod go.sum ./
#프로젝트의 모든 종속성 다운로드
RUN go mod download

#전체 소스 파일을 /app으로 복사하고 최적화된 바이너리 생성
COPY . .
RUN go build -trimmpath -ldflags="-w -s" -o app

# 최종 배포 단계 : 실제 운영 환경에서 실행될 최소한의 이미지 생성
#경량화된 Degian 이미지 사용
FROM debian:bullseye-slim AS deploy

#시스템 패키지 최신화
RUN apt-get update

#첫 번째 단계에서 만든 실행 파일만 복사
COPY --from=deploy-builder /app/app .

#컨테이너 시작 시 애플리케이션 자동 실행
CMD ["./app"]

#개발자의 로컬 환경을 위한 설정
FROM golang:1.23 AS dev

#/app 디렉토리를 작업공간으로 설정
WORKDIR /app

#air 도구 설치 (코드 변경 시 자동 재빌드 지원)
RUN go install github.com/air-verse/air@latest

#개발 서버 자동 시작
CMD ["air"]

## docker-compose down -v <= 실행중인 컨테이너를 멈추고 매핑된 볼륨들을 제거