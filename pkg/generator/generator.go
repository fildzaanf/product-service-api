package generator

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"path/filepath"
	"time"
)

func GenerateEmailTemplate(fileTemplate string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("go-commerce-api/pkg/email/template/%s", fileTemplate))
	if err != nil {
		return "", errors.New("invalid template file")
	}

	emailTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var templateBuffer bytes.Buffer
	if err := emailTemplate.Execute(&templateBuffer, data); err != nil {
		return "", err
	}

	return templateBuffer.String(), nil
}

func GeneratePaymentCode() string {
	rand.Seed(time.Now().UnixNano())

	part1 := rand.Intn(9000) + 1000
	part2 := rand.Intn(9000) + 1000
	part3 := rand.Intn(9000) + 1000

	return fmt.Sprintf("GC-%d-%d-%d", part1, part2, part3)
}
