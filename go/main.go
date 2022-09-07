package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BattleHistoryCurrent struct {
	CurrentBattle   *CurrentBattle         `json:"current_battle"`
	PreviousBattles []*BattleHistoryRecord `json:"previous_battles"`
}
type CurrentBattle struct {
	Number    int    `json:"number"`
	StartedAt int64  `json:"started_at"`
	ExpiresAt int64  `json:"expires_at"`
	Signature string `json:"signature"`
}
type BattleHistoryRecord struct {
	Number    int    `json:"number"`
	StartedAt int64  `json:"started_at"`
	EndedAt   *int64 `json:"ended_at"`
	Winner    int64  `json:"winner"`
	RunnerUp  int64  `json:"runner_up"`
	Loser     int64  `json:"loser"`
	Signature string `json:"signature"`
}

func main() {
	resp, err := http.Get("http://api.supremacygame.dev/api/battle_history")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	result := &BattleHistoryCurrent{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", result.CurrentBattle)
	for _, battle := range result.PreviousBattles {
		fmt.Printf("%+v\n", battle)
	}
}
