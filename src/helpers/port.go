package helpers

import (
	"fmt"
	"os"
	"strconv"
)

func GetPortFromEnv() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("PORT environment variable is required")
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("invalid port number: %v", err)
	}

	if portNum < 0 || portNum > 65535 {
		return "", fmt.Errorf("port number %d is out of valid range (0-65535)", portNum)
	}

	return port, nil
}
