package transaction_service

func getCardBalanceInPan(pan int) float32 {
	d := 100000000

	if pan < d {
		return 0
	}
	res := pan % d

	return float32(res)
}
