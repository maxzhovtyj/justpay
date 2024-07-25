package order

import "github.com/google/uuid"

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

type Order struct {
	ID     uuid.UUID
	Status Status
}

type Event struct {
}
