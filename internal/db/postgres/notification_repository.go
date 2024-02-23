package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) AddNotificationTimeIfNotExists(userID int64, notificationTime time.Time) error {
	existingTimes, err := r.GetNotificationTimes(userID)
	if err != nil {
		return fmt.Errorf("failed to check existing notification times: %w", err)
	}

	for _, existingTime := range existingTimes {
		if existingTime.Equal(notificationTime) {
			return fmt.Errorf("notification time already exists")
		}
	}

	err = r.AddNotificationTime(userID, notificationTime)
	if err != nil {
		return fmt.Errorf("failed to add notification time in the database: %w", err)
	}

	return nil
}

func (r *NotificationRepository) AddNotificationTime(userID int64, notificationTime time.Time) error {
	query := "INSERT INTO notifications (user_id, notification_time) VALUES ($1, $2)"
	_, err := r.db.Exec(query, userID, notificationTime)
	if err != nil {
		return fmt.Errorf("failed to add notification time in the database: %w", err)
	}

	return nil
}

func (r *NotificationRepository) GetNotificationTimes(userID int64) ([]time.Time, error) {
	query := "SELECT notification_time FROM notifications WHERE user_id = $1"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification times: %w", err)
	}
	defer rows.Close()

	var notificationTimes []time.Time
	for rows.Next() {
		var notificationTime time.Time
		err := rows.Scan(&notificationTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification time: %w", err)
		}
		notificationTimes = append(notificationTimes, notificationTime)
	}

	return notificationTimes, nil
}

func (r *NotificationRepository) RemoveAllNotificationTimes(userID int64) error {
	query := "DELETE FROM notifications WHERE user_id = $1"
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to remove notification times: %w", err)
	}

	return nil
}

func (r *NotificationRepository) GetUsersWithNotification(currentTime time.Time) ([]UserNotificationTime, error) {
	query := `
		SELECT user_id, notification_time
		FROM notifications
		WHERE notification_time = $1
	`
	rows, err := r.db.Query(query, currentTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get users with notifications: %w", err)
	}
	defer rows.Close()

	var usersWithNotifications []UserNotificationTime
	for rows.Next() {
		var userWithNotification UserNotificationTime
		err := rows.Scan(&userWithNotification.UserID, &userWithNotification.NotificationTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user with notification time: %w", err)
		}
		usersWithNotifications = append(usersWithNotifications, userWithNotification)
	}

	return usersWithNotifications, nil
}

func (r *NotificationRepository) RemoveNotificationTime(userID int64, notificationTime time.Time) error {
	query := "DELETE FROM notifications WHERE user_id = $1 AND notification_time = $2"
	_, err := r.db.Exec(query, userID, notificationTime)
	if err != nil {
		return fmt.Errorf("failed to remove notification time: %w", err)
	}

	return nil
}

type UserNotificationTime struct {
	UserID           int64
	NotificationTime time.Time
}
