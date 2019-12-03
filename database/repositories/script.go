package repositories

import (
	"gpu-demonstration-api/database"
	"gpu-demonstration-api/models"
)

type Script struct{}

func (sr *Script) Insert(s models.Script) error {
	_, err := database.DB.Exec(
		`INSERT INTO "SCRIPT" (
			"M_NAME",
			"DESCRIPTION",
			"PROCESSOR",
			"LOCATION_PATH",
			"LOCATION_TYPE"
		) VALUES (
			:1, :2, :3, :4, :5
		)`,
		s.Name,
		s.Description,
		s.Processor,
		s.LocationPath,
		s.LocationType,
	)
	if err != nil {
		return err
	}
	return nil
}

func (sr *Script) Update(s models.Script) error {
	_, err := database.DB.Exec(
		`UPDATE "SCRIPT" SET
			"M_NAME" = :1,
			"DESCRIPTION" = :2,
			"PROCESSOR" = :3,
			"LOCATION_PATH" = :4,
			"LOCATION_TYPE" = :5
		WHERE SCRIPT_ID = :6`,
		s.Name,
		s.Description,
		s.Processor,
		s.LocationPath,
		s.LocationType,
		s.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (sr *Script) Delete(scriptID int) error {
	_, err := database.DB.Exec(
		`DELETE FROM "SCRIPT" WHERE "SCRIPT_ID" = :1`,
		scriptID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (sr *Script) SelectAll() ([]models.Script, error) {
	var ss []models.Script
	rows, err := database.DB.Query(
		`SELECT
			"SCRIPT_ID",
			"M_NAME",
			"PROCESSOR",
			"DESCRIPTION",
			"LOCATION_PATH",
			"LOCATION_TYPE"
		FROM "SCRIPT"
		WHERE "IS_ACTIVE" = 1`,
	)
	if err != nil {
		return ss, err
	}

	defer rows.Close()
	for rows.Next() {
		var s models.Script
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Processor,
			&s.Description,
			&s.LocationPath,
			&s.LocationType,
		)
		if err != nil {
			return ss, err
		}
		ss = append(ss, s)
	}
	err = rows.Err()
	if err != nil {
		return ss, err
	}
	return ss, nil
}

func (sr *Script) SelectByID(scriptID int) (models.Script, error) {
	var s models.Script
	rows, err := database.DB.Query(
		`SELECT
			"SCRIPT_ID",
			"M_NAME",
			"PROCESSOR",
			"DESCRIPTION",
			"LOCATION_PATH",
			"LOCATION_TYPE"
		FROM "SCRIPT"
		WHERE "SCRIPT_ID" = :1`,
		scriptID,
	)
	if err != nil {
		return s, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Processor,
			&s.Description,
			&s.LocationPath,
			&s.LocationType,
		)
		if err != nil {
			return s, err
		}
	}
	err = rows.Err()
	if err != nil {
		return s, err
	}
	return s, nil
}
