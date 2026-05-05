package utils

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromPDF(pdfData []byte) (string, error) {
	if len(pdfData) == 0 {
		return "", errors.New("PDF data is empty")
	}

	if len(pdfData) > 50*1024*1024 {
		return "", errors.New("PDF exceeds maximum allowed size of 50 MB")
	}

	if !bytes.HasPrefix(pdfData, []byte("%PDF")) {
		return "", errors.New("Data does not appear to be a valid pdf")
	}

	reader, err := pdf.NewReader(bytes.NewReader(pdfData), int64(len(pdfData)))
	if err != nil {
		return "", fmt.Errorf("open PDF: %w", err)
	}

	var sb strings.Builder

	for i := 1; i <= reader.NumPage(); i++ {
		text, err := reader.Page(i).GetPlainText(nil)
		if err != nil {
			continue
		}
		sb.WriteString(text)
	}

	result := strings.TrimSpace(sb.String())
	if result == "" {
		return "", errors.New("no extractable text found in PDF")
	}

	return result, nil
}
