package scrimmage

type BetOutcome string

const (
	BetOutcome_Win       BetOutcome = "win"
	BetOutcome_Lose      BetOutcome = "lose"
	BetOutcome_Push      BetOutcome = "push"
	BetOutcome_Cashout   BetOutcome = "cashout"
	BetOutcome_Postponed BetOutcome = "postponed"
)

type BetType string

const (
	BetType_Single  BetType = "single"
	BetType_Parlays BetType = "parlays"
)

type SingleBetType string

const (
	SingleBetType_Over      SingleBetType = "over"
	SingleBetType_Under     SingleBetType = "under"
	SingleBetType_Spread    SingleBetType = "spread"
	SingleBetType_Moneyline SingleBetType = "moneyline"
	SingleBetType_Prop      SingleBetType = "prop"
)

type BetLeague string

type BetSport string

// BetDate is milisecond the number of milliseconds elapsed since the epoch
type BetDate int64

type BetDataType string

const (
	BetDataType_BetExecuted = "betExecuted"
	BetDataType_BetMade     = "betMade"
)

type RewardeHeader struct {
	ID     *string `json:"id"`
	Type   *string `json:"type"`
	UserID *string `json:"userId"`
}

type BetEvent struct {
	BetType     BetType     `json:"betType"`
	Odds        float64     `json:"odds"`
	Description string      `json:"description"`
	WagerAmount float64     `json:"wagerAmount"`
	NetProfit   *float64    `json:"netProfit,omitempty"`
	Outcome     *BetOutcome `json:"outcome,omitempty"`
	IsLive      bool        `json:"isLive"`
	BetDate     BetDate     `json:"betDate"`
	Bets        []SingleBet `json:"bets"`
}

type SingleBet struct {
	Type           SingleBetType `json:"type"`
	Odds           float64       `json:"odds"`
	TeamBetOn      *string       `json:"teamBetOn"`
	TeamBetAgainst *string       `json:"teamBetAgainst"`
	Player         *string       `json:"player"`
	League         BetLeague     `json:"league"`
	Sport          BetSport      `json:"sport"`
}

type CreateIntegrationRewardRequest struct {
	EventID  *string     `json:"eventId,omitempty"`
	UserID   string      `json:"userId"`
	DataType BetDataType `json:"dataType"`
	Body     BetEvent    `json:"body"`
}

type CreateIntegrationRewardResponse struct {
	Namespace string       `json:"namespace"`
	UserId    string       `json:"userId"`
	DataType  *BetDataType `json:"dataType"`
	EventID   *string      `json:"eventId"`
	Body      *BetEvent    `json:"body"`
}
