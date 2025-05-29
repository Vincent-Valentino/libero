package scheduler

import (
	"context"
	"libero-backend/internal/service"
	"log"
	"time"
)

// Scheduler manages periodic background tasks.
type Scheduler struct {
	fixturesService service.FixturesService
	ctx             context.Context
	cancel          context.CancelFunc
}

// New creates a new scheduler.
func New(fixturesService service.FixturesService) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		fixturesService: fixturesService,
		ctx:             ctx,
		cancel:          cancel,
	}
}

// Start begins all scheduled tasks.
func (s *Scheduler) Start() {
	// Start the task to refresh today's fixtures every 4 hours
	go s.scheduleTodayFixtures()

	// Start the task to refresh fixtures summary for all major competitions
	go s.scheduleFixturesSummaries()
}

// Stop terminates all scheduled tasks.
func (s *Scheduler) Stop() {
	s.cancel()
}

// scheduleTodayFixtures refreshes today's fixtures every 4 hours.
func (s *Scheduler) scheduleTodayFixtures() {
	// First run immediately
	s.fetchTodayFixtures()

	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.fetchTodayFixtures()
		case <-s.ctx.Done():
			log.Println("Today's fixtures scheduler stopped")
			return
		}
	}
}

// scheduleFixturesSummaries refreshes fixtures summaries for major competitions every 6 hours.
func (s *Scheduler) scheduleFixturesSummaries() {
	// Major competition codes
	comps := []string{"PL", "PD", "SA", "BL1", "FL1", "CL", "EL"}

	// Wait 15 seconds before starting to avoid overwhelming the API on startup
	time.Sleep(15 * time.Second)

	// First run immediately
	for _, comp := range comps {
		s.fetchFixturesSummary(comp)
		// Wait between requests
		time.Sleep(5 * time.Second)
	}

	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, comp := range comps {
				s.fetchFixturesSummary(comp)
				// Wait between requests to not hit rate limits
				time.Sleep(5 * time.Second)
			}
		case <-s.ctx.Done():
			log.Println("Fixtures summaries scheduler stopped")
			return
		}
	}
}

// fetchTodayFixtures gets today's fixtures and logs any errors.
func (s *Scheduler) fetchTodayFixtures() {
	log.Println("Scheduler: Refreshing today's fixtures")
	_, err := s.fixturesService.GetTodaysFixtures()
	if err != nil {
		log.Printf("Scheduler: Error refreshing today's fixtures: %v", err)
	} else {
		log.Println("Scheduler: Today's fixtures refreshed successfully")
	}
}

// fetchFixturesSummary gets fixtures summary for a competition and logs any errors.
func (s *Scheduler) fetchFixturesSummary(competitionCode string) {
	log.Printf("Scheduler: Refreshing fixtures summary for %s", competitionCode)
	_, err := s.fixturesService.GetFixturesSummary(competitionCode)
	if err != nil {
		log.Printf("Scheduler: Error refreshing fixtures summary for %s: %v", competitionCode, err)
	} else {
		log.Printf("Scheduler: Fixtures summary for %s refreshed successfully", competitionCode)
	}
}
