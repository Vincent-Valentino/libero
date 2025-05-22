package scheduler

import (
	"libero-backend/internal/service"
	"log"
	"time"
)

// Scheduler handles periodic tasks for the application
type Scheduler struct {
	fixturesService service.FixturesService
	stopChan        chan struct{}
	running         bool
}

// New creates a new scheduler with required services
func New(fixturesService service.FixturesService) *Scheduler {
	return &Scheduler{
		fixturesService: fixturesService,
		stopChan:        make(chan struct{}),
	}
}

// Start begins the scheduled tasks
func (s *Scheduler) Start() {
	if s.running {
		return
	}
	s.running = true
	go s.run()
	log.Println("Scheduler started")
}

// Stop terminates the scheduled tasks
func (s *Scheduler) Stop() {
	if !s.running {
		return
	}
	s.running = false
	s.stopChan <- struct{}{}
	log.Println("Scheduler stopped")
}

// run is the main scheduler loop
func (s *Scheduler) run() {
	// Refresh fixtures data every 4 hours
	refreshTicker := time.NewTicker(4 * time.Hour)
	defer refreshTicker.Stop()

	// Immediately fetch data on startup
	s.refreshFixtures()

	for {
		select {
		case <-refreshTicker.C:
			s.refreshFixtures()
		case <-s.stopChan:
			return
		}
	}
}

// refreshFixtures updates fixtures data via the fixtures service
func (s *Scheduler) refreshFixtures() {
	log.Println("Scheduler: Refreshing fixtures data")
	// This is a placeholder that would normally call service methods
	// Since this is a minimal implementation, we're just logging
	// In a real implementation, this would call methods like:
	// _, err := s.fixturesService.GetTodaysFixtures()
	// if err != nil {
	//     log.Printf("Error refreshing fixtures: %v", err)
	// }
} 