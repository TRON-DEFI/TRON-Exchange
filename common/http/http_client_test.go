package http

import (
	"fmt"
	"testing"
)

func TestAA(t *testing.T) {
	url := "http://18.216.57.65:20110/api/dapp/querySmart"
	param := "startTime=9898&endTime=1540535940011&contractAddress=TZEuGDUMUpBeNsdeh8vRs7Af9ZvhF4hGft&start=0&limit=3"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDA3OTQ5OTgsImlhdCI6MTU0MDc5MzE5OCwiaWQiOjEsIm5iZiI6MTU0MDc5MzE5OCwidXNlcm5hbWUiOiJ0cm9uIn0.OuqIoCQtGfXr1DPgseIazrbrYUhrOKE8PGPyDRf8OUs"
	url += "?" + param
	result, err := Get(url, token, true)
	if err != nil {
		return
	}
	fmt.Printf("result:%s\n", result)
}
