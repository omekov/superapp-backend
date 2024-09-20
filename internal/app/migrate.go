package app

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/xuri/excelize/v2"
)

func startMigrate(db *sql.DB, kgdURL string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Не удалось создать драйвер миграции: %v", err)
	}

	// 3. Создаём источник миграций (например, файловый источник)
	sourceDriver, err := (&file.File{}).Open("file://../../migrations")
	if err != nil {
		return fmt.Errorf("Не удалось открыть источник миграций: %v", err)
	}

	// 4. Инициализируем объект миграции с использованием драйвера и источника
	m, err := migrate.NewWithInstance(
		"file",       // Имя источника
		sourceDriver, // Экземпляр источника
		"sqlite3",    // Имя базы данных
		driver,       // Экземпляр драйвера базы данных
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	fileBytes, err := downloadFile(kgdURL)
	if err != nil {
		return fmt.Errorf("downloadFile: %v", err)
	}

	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		return fmt.Errorf("Ошибка при открытии Excel-файла: %v", err)
	}
	defer f.Close()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return fmt.Errorf("Файл не содержит листов")
	}

	// Выводим все названия листов
	fmt.Println("Доступные листы в файле:", sheetList)
	if len(sheetList) == 0 {
		return fmt.Errorf("excel нет листов")
	}

	// Пример чтения данных из файла
	rows, err := f.GetRows(sheetList[0])
	if err != nil {
		return fmt.Errorf("GetRows:%s", err.Error())
	}

	popularMark := map[string]int{"BMW": 1000000, "NISSAN": 1000001, "KIA": 1000002, "VOLKSWAGEN": 1000003, "MERCEDES-BENZ": 1000004, "HYUNDAI": 1000005, "TOYOTA": 1000006}
	for i, row := range rows {
		rate := i + 1
		if i == 0 {
			continue
		}
		id, err := strconv.Atoi(strings.Trim(strings.Replace(row[0], ",", "", -1), " "))
		if err != nil {
			return err
		}
		mark := strings.Trim(row[1], " ")
		model := strings.Trim(row[2], " ")
		motor := strings.Trim(strings.Replace(row[3], ",", "", -1), " ")
		var volume int
		if strings.Contains(strings.ToLower(motor), strings.ToLower("Элект")) {
			volume = 0
		} else {
			volume, err = strconv.Atoi(strings.Trim(strings.Replace(row[3], ",", "", -1), " "))
			if err != nil {
				return err
			}
		}
		year, err := strconv.Atoi(strings.Trim(strings.Replace(row[4], ",", "", -1), " "))
		if err != nil {
			return err
		}
		amount, err := strconv.Atoi(strings.Trim(strings.Replace(row[5], ",", "", -1), " "))
		if err != nil {
			return err
		}
		if _, ok := popularMark[mark]; ok {
			rate = popularMark[mark]
		}
		_, err = db.Exec(
			"INSERT INTO data (id, mark, model, volume, year, amount, popular_rate) VALUES (?, ?, ?, ?, ?, ?, ?)",
			id,
			mark,
			model,
			volume,
			year,
			amount,
			rate,
		)
		if err != nil {
			return err
		}
		// модели которые
		// models = []string{
		// 	"camry","corolla", "land cruiser prado", "land cruiser", "rav4", "highlander",
		// 	"elantra", "sonata", "accent", "tucson", "santa fe", "grandeur", "creata", "palisade",
		// 	"E200", "E220", "E230", "E280", "E320", "S320", "S500", "190",
		// 	"Passat",
		// }
		// db.QueryContext(ctx, `SELECT year, model`)
	}

	return nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
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
