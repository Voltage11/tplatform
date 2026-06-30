package cache

import (
	"context"
	"sync"
	"time"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/google/uuid"
)

type UserCache interface {
	GetByID(id uuid.UUID) *domain.User
	GetByEmail(email string) *domain.User
	Set(user *domain.User)
	Delete(id uuid.UUID)
	ClearExpired()
	Stop()
}

type userItem struct {
	*domain.User
	expiredAt time.Time
}

type userCache struct {
	mu         sync.RWMutex
	ttl        time.Duration
	users      map[uuid.UUID]*userItem
	usersEmail map[string]*userItem
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewUserCache(ttl time.Duration) UserCache {
	ctx, cancel := context.WithCancel(context.Background())

	c := &userCache{
		ttl:        ttl,
		users:      make(map[uuid.UUID]*userItem),
		usersEmail: make(map[string]*userItem),
		ctx:        ctx,
		cancel:     cancel,
	}

	// Запускаем горутину очистки
	go c.cleanupLoop()

	return c
}

func (c *userCache) cleanupLoop() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.ClearExpired()
		}
	}
}

func (c *userCache) GetByID(id uuid.UUID) *domain.User {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.users[id]
	if !ok {
		return nil
	}

	if time.Now().UTC().Before(item.expiredAt) {
		return item.User
	}

	// Истекший кэш удаляем при следующем вызове ClearExpired
	return nil
}

func (c *userCache) GetByEmail(email string) *domain.User {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.usersEmail[email]
	if !ok {
		return nil
	}

	if time.Now().UTC().Before(item.expiredAt) {
		return item.User
	}

	return nil
}

func (c *userCache) Set(user *domain.User) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если пользователь уже есть с другим email, удаляем старую запись в usersEmail
	if oldItem, ok := c.users[user.ID]; ok && oldItem.Email != user.Email {
		delete(c.usersEmail, oldItem.Email)
	}

	item := &userItem{
		User:      user,
		expiredAt: time.Now().UTC().Add(c.ttl),
	}

	c.users[user.ID] = item
	c.usersEmail[user.Email] = item
}

func (c *userCache) Delete(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.users[id]; ok {
		delete(c.users, id)
		delete(c.usersEmail, item.Email)
	}
}

func (c *userCache) ClearExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UTC()
	for id, item := range c.users {
		if now.After(item.expiredAt) {
			delete(c.users, id)
			delete(c.usersEmail, item.Email)
		}
	}
}

func (c *userCache) Stop() {
	c.cancel()
}
