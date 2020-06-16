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
		pm.mu.Lock()
		fmt.Println("BBB4--", "RLock")
		pm.dialoutBlacklist[ip] = time.Now().Unix() + dialoutBlacklistCooldownSec
		log.Info("add black list", "len", len(pm.dialoutBlacklist), "ip", ip, "data", pm.dialoutBlacklist)
		pm.mu.Unlock()
		fmt.Println("BBB4--", "RLock")
	}
}

func (pm *blackNodes) ChkDialoutBlacklist(ip string) bool {
	pm.mu.RLock()
	fmt.Println("BBB--2222", "RLock")
	tm, ok := pm.dialoutBlacklist[ip]
	pm.mu.RUnlock()
	fmt.Println("BBB--4444", "RUnLock")
	if ok {
		if time.Now().Unix() < tm {
			return true
		}
		pm.mu.Lock()
		fmt.Println("BBB--666", "Lock")
		delete(pm.dialoutBlacklist, ip)
		pm.mu.Unlock()
		fmt.Println("BBB--9999", "UnLock")
	}
	return false
}

func (pm *blackNodes) chkDialinBlacklist(ip string) bool {
	pm.mu.RLock()
	fmt.Println("BBB1--", "RLock")
	tm, ok := pm.dialinBlacklist[ip]
	pm.mu.RUnlock()
	fmt.Println("BBB1-", "RLock")
	if ok {
		if time.Now().Unix() < tm {
			return true
		}
		pm.mu.Lock()
		fmt.Println("BBB2--", "RLock")
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
	b.mu.Lock()
	fmt.Println("BBB3--", "Lock")
	defer b.mu.Unlock()
	defer fmt.Println("BBB3--", "Unlock")
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
