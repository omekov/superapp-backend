package repository

import (
	"database/sql"
	"fmt"
)

type Repo struct {
	getMarksStmt          *sql.Stmt
	getModelsStmt         *sql.Stmt
	getVolumesStmt        *sql.Stmt
	getSpecificationsStmt *sql.Stmt
	db                    *sql.DB
}

func NewRepository(db *sql.DB) (Repo, error) {
	getMarksStmt, err := db.Prepare("SELECT DISTINCT mark FROM data ORDER BY popular_rate DESC;")
	if err != nil {
		return Repo{}, fmt.Errorf("getMarkStmt -> %v", err)
	}

	getModelsStmt, err := db.Prepare("SELECT DISTINCT model FROM data WHERE mark = ? ORDER BY model ASC;")
	if err != nil {
		return Repo{}, fmt.Errorf("getModelStmt -> %v", err)
	}

	getVolumesStmt, err := db.Prepare("SELECT DISTINCT volume FROM data WHERE mark = ? and model = ? ORDER BY model ASC;")
	if err != nil {
		return Repo{}, fmt.Errorf("getVolumeStmt -> %v", err)
	}

	getSpecificationsStmt, err := db.Prepare("SELECT year, amount FROM data WHERE mark = ? and model = ? and volume = ? ORDER BY year DESC;")
	if err != nil {
		return Repo{}, fmt.Errorf("getSpecificationStmt -> %v", err)
	}
	return Repo{
		getMarksStmt:          getMarksStmt,
		getModelsStmt:         getModelsStmt,
		getVolumesStmt:        getVolumesStmt,
		getSpecificationsStmt: getSpecificationsStmt,
		db:                    db,
	}, nil
}
