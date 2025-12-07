package ocr

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

type Client struct {
	client *gosseract.Client
}

func NewClient() *Client {
	client := gosseract.NewClient()

	client.SetLanguage("eng", "ind")

	client.SetPageSegMode(3)

	return &Client{
		client: client,
	}
}

func (c *Client) Close() error {
	return c.client.Close()
}

// ProcessImage performs OCR on an image file path
func (c *Client) ProcessImage(imagePath string) (string, error) {
	if err := c.client.SetImage(imagePath); err != nil {
		return "", fmt.Errorf("failed to set image for OCR: %w", err)
	}

	text, err := c.client.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	return text, nil
}

// ProcessImageFromBytes performs OCR on image data in memory
func (c *Client) ProcessImageFromBytes(imageData []byte) (string, error) {
	if err := c.client.SetImageFromBytes(imageData); err != nil {
		return "", fmt.Errorf("failed to set image bytes for OCR: %w", err)
	}

	text, err := c.client.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	return text, nil
}
