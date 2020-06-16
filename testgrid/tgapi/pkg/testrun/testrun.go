package testrun

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedhq/kurl/testgrid/tgapi/pkg/persistence"
	"github.com/replicatedhq/kurl/testgrid/tgapi/pkg/testrun/types"
)

func TryGet(id string) (*types.TestRun, error) {
	pg := persistence.MustGetPGSession()

	query := `select count(1) from testrun where ref = $1`
	row := pg.QueryRow(query, id)

	count := 0
	if err := row.Scan(&count); err != nil {
		return nil, errors.Wrap(err, "failed to scan count")
	}

	if count == 0 {
		return nil, nil
	}

	query = `select ref, created_at from testrun where ref = $1`
	row = pg.QueryRow(query, id)

	testRun := types.TestRun{}

	if err := row.Scan(&testRun.ID, &testRun.CreatedAt); err != nil {
		return nil, errors.Wrap(err, "failed to scan")
	}

	return &testRun, nil
}

func Delete(id string) error {
	pg := persistence.MustGetPGSession()

	query := `delete from testrun where ref = $1`
	if _, err := pg.Exec(query, id); err != nil {
		return errors.Wrap(err, "failed to delete")
	}

	query = `delete from testinstance where testrun_ref = $1`
	if _, err := pg.Exec(query, id); err != nil {
		return errors.Wrap(err, "failed to delete")
	}

	return nil
}

func Create(id string) error {
	pg := persistence.MustGetPGSession()

	query := `insert into testrun (ref, created_at) values ($1, $2)`
	if _, err := pg.Exec(query, id, time.Now()); err != nil {
		return errors.Wrap(err, "failed to insert")
	}

	return nil
}

func List() ([]types.TestRun, error) {
	pg := persistence.MustGetPGSession()

	query := `select ref, created_at from testrun order by created_at desc`
	rows, err := pg.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list runs")
	}

	runs := []types.TestRun{}
	for rows.Next() {
		run := types.TestRun{}

		if err := rows.Scan(&run.ID, &run.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan run")
		}

		runs = append(runs, run)
	}

	return runs, nil
}