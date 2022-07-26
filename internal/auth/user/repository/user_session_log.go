package repository

import (
	"context"
)

func (r *userRepo) CreateUserSessionLog(ctx context.Context, log UserSessionLog) (uint64, error) {
	r.rw.Lock()
	defer r.rw.Unlock()

	var userSessionLogID uint64
	if r.db == nil {
		return userSessionLogID, ErrNotConnection
	}

	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO user_session_logs
	(
		session_id,
		username,
		user_agent,
		client_ip,
		http_method,
		http_path,
		http_req_body
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`)
	if err != nil {
		return userSessionLogID, err
	}

	err = stmt.QueryRowContext(ctx,
		log.SessionID,
		log.Username,
		log.UserAgent,
		log.ClientIP,
		log.HTTPMethod,
		log.HTTPPath,
		log.HTTPReqBody,
	).Scan(&userSessionLogID)
	if err != nil {
		return userSessionLogID, err
	}

	return userSessionLogID, nil
}

func (r *userRepo) UpdateUserSessionLog(ctx context.Context, userSessionLogID uint64, log UserSessionLog) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.db == nil {
		return ErrNotConnection
	}

	stmt, err := r.db.PrepareContext(ctx, `UPDATE user_session_logs SET
		http_status=$2,
		http_res_body=$3
	WHERE id=$1`)
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, userSessionLogID, log.HTTPStatus, log.HTTPResBody)
	if err != nil {
		return err
	}

	rowNum, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowNum == 0 {
		return ErrNoRowsUpdated
	}
	return nil
}
