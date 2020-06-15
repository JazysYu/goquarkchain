package nodefilter

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

const (
	// period between adding a dialout_blacklisted node and removing it
	dialoutBlacklistCooldownSec int64 = 24 * 3600
	dialinBlacklistCooldownSec  int64 = 8 * 3600
	unblacklistInterval               = time.Duration(15 * 60)
)

type BlackFilter interface {
	AddDialoutBlacklist(string)
	ChkDialoutBlacklist(string) bool
	PeriodicallyUnblacklist()
}

func NewHandleBlackListErr(text string) error {
	return &BlackErr{text}
}

// errorString is a trivial implementation of error.
type BlackErr struct {
	s string
}

func (e *BlackErr) Error() string {
	return e.s
}

type blackNodes struct {
	currTime time.Time
	// BootstrapNodes | preferedNodes
	WhitelistNodes map[string]*enode.Node
	// IP to unblacklist time, we blacklist by IP
	mu               sync.RWMutex
	dialoutBlacklist map[string]int64
	dialinBlacklist  map[string]int64
}

func NewBlackList(whitelistNodes map[string]*enode.Node) BlackFilter {
	return &blackNodes{
		currTime:         time.Now(),
		WhitelistNodes:   whitelistNodes,
		dialoutBlacklist: make(map[string]int64),
		dialinBlacklist:  make(map[string]int64),
	}
}

func (pm *blackNodes) AddDialoutBlacklist(ip string) {
	if _, ok := pm.WhitelistNodes[ip]; !ok {
		fmt.Println("BBB4--", "RLock")
		pm.mu.Lock()
		pm.dialoutBlacklist[ip] = time.Now().Unix() + dialoutBlacklistCooldownSec
		log.Info("add black list", "len", len(pm.dialoutBlacklist), "ip", ip, "data", pm.dialoutBlacklist)
		pm.mu.Unlock()
		fmt.Println("BBB4--", "RLock")
	}
}

func (pm *blackNodes) ChkDialoutBlacklist(ip string) bool {
	fmt.Println("BBB--1111", "RLock")
	pm.mu.RLock()
	fmt.Println("BBB--2222")
	tm, ok := pm.dialoutBlacklist[ip]
	fmt.Println("BBB--333")
	pm.mu.RUnlock()
	fmt.Println("BBB--4444", "RUnLock")
	if ok {
		if time.Now().Unix() < tm {
			fmt.Println("BBB--555")
			return true
		}
		fmt.Println("BBB--666", "Lock")
		pm.mu.Lock()
		fmt.Println("BBB--7777")
		delete(pm.dialoutBlacklist, ip)
		fmt.Println("BBB--8888")
		pm.mu.Unlock()
		fmt.Println("BBB--9999", "UnLock")
	}
	return false
}

func (pm *blackNodes) chkDialinBlacklist(ip string) bool {
	fmt.Println("BBB1--", "RLock")
	pm.mu.RLock()
	tm, ok := pm.dialinBlacklist[ip]
	pm.mu.RUnlock()
	fmt.Println("BBB1-", "RLock")
	if ok {
		if time.Now().Unix() < tm {
			return true
		}
		fmt.Println("BBB2--", "RLock")
		pm.mu.Lock()
		delete(pm.dialinBlacklist, ip)
		pm.mu.Unlock()
		fmt.Println("BBB2--", "RLock")
	}
	return false
}

func (b *blackNodes) PeriodicallyUnblacklist() {
	now := time.Now()
	if now.Sub(b.currTime) < unblacklistInterval {
		return
	}
	fmt.Println("BBB3--", "RLock")
	b.mu.Lock()
	defer b.mu.Unlock()
	defer fmt.Println("BBB3--", "RLock")
	b.currTime = now
	for ip, tm := range b.dialoutBlacklist {
		if now.Unix() >= tm {
			delete(b.dialoutBlacklist, ip)
		}
	}
	for ip, tm := range b.dialinBlacklist {
		if now.Unix() >= tm {
			delete(b.dialinBlacklist, ip)
		}
	}
}
