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
	ServerUrl string `json:"serverUrl"`
	LastSync  string `json:"lastSync"`
	UserEmail string `json:"userEmail"`
	UserId    string `json:"userId"`
	Status    Status `json:"status"`
}

type BWManager struct {
	items       []Item
	token       string
	VaultStatus VaultStatus
}

func NewBWManager() *BWManager {
	var bwm BWManager
	bwm.token = os.Getenv("BW_SESSION")
	bwm.UpdateStatus()
	return &bwm
}

func (bwm *BWManager) Login(un string, pw string) error {
	if bwm.VaultStatus.Status != Unauthenticated {
		return nil
	}
	out, err := exec.Command("bw", "login", un, pw, "--raw").Output()
	if err != nil {
		return errors.New(err.Error())
	}
	bwm.UpdateStatus()
	(*bwm).token = string(out)
	return nil
}

func (bwm *BWManager) Unlock(pw string) error {
	if bwm.VaultStatus.Status == Unauthenticated {
		return errors.New("Not Logged in")
	}
	out, err := exec.Command("bw", "unlock", pw, "--raw").Output()
	if err != nil {
		return errors.New(err.Error())
	}
	bwm.UpdateStatus()
	(*bwm).token = string(out)
	return nil
}

func (bwm *BWManager) Logout() error {
	if bwm.VaultStatus.Status == Unauthenticated {
		return errors.New("Not Logged in")
	}
	_, err := exec.Command("bw", "logout").Output()
	if err != nil {
		return errors.New(err.Error())
	}
	(*bwm).token = ""
	bwm.UpdateStatus()
	return nil
}

func (bwm *BWManager) UpdateStatus() error {
	out, err := exec.Command("bw", "status").Output()
	if err != nil {
		return errors.New(err.Error())
	}
	json.Unmarshal(out, &bwm.VaultStatus)
	return nil
}

func (bwm *BWManager) UpdateList() error {
	if bwm.VaultStatus.Status == Unauthenticated {
		fmt.Println("Not Logged In")
		return errors.New("Not Logged In")
	}
	out, err := exec.Command("bw", "list", "items", "--session", bwm.token).Output()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return errors.New(err.Error())
	}
	json.Unmarshal(out, &bwm.items)

	return nil
}

func (bwm *BWManager) GetList() ([]Item, error) {
	if bwm.VaultStatus.Status == Unauthenticated {
		return nil, errors.New("Not Logged In")
	}
	return bwm.items, nil
}

func (bwm *BWManager) GetPassword(id string) (string, error) {
	for _, v := range bwm.items {
		if v.ID == id {
			return v.Login.Password, nil
		}
	}
	return "", errors.New("No entry matching ID found.")
}
