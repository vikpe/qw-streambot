package sstat

import (
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

func GetMvdsvServer(address string) mvdsv.Mvdsv {
	nullResult := mvdsv.Mvdsv{}

	if "" == address {
		return nullResult
	}

	genericServer, err := serverstat.GetInfo(address)

	if err != nil {
		return nullResult
	}

	return convert.ToMvdsv(genericServer)
}
