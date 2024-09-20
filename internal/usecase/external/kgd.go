package external

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c Client) DownloadFile() ([]byte, error) {
	resp, err := c.httpClient.Get(c.kgdURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка загрузки файла: статус %d", resp.StatusCode)
	}
	// Читаем содержимое файла в память
	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении данных файла: %v", err)
	}

	return fileBytes, nil
}
