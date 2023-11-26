package sender

import (
	"XAXAtonMTC/pkg/packetsender"
	"log"
	"net/http"
	"sync"
)

type SendersMemoryRepository struct {
	mu      *sync.RWMutex
	senders map[uint32]*Sender
}

func NewSendersMemoryRepository() *SendersMemoryRepository {
	return &SendersMemoryRepository{
		mu:      &sync.RWMutex{},
		senders: make(map[uint32]*Sender),
	}
}

func (repo *SendersMemoryRepository) AddSender(userID uint32, w http.ResponseWriter, r *http.Request) error {
	s, err := NewSender(w, r)
	if err != nil {
		return err
	}

	repo.mu.Lock()
	repo.senders[userID] = s
	repo.mu.Unlock()

	return nil
}

func (repo *SendersMemoryRepository) SendPacket(usersID []uint32, check func(uint32) bool, packet *packetsender.Packet) {
	wg := sync.WaitGroup{}

	repo.mu.RLock()
	defer repo.mu.RUnlock()
	for _, userID := range usersID {
		wg.Add(1)
		go func(userID uint32) {
			defer wg.Done()

			if !check(userID) {
				return
			}

			err := repo.senders[userID].SendPacket(packet)
			if err != nil {
				log.Println(err)
			}
		}(userID)
	}

	wg.Wait()
}

func (repo *SendersMemoryRepository) DeleteSender(userID uint32) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	sender, ok := repo.senders[userID]
	if !ok {
		return
	}

	sender.wsConn.Close()
	delete(repo.senders, userID)
}

func (repo *SendersMemoryRepository) CloseConnections(usersID ...uint32) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, userID := range usersID {
		sender, ok := repo.senders[userID]
		if !ok {
			continue
		}

		sender.wsConn.Close()
		delete(repo.senders, userID)
	}
}
