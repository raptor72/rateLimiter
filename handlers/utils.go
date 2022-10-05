package handlers

import (
	"net"

	bl "github.com/raptor72/rateLimiter/api/blacklists"
	wl "github.com/raptor72/rateLimiter/api/whitelists"
	"github.com/raptor72/rateLimiter/config"
	log "github.com/sirupsen/logrus"
)

func InWhiteList(cfg *config.Config, ip net.IP) (bool, error) {
	db, err := cfg.NewDB()
	if err != nil {
		log.WithError(err).Error("new DB error")
		return false, err
	}
	wlStorage := wl.NewPgsqlStorage(db)
	whitelists, err := wlStorage.Select()
	if err != nil {
		log.WithError(err).Error("failed to select white lists")
		return false, err
	}

	for _, cidr := range whitelists {
		ipv4Addr, ipv4Net, err := net.ParseCIDR(cidr.Address)
		if err != nil {
			log.WithError(err).Errorf("got error during parsing exists white list %v", err)
			continue
		}
		if ipv4Net.Contains(ip) || ipv4Addr.Equal(ip) {
			return true, nil
		}
	}
	return false, nil
}

func InBlackList(cfg *config.Config, ip net.IP) (bool, error) {
	db, err := cfg.NewDB()
	if err != nil {
		log.WithError(err).Error("new DB error")
		return false, err
	}
	blStorage := bl.NewPgsqlStorage(db)
	blacklists, err := blStorage.Select()
	if err != nil {
		log.WithError(err).Error("failed to select black lists")
		return false, err
	}

	for _, cidr := range blacklists {
		ipv4Addr, ipv4Net, err := net.ParseCIDR(cidr.Address)
		if err != nil {
			log.WithError(err).Errorf("got error during parsing exists black list %v", err)
			continue
		}
		if ipv4Net.Contains(ip) || ipv4Addr.Equal(ip) {
			return true, nil
		}
	}
	return false, nil
}
