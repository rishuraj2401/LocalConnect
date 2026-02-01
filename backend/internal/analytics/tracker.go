package analytics

import (
	"context"
	"sync"
	"time"
)

// Event represents an analytics event
type Event struct {
	Type      string
	UserID    string
	ProfileID string
	Metadata  map[string]interface{}
	Timestamp time.Time
}

// Tracker collects and aggregates analytics events using channels
type Tracker struct {
	events chan Event
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	stats  *Statistics
}

// Statistics holds aggregated analytics data
type Statistics struct {
	mu              sync.RWMutex
	profileViews    map[string]int64
	categoryViews   map[string]int64
	searchQueries   int64
	contactRequests int64
}

// NewTracker creates a new analytics tracker
func NewTracker(bufferSize int) *Tracker {
	ctx, cancel := context.WithCancel(context.Background())
	t := &Tracker{
		events: make(chan Event, bufferSize),
		ctx:    ctx,
		cancel: cancel,
		stats: &Statistics{
			profileViews:  make(map[string]int64),
			categoryViews: make(map[string]int64),
		},
	}
	
	// Start multiple aggregator workers
	workers := 2
	for i := 0; i < workers; i++ {
		t.wg.Add(1)
		go t.aggregator()
	}
	
	// Start periodic stats flusher
	t.wg.Add(1)
	go t.statsFlusher()
	
	return t
}

// Track queues an event for tracking
func (t *Tracker) Track(event Event) {
	event.Timestamp = time.Now()
	select {
	case t.events <- event:
	case <-t.ctx.Done():
	default:
		// Buffer full, drop event (in production, you might want to log this)
	}
}

// aggregator processes events and updates statistics
func (t *Tracker) aggregator() {
	defer t.wg.Done()
	
	for {
		select {
		case event := <-t.events:
			t.processEvent(event)
		case <-t.ctx.Done():
			// Drain remaining events
			for {
				select {
				case event := <-t.events:
					t.processEvent(event)
				default:
					return
				}
			}
		}
	}
}

func (t *Tracker) processEvent(event Event) {
	t.stats.mu.Lock()
	defer t.stats.mu.Unlock()
	
	switch event.Type {
	case "profile_view":
		t.stats.profileViews[event.ProfileID]++
	case "category_search":
		if category, ok := event.Metadata["category"].(string); ok {
			t.stats.categoryViews[category]++
		}
		t.stats.searchQueries++
	case "contact_request":
		t.stats.contactRequests++
	}
}

// statsFlusher periodically flushes stats to database
func (t *Tracker) statsFlusher() {
	defer t.wg.Done()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			t.flushStats()
		case <-t.ctx.Done():
			t.flushStats()
			return
		}
	}
}

func (t *Tracker) flushStats() {
	t.stats.mu.Lock()
	defer t.stats.mu.Unlock()
	
	// In a real application, this would write to database
	// For now, we'll just reset counters
	
	// You could persist these stats to PostgreSQL here
	t.stats.profileViews = make(map[string]int64)
	t.stats.categoryViews = make(map[string]int64)
	t.stats.searchQueries = 0
	t.stats.contactRequests = 0
}

// GetProfileViews returns the view count for a profile
func (t *Tracker) GetProfileViews(profileID string) int64 {
	t.stats.mu.RLock()
	defer t.stats.mu.RUnlock()
	return t.stats.profileViews[profileID]
}

// Shutdown gracefully shuts down the tracker
func (t *Tracker) Shutdown() {
	t.cancel()
	close(t.events)
	t.wg.Wait()
}
