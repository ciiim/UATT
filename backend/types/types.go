package types

type ActionUID = int

type ActionTypeID = int

type ActionType string

type ModuleTypeID = int

type ModuleUID = int

const (
	IOAction      ActionType = "IO"
	ControlAction ActionType = "Control"
	DebugAction   ActionType = "Debug"
)

const (
	PrintAT ActionTypeID = 90
	DelayAT ActionTypeID = 91
	ShowAT  ActionTypeID = 92

	SendAT    ActionTypeID = 1
	ReceiveAT ActionTypeID = 2

	DeclareAT ActionTypeID = 23

	IfAT   ActionTypeID = 24
	ElseAT ActionTypeID = 25

	ForLabelAT ActionTypeID = 26

	EndBlockAT ActionTypeID = 27

	GotoLabelAT ActionTypeID = 28
	GotoAT      ActionTypeID = 29

	ChangeBaudRateAT ActionTypeID = 30

	StopAT ActionTypeID = 31

	AssignAT ActionTypeID = 32
)

const (
	FillMT   ModuleTypeID = 10
	FixedMT  ModuleTypeID = 11
	CalcMT   ModuleTypeID = 12
	CustomMT ModuleTypeID = 13
)
