package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()
	r.POST("/", authHandler) // 인증 엔드포인트 등록
	r.Run(":3000")           // 서버 실행
}

// authHandler는 인증 요청을 처리하는 핸들러입니다.
func authHandler(c *gin.Context) {
	authToken, err := getEnv("AUTH_TOKEN") // 환경 변수에서 인증 토큰 가져오기
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// 요청 헤더의 Authorization 값이 일치하지 않으면 403 응답 반환
	if c.GetHeader("Authorization") != authToken {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid auth token"})
		return
	}

	awsRegion, err := getEnv("AWS_REGION") // AWS 리전 가져오기
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// AWS SDK 설정 로드
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		log.Println("Failed to load AWS config:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	dbClient := dynamodb.NewFromConfig(awsConfig) // DynamoDB 클라이언트 생성

	// 새로운 UUID와 타임스탬프 생성
	newUUID := uuid.New().String()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// DynamoDB에 저장할 아이템 구성
	item := map[string]types.AttributeValue{
		"id":        &types.AttributeValueMemberS{Value: newUUID},
		"timestamp": &types.AttributeValueMemberN{Value: timestamp},
	}

	tableName, err := getEnv("DDB_TABLE_NAME") // DynamoDB 테이블 이름 가져오기
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// DynamoDB에 아이템 저장
	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		log.Println("Failed to write to DynamoDB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// 성공 응답 반환
	c.JSON(http.StatusOK, gin.H{"message": "Authentication succeeded"})
}

// getEnv는 환경 변수를 가져오며, 없을 경우 오류를 반환합니다.
func getEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("ERROR: Environment variable '%s' is not set", key)
	}
	return value, nil
}
