package main

import (
	"context"
	"log"
	"time"

	scrimmage "github.com/Scrimmage-co/golang-sdk"
)

func main() {
	// init SDK
	sdk, err := scrimmage.InitRewarder(context.Background(),
		"API_SERVER_ENDPOINT",
		"PRIVATE_KEY",
		"production",
		scrimmage.WithSecure(true),
		scrimmage.WithValidateAPIServerEndpoint(true),
	)

	if err != nil {
		log.Fatalf("failed to initialize sdk, reason : %s", err)
	}

	// get user token
	userToken, err := sdk.User.GetUserToken(context.Background(), scrimmage.GetUserTokenRequest{
		UserID: "userId",
	})
	if err != nil {
		log.Fatalf("failed to get user token, reason : %s", err)
	}

	log.Println(userToken)

	// track reward
	trackResult, err := sdk.Reward.TrackRewardable(
		context.Background(),
		"userId",
		scrimmage.BetDataType_BetExecuted,
		scrimmage.BetEvent{
			BetType:     scrimmage.BetType_Single,
			IsLive:      false,
			Odds:        1.5,
			Description: "lorem ipsum",
			WagerAmount: 1000,
			NetProfit:   scrimmage.GetPtrOf[float64](500),
			Outcome:     scrimmage.GetPtrOf[scrimmage.BetOutcome]("win"),
			BetDate:     scrimmage.BetDate(time.Now().UnixMilli()),
			Bets: []scrimmage.SingleBet{
				{
					Type:           scrimmage.SingleBetType_Spread,
					Odds:           1.5,
					TeamBetOn:      scrimmage.GetPtrOf("team a"),
					TeamBetAgainst: scrimmage.GetPtrOf("team b"),
					League:         scrimmage.BetLeague("nba"),
					Sport:          scrimmage.BetSport("basketball"),
				},
			},
		},
	)

	if err != nil {
		log.Fatalf("failed to track reward, reason : %s", err)
	}

	log.Println(trackResult)

	// track reward once
	trackResultOnce, err := sdk.Reward.TrackRewardableOnce(
		context.Background(),
		"userId",
		scrimmage.BetDataType_BetExecuted,
		scrimmage.GetPtrOf("UniqueeventId"),
		scrimmage.BetEvent{
			BetType:     scrimmage.BetType_Single,
			IsLive:      false,
			Odds:        1.5,
			Description: "lorem ipsum",
			WagerAmount: 1000,
			NetProfit:   scrimmage.GetPtrOf[float64](500),
			Outcome:     scrimmage.GetPtrOf[scrimmage.BetOutcome]("win"),
			BetDate:     scrimmage.BetDate(time.Now().UnixMilli()),
			Bets: []scrimmage.SingleBet{
				{
					Type:           scrimmage.SingleBetType_Spread,
					Odds:           1.5,
					TeamBetOn:      scrimmage.GetPtrOf("team a"),
					TeamBetAgainst: scrimmage.GetPtrOf("team b"),
					League:         scrimmage.BetLeague("nba"),
					Sport:          scrimmage.BetSport("basketball"),
				},
			},
		},
	)

	if err != nil {
		log.Fatalf("failed to track reward, reason : %s", err)
	}

	log.Println(trackResultOnce)
}
