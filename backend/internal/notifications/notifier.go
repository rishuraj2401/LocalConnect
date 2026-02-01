package notifications

import (
	"context"
	"log"
	"sync"
	"time"
)

// Notification represents a notification to be sent
type Notification struct {
	UserID  string
	Type    string
	Message string
	Data    map[string]interface{}
}

// Notifier handles async notification sending with batching
type Notifier struct {
	queue      chan Notification
	batchSize  int
	batchDelay time.Duration
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewNotifier creates a new notification service
func NewNotifier(bufferSize, batchSize int, batchDelay time.Duration) *Notifier {
	ctx, cancel := context.WithCancel(context.Background())
	n := &Notifier{
		queue:      make(chan Notification, bufferSize),
		batchSize:  batchSize,
		batchDelay: batchDelay,
		ctx:        ctx,
		cancel:     cancel,
	}
	
	// Start worker goroutines for batch processing
	workers := 3
	for i := 0; i < workers; i++ {
		n.wg.Add(1)
		go n.worker(i)
	}
	
	return n
}

// Send queues a notification for async processing
func (n *Notifier) Send(notification Notification) {
	select {
	case n.queue <- notification:
	case <-n.ctx.Done():
	default:
		// Queue is full, process immediately in goroutine
		go n.processNotification(notification)
	}
}

// worker processes notifications in batches
func (n *Notifier) worker(id int) {
	defer n.wg.Done()
	
	batch := make([]Notification, 0, n.batchSize)
	ticker := time.NewTicker(n.batchDelay)
	defer ticker.Stop()
	
	processBatch := func() {
		if len(batch) > 0 {
			n.processBatch(batch)
			batch = batch[:0]
		}
	}
	
	for {
		select {
		case notification := <-n.queue:
			batch = append(batch, notification)
			if len(batch) >= n.batchSize {
				processBatch()
			}
		case <-ticker.C:
			processBatch()
		case <-n.ctx.Done():
			processBatch()
			return
		}
	}
}

// processBatch processes a batch of notifications
func (n *Notifier) processBatch(batch []Notification) {
	// Group notifications by user for efficiency
	userNotifications := make(map[string][]Notification)
	for _, notif := range batch {
		userNotifications[notif.UserID] = append(userNotifications[notif.UserID], notif)
	}
	
	// Process notifications concurrently per user
	var wg sync.WaitGroup
	for userID, notifs := range userNotifications {
		wg.Add(1)
		go func(uid string, notifications []Notification) {
			defer wg.Done()
			for _, notif := range notifications {
				n.processNotification(notif)
			}
		}(userID, notifs)
	}
	wg.Wait()
}

// processNotification sends a single notification
func (n *Notifier) processNotification(notification Notification) {
	// In a real application, this would send emails, push notifications, etc.
	// For now, we'll just log it
	log.Printf("Notification to %s [%s]: %s", notification.UserID, notification.Type, notification.Message)
	
	// Simulate some processing time
	time.Sleep(10 * time.Millisecond)
}

// Shutdown gracefully shuts down the notifier
func (n *Notifier) Shutdown() {
	n.cancel()
	close(n.queue)
	n.wg.Wait()
}
