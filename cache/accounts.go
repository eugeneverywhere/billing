package cache

import (
	"github.com/eugeneverywhere/billing/db"
	"github.com/eugeneverywhere/billing/db/models"
	"github.com/lillilli/logger"
	"sync"
)

// AccountsCache - accounts cache interface
type AccountsCache interface {
	GetAccountsByID() map[int64]*models.Account
	GetAccountsByExtID() map[string]*models.Account
	Update()
}

type accountsCache struct {
	log          logger.Logger
	db           db.DB
	cacheByID    map[int64]*models.Account
	cacheByExtID map[string]*models.Account
	sync.Mutex
}

// NewAccountsCache - return new cache instance
func NewAccountsCache(db db.DB) AccountsCache {
	return &accountsCache{
		db:           db,
		cacheByID:    make(map[int64]*models.Account),
		cacheByExtID: make(map[string]*models.Account),
		log:          logger.NewLogger("accounts cache"),
	}
}

func (c *accountsCache) GetAccountsByID() map[int64]*models.Account {
	c.Lock()

	defer c.Unlock()
	return c.cacheByID
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
	mapByID, mapByExtID := makeMaps(accounts)
	c.Lock()
	c.cacheByID = mapByID
	c.cacheByExtID = mapByExtID
	c.Unlock()
}

func makeMaps(array []*models.Account) (byID map[int64]*models.Account, byExtID map[string]*models.Account) {
	byID = make(map[int64]*models.Account)
	byExtID = make(map[string]*models.Account)
	for _, v := range array {
		byID[v.ID] = v
		byExtID[v.ExternalID] = v
	}
	return byID, byExtID
}
