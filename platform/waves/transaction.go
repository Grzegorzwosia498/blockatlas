package waves

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"strconv"

	"github.com/trustwallet/blockatlas/coin"
)

func (p *Platform) GetTxsByAddress(address string) (blockatlas.TxPage, error) {
	addressTxs, err := p.client.GetTxs(address, 25)
	if err != nil {
		return nil, err
	}

	txs := NormalizeTxs(addressTxs)

	return txs, nil
}

func NormalizeTxs(srcTxs []Transaction) (txs []blockatlas.Tx) {
	for _, srcTx := range srcTxs {
		tx, ok := NormalizeTx(&srcTx)
		if !ok || len(txs) >= blockatlas.TxPerPage {
			continue
		}
		txs = append(txs, tx)
	}
	return
}

func NormalizeTx(srcTx *Transaction) (tx blockatlas.Tx, ok bool) {
	var result blockatlas.Tx

	if srcTx.Type == 4 && len(srcTx.AssetId) == 0 {
		result = blockatlas.Tx{
			ID:     srcTx.Id,
			Coin:   coin.WAVES,
			From:   srcTx.Sender,
			To:     srcTx.Recipient,
			Fee:    blockatlas.Amount(strconv.Itoa(int(srcTx.Fee))),
			Date:   int64(srcTx.Timestamp) / 1000,
			Block:  srcTx.Block,
			Memo:   srcTx.Attachment,
			Status: blockatlas.StatusCompleted,
			Meta: blockatlas.Transfer{
				Value:    blockatlas.Amount(strconv.Itoa(int(srcTx.Amount))),
				Symbol:   coin.Coins[coin.WAVES].Symbol,
				Decimals: coin.Coins[coin.WAVES].Decimals,
			},
		}
		return result, true
	}

	return result, false
}
