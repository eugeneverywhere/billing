package cache

import (
	"github.com/eugeneverywhere/billing/db"
	"github.com/eugeneverywhere/billing/db/models"
	"github.com/lillilli/logger"
	"sync"
)

// AccountsCache - accounts cache interface
type AccountsCache interface {
	GetAccountsByExtID() map[string]*models.Account
	Update()
}

type accountsCache struct {
	log          logger.Logger
	db           db.DB
	cacheByExtID map[string]*models.Account
	sync.Mutex
}

// NewAccountsCache - return new cache instance
func NewAccountsCache(db db.DB) AccountsCache {
	return &accountsCache{
		db:           db,
		cacheByExtID: make(map[string]*models.Account),
		log:          logger.NewLogger("accounts cache"),
	}
}

func (c *accountsCache) GetAccountsByExtID() map[string]*models.Account {
	c.Lock()

	defer c.Unlock()
	return c.cacheByExtID
}

// Update - update accounts
func (c *accountsCache) Update() {
	accounts := make([]*models.Account, 0)
	accounts, err := c.db.GetAllAccounts()
	if err != nil {
		c.log.Errorf("update accounts cache failed: %v", err.Error())
	}
	mapByExtID := makeMaps(accounts)
	c.Lock()
	c.cacheByExtID = mapByExtID
	c.Unlock()
}

func makeMaps(array []*models.Account) (byExtID map[string]*models.Account) {
	byExtID = make(map[string]*models.Account)
	for _, v := range array {
		byExtID[v.ExternalID] = v
	}
	return byExtID
}
