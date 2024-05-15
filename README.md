# Official Scrimmage SDK for Go

[![go.dev](https://img.shields.io/badge/go.dev-pkg-007d9c.svg?style=flat)](https://pkg.go.dev/github.com/Scrimmage-co/golang-sdk)

This library is a part of the [Scrimmage Rewards Program](https://scrimmage.co)
that is providing a solution for loyalty programs and rewards.


## Requirements
This SDK requires a minimum version of 1.18.

## Installation

`scrimmage` can be installed like any other Go library through `go get`:

```console
$ go get github.com/Scrimmage-co/golang-sdk
```

## Usage

### 1. Initialize the SDK:
To initialize the library, provide your unique server endpoint and secret key. Obtain the `API_SERVER_ENDPOINT` from your admin dashboard URL, which should resemble "your_company_name.apps.scrimmage.co". The secret key is generated during [Step 2 of Getting Started.](https://docs.scrimmage.co/docs/getting-started#2-create-secret-key)


   ```go
   import scrimmage "github.com/Scrimmage-co/golang-sdk"
   ...
    sdk, err := scrimmage.InitRewarder(
		context.Background(),
		apiServerEndpoint,
		privateKey,
		namespace
    )
   ...
   ```
### 2. User Authentication

Scrimmage uses an anonymous user ID for each user to maintain their reward profile without storing login credentials. Authenticate users every time they access the reward widget or iframe:

```go
...
sdk.User.GetUserToken(context.Background, scrimmage. GetUserTokenRequest{
    UserID:     "userId",
    Tags:       ["tag-1", "tag-2"], // optional
    Properties: map[string]interface{
        "at": 1718231233,
        "device": "personal"
    } // optional
    })
...
```

Pass this token to the frontend for user authentication. The token is valid for 24 hours.


```
https://<API_SERVER_ENDPOINT>?token=<TOKEN>
```

### 3. Sending an Event

To track user activities, use the following methods. Customize the interface as needed:

```go
result, err := sdk.Reward.TrackRewardableOnce(
    context.Background(),
    "userId",
    scrimmage.BetDataType_BetExecuted,
    scrimmage.GetPtrOf("uniqueEventId"),
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
```

Please insert this code wherever your bets (events) are executed. Once this code is inserted, it will open up a one-way connection for you to send bet details to Scrimmage.


## License

Licensed under The MIT License, see [LICENSE](https://github.com/Scrimmage-co/golang-sdk/blob/main/LICENSE).