package main

import (
	"fmt"
	"log"
	"os"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

const (
	apiKey    string = "API Key"
	apiSecret string = "API Secret"
)

var (
	/*
		instrumentsIDs = []uint32{3861249, 40193, 60417, 1510401, 4267265, 81153, 136442372, 4268801,
			134657, 2714625, 140033, 177665, 5215745, 2800641, 225537, 232961, 315393, 1850625, 341249,
			119553, 345089, 348929, 356865, 340481, 1270529, 424961, 1346049, 408065, 3001089, 492033,
			2939649, 519937, 2815745, 2977281, 4598529, 633601, 3834113, 738561, 5582849, 794369, 779521,
			857857, 2953217, 878593, 884737, 895745, 3465729, 897537, 2889473, 2952193, 969473}
	*/
	instrumentsIDs = []uint32{738561}
)

func calcEMA(prevEMA, currValue, period float64) float64 {
	currEMA := (currValue * (2.0 / (period + 1))) + prevEMA*(1.0-(2.0/(period+1)))
	// fmt.Println(currValue, ":", currEMA)
	return currEMA
}

func main() {

	// Get margins
	from := time.Date(2022, 6, 3, 10, 0, 0, 0, time.Now().Location())
	to := time.Date(2022, 6, 3, 15, 30, 0, 0, time.Now().Location())

	// Create a new Kite connect instance
	kc := kiteconnect.New(apiKey)

	// Get user details and access token
	data, err := kc.GenerateSession("Secret Token", apiSecret)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	// Set access token
	kc.SetAccessToken(data.AccessToken)
	log.Println("data.AccessToken", data.AccessToken)

	var prevEMA float64
	var emaPeriods = []float64{5.0, 10.0, 20.0, 40.0, 80.0}

	for _, instID := range instrumentsIDs {
		historicalData, err := kc.GetHistoricalData(int(instID), "minute", from, to, false, false)
		if err != nil {
			fmt.Printf("Error getting Historical Data: %v\n", err)
			os.Exit(1)
		}

		ema5Data := make([]float64, len(historicalData))
		ema10Data := make([]float64, len(historicalData))
		ema20Data := make([]float64, len(historicalData))
		ema40Data := make([]float64, len(historicalData))
		ema80Data := make([]float64, len(historicalData))

		emaDataMap := map[float64][]float64{5.0: ema5Data, 10.0: ema10Data, 20.0: ema20Data, 40.0: ema40Data, 80.0: ema80Data}

		for _, period := range emaPeriods {
			prevEMA = 0.0
			for i, data := range historicalData {
				if i < int(period) {
					prevEMA += data.Close
					continue
				}

				if i == int(period) {
					prevEMA /= period
					for j := 0; j < i; j++ {
						emaDataMap[period][j] = prevEMA
					}
				}

				prevEMA = calcEMA(prevEMA, data.Close, period)
				emaDataMap[period][i] = prevEMA
			}
			fmt.Println("Instrument ID: ", instID, " Period: ", period, " EMA: ", emaDataMap[period][len(historicalData)-1])
		}
	}
}
