package services

import (
	"context"
	"fmt"
	"github.com/looplab/fsm"
	"leetgode/internal/db/postgres"
)

type StateService struct {
	FSM            *fsm.FSM
	userRepository *postgres.UserRepository
}

func NewStateService(userRepository *postgres.UserRepository) *StateService {
	return &StateService{
		FSM: fsm.NewFSM(
			"Run",
			fsm.Events{
				{Name: "Run", Src: []string{}, Dst: "Run"},
				{
					Name: "MainMenu",
					Src: []string{
						"ShowMyStats", "ShowTodayProblem", "Run", "PleaseAuthLeetcode", "AuthLeetcode",
						"RemindersMenu", "AddReminder", "WaitingAddReminderTime", "DeleteReminder",
						"WaitingDeleteReminderTime", "ShowReminders"},
					Dst: "MainMenu",
				},
				{Name: "ShowMyStats", Src: []string{"MainMenu", "PleaseAuthLeetcode", "Run"}, Dst: "ShowMyStats"},
				{Name: "ShowTodayProblem", Src: []string{"MainMenu", "Run"}, Dst: "ShowTodayProblem"},
				{Name: "PleaseAuthLeetcode", Src: []string{"ShowMyStats", "MainMenu"}, Dst: "PleaseAuthLeetcode"},
				{Name: "AuthLeetcode", Src: []string{"PleaseAuthLeetcode"}, Dst: "AuthLeetcode"},
				{Name: "RemindersMenu", Src: []string{"MainMenu", "AddReminder", "WaitingAddReminderTime", "DeleteReminder",
					"WaitingDeleteReminderTime", "ShowReminders"}, Dst: "RemindersMenu"},
				{Name: "AddReminder", Src: []string{"RemindersMenu"}, Dst: "AddReminder"},
				{Name: "WaitingAddReminderTime", Src: []string{"AddReminder"}, Dst: "WaitingAddReminderTime"},
				{Name: "DeleteReminder", Src: []string{"RemindersMenu"}, Dst: "DeleteReminder"},
				{Name: "WaitingDeleteReminderTime", Src: []string{"DeleteReminder"}, Dst: "WaitingDeleteReminderTime"},
				{Name: "ShowReminders", Src: []string{"RemindersMenu"}, Dst: "ShowReminders"},
			},
			fsm.Callbacks{},
		),
		userRepository: userRepository,
	}
}

func (ss *StateService) TryChangeStateTo(userID int64, newState string) bool {
	currentState, err := ss.GetUserState(userID)
	if err != nil {
		return false
	}
	ss.FSM.SetState(currentState)

	err = ss.FSM.Event(context.Background(), newState)
	if err != nil {
		return false
	}

	err = ss.SetUserState(userID, newState)
	if err != nil {
		return false
	}

	return true
}

func (ss *StateService) GetUserState(userID int64) (string, error) {
	user, err := ss.userRepository.GetUserByID(userID)
	if err != nil {
		err := ss.CreateUser(userID, "Run")
		if err != nil {
			return "", fmt.Errorf("failed to create user: %w", err)
		}
		return "Run", nil
	}
	return user.CurrentState, nil
}

func (ss *StateService) CreateUser(userID int64, initialState string) error {
	err := ss.userRepository.AddUser(userID, initialState)
	if err != nil {
		return fmt.Errorf("failed to add user to the database: %w", err)
	}
	return nil
}

func (ss *StateService) SetUserState(userID int64, newState string) error {
	err := ss.userRepository.UpdateUserState(userID, newState)
	if err != nil {
		return fmt.Errorf("failed to update user state in the database: %w", err)
	}
	return nil
}
