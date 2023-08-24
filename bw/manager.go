package bw

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type ItemType int

const (
	Login      ItemType = 1
	SecureNote ItemType = 2
	Card       ItemType = 3
	Identity   ItemType = 4
)

func (it ItemType) String() string {
	switch it {
	case Login:
		return "Login"
	case SecureNote:
		return "SecureNote"
	case Card:
		return "Card"
	case Identity:
		return "Identity"
	default:
		panic("undefined item type")
	}
}

type Status string

const (
	Unlocked        Status = "unlocked"
	Locked          Status = "locked"
	Unauthenticated Status = "unauthenticated"
)

type ItemLoginURI struct {
	URI string `json:"uri"`
}

type ItemLogin struct {
	URIs                 []ItemLoginURI `json:"uris"`
	Username             string         `json:"username"`
	Password             string         `json:"password"`
	PasswordRevisionDate time.Time      `json:"passwordRevisionDate"`
	// TODO: totp
}

type Item struct {
	Object         string    `json:"object"` // TODO: enum
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	FolderID       string    `json:"folderId"`
	Type           ItemType  `json:"type"`
	Reprompt       int       `json:"reprompt"`
	Name           string    `json:"name"`
	Notes          string    `json:"notes"`
	Favorite       bool      `json:"favorite"`
	Login          ItemLogin `json:"login"`
	RevisionDate   time.Time `json:"revisionDate"`
	CreationDate   time.Time `json:"creationDate"`
	DeletedDate    time.Time `json:"deletedDate"`
}

type VaultStatus struct {
	ServerURL string `json:"serverUrl"`
	LastSync  string `json:"lastSync"`
	UserEmail string `json:"userEmail"`
	UserID    string `json:"userId"`
	Status    Status `json:"status"`
}

type Manager struct {
	items       []Item
	token       string
	VaultStatus VaultStatus
}

var ErrNotLoggedIn = errors.New("not logged in")

func NewBWManager() *Manager {
	var bwm Manager
	bwm.token = os.Getenv("BW_SESSION")
	return &bwm
}

func (bwm *Manager) Login(un string, pw string) error {
	if bwm.VaultStatus.Status != Unauthenticated {
		return nil
	}
	out, err := exec.Command("bw", "login", un, pw, "--raw").Output()
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	bwm.token = string(out)
	err = bwm.UpdateStatus()
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	return nil
}

func (bwm *Manager) Unlock(pw string) error {
	if bwm.VaultStatus.Status == Unauthenticated {
		return ErrNotLoggedIn
	}
	out, err := exec.Command("bw", "unlock", pw, "--raw").Output()
	if err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}
	bwm.token = string(out)
	err = bwm.UpdateStatus()
	if err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}
	return nil
}

func (bwm *Manager) Logout() error {
	if bwm.VaultStatus.Status == Unauthenticated {
		return ErrNotLoggedIn
	}
	_, err := exec.Command("bw", "logout").Output()
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	bwm.token = ""
	err = bwm.UpdateStatus()
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

func (bwm *Manager) UpdateStatus() error {
	out, err := exec.Command("bw", "status").Output()
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	err = json.Unmarshal(out, &bwm.VaultStatus)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}

func (bwm *Manager) UpdateList() error {
	if bwm.VaultStatus.Status == Unauthenticated {
		return ErrNotLoggedIn
	}
	out, err := exec.Command("bw", "list", "items", "--session", bwm.token).Output() // #nosec G204
	if err != nil {
		return fmt.Errorf("failed to update list: %w", err)
	}
	err = json.Unmarshal(out, &bwm.items)
	if err != nil {
		return fmt.Errorf("failed to update list: %w", err)
	}
	return nil
}

func (bwm *Manager) GetList() ([]Item, error) {
	if bwm.VaultStatus.Status == Unauthenticated {
		return nil, ErrNotLoggedIn
	}
	return bwm.items, nil
}
