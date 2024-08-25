package blockchain

type (
	Blockchain struct {
		Id             int64  `json:"id"`
		Name           string `json:"name"`
		NativeCurrency string `json:"native_currency"`
		ChainId        int64  `json:"chain_id"`
	}
)
