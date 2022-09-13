package sstat

import (
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

func GetMvdsvServer(address string) mvdsv.Mvdsv {
	server, err := serverstat.NewClient().GetInfo(address)

	if err != nil {
		return mvdsv.Mvdsv{}
	}

	return convert.ToMvdsv(server)
}
