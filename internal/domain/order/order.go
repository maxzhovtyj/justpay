package order

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrEventAlreadyExists  = errors.New("event already exists")
	ErrFinalStatusReceived = errors.New("final status received")
)

// "cool_order_created|confirmed_by_mayor|sbu_verification_pending|changed_my_mind|failed|give_my_money_back|chinazes",

type Status int

const (
	CoolOrderCreatedStatus       Status = 0
	ConfirmedByMayorStatus              = 1
	SBUVerificationPendingStatus        = 2
	ChangedMyMindStatus                 = 3
	FailedStatus                        = 4
	GiveMyMoneyBackStatus               = 5
	ChinazezStatus                      = 6
)

var statusNameToCode = map[string]Status{
	"cool_order_created":       CoolOrderCreatedStatus,
	"confirmed_by_mayor":       ConfirmedByMayorStatus,
	"sbu_verification_pending": SBUVerificationPendingStatus,
	"changed_my_mind":          ChangedMyMindStatus,
	"failed":                   FailedStatus,
	"give_my_money_back":       GiveMyMoneyBackStatus,
	"chinazes":                 ChinazezStatus,
}

var isFinalStatus = map[Status]bool{
	CoolOrderCreatedStatus:       false,
	ConfirmedByMayorStatus:       false,
	SBUVerificationPendingStatus: false,
	ChangedMyMindStatus:          true,
	FailedStatus:                 true,
	GiveMyMoneyBackStatus:        true,
	ChinazezStatus:               false,
}

func (s Status) IsFinal() bool {
	return isFinalStatus[s]
}

var statusCodeToName = map[Status]string{
	CoolOrderCreatedStatus:       "cool_order_created",
	ConfirmedByMayorStatus:       "confirmed_by_mayor",
	SBUVerificationPendingStatus: "sbu_verification_pending",
	ChangedMyMindStatus:          "changed_my_mind",
	FailedStatus:                 "failed",
	GiveMyMoneyBackStatus:        "give_my_money_back",
	ChinazezStatus:               "chinazes",
}

func (s Status) String() string {
	return statusCodeToName[s]
}

func GetStatus(s string) Status {
	return statusNameToCode[s]
}

type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Status    Status
	IsFinal   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Event struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	UserID    uuid.UUID
	Status    Status
	IsFinal   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FiltrationOptions struct {
	Status
	UserID    string
	Limit     string
	Offset    string
	IsFinal   string
	SortBy    string
	SortOrder string
}

func NewEvent(eventID, orderID, userID, status string, createdAt, updatedAt time.Time) (Event, error) {
	id, err := uuid.Parse(eventID)
	if err != nil {
		return Event{}, err
	}

	oid, err := uuid.Parse(orderID)
	if err != nil {
		return Event{}, err
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return Event{}, err
	}

	s := GetStatus(status)

	return Event{
		ID:        id,
		OrderID:   oid,
		UserID:    uid,
		Status:    s,
		IsFinal:   s.IsFinal(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
