package distributed

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Payload   map[string]any    `json:"payload"`
	Priority  int               `json:"priority"`
	Timeout   time.Duration     `json:"timeout"`
	MaxRetry  int               `json:"max_retry"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

type TaskResult struct {
	TaskID    string         `json:"task_id"`
	Output    map[string]any `json:"output"`
	Error     string         `json:"error,omitempty"`
	Worker    string         `json:"worker"`
	Latency   time.Duration  `json:"latency"`
	Attempt   int            `json:"attempt"`
	Retries   int            `json:"retries"`
	Succeeded bool           `json:"succeeded"`
}

type Worker interface {
	ID() string
	Execute(ctx context.Context, task *Task) (*TaskResult, error)
}

type Coordinator struct {
	workers     []Worker
	taskCh      chan *Task
	resultCh    chan *TaskResult
	mu          sync.RWMutex
	wg          sync.WaitGroup
	handlers    map[string]func(ctx context.Context, task *Task) (*TaskResult, error)
	metrics     *CoordinatorMetrics
	draining    atomic.Bool
	drainWg     sync.WaitGroup
}

type CoordinatorMetrics struct {
	TasksSubmitted atomic.Int64
	TasksCompleted atomic.Int64
	TasksFailed    atomic.Int64
	TasksRetried   atomic.Int64
	TasksTimeout   atomic.Int64
	TotalLatency   atomic.Int64
}

func NewCoordinator(workers []Worker, queueSize int) *Coordinator {
	if queueSize <= 0 {
		queueSize = 100
	}
	return &Coordinator{
		workers:  workers,
		taskCh:   make(chan *Task, queueSize),
		resultCh: make(chan *TaskResult, queueSize),
		handlers: make(map[string]func(ctx context.Context, task *Task) (*TaskResult, error)),
		metrics:  &CoordinatorMetrics{},
	}
}

func (c *Coordinator) Submit(task *Task) {
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	if task.MaxRetry < 0 {
		task.MaxRetry = 0
	}
	c.metrics.TasksSubmitted.Add(1)
	c.taskCh <- task
}

func (c *Coordinator) SubmitPriority(task *Task) {
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	c.metrics.TasksSubmitted.Add(1)
	c.taskCh <- task
}

func (c *Coordinator) RegisterHandler(taskType string, handler func(ctx context.Context, task *Task) (*TaskResult, error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[taskType] = handler
}

func (c *Coordinator) Start(ctx context.Context) {
	for _, w := range c.workers {
		c.wg.Add(1)
		go c.workerLoop(ctx, w)
	}
}

func (c *Coordinator) workerLoop(ctx context.Context, w Worker) {
	defer c.wg.Done()
	for {
		if c.draining.Load() {
			return
		}
		select {
		case <-ctx.Done():
			return
		case task, ok := <-c.taskCh:
			if !ok {
				return
			}
			c.drainWg.Add(1)
			result := c.executeWithRetry(ctx, w, task)
			c.drainWg.Done()
			if result != nil {
				c.resultCh <- result
			}
		}
	}
}

func (c *Coordinator) executeWithRetry(ctx context.Context, w Worker, task *Task) *TaskResult {
	maxAttempts := task.MaxRetry + 1
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if c.draining.Load() {
			return &TaskResult{
				TaskID:  task.ID,
				Error:   "coordinator draining",
				Worker:  w.ID(),
				Attempt: attempt,
			}
		}
		var execCtx context.Context
		var cancel context.CancelFunc
		if task.Timeout > 0 {
			execCtx, cancel = context.WithTimeout(ctx, task.Timeout)
		} else {
			execCtx, cancel = context.WithCancel(ctx)
		}
		start := time.Now()
		result, err := w.Execute(execCtx, task)
		latency := time.Since(start)
		cancel()
		if err != nil {
			lastErr = err
			c.metrics.TasksRetried.Add(1)
			if attempt < maxAttempts {
				backoff := computeBackoff(attempt)
				select {
				case <-ctx.Done():
					return &TaskResult{
						TaskID:  task.ID,
						Error:   ctx.Err().Error(),
						Worker:  w.ID(),
						Latency: latency,
						Attempt: attempt,
						Retries: attempt - 1,
					}
				case <-time.After(backoff):
					continue
				}
			}
			c.metrics.TasksFailed.Add(1)
			return &TaskResult{
				TaskID:    task.ID,
				Error:     err.Error(),
				Worker:    w.ID(),
				Latency:   latency,
				Attempt:   attempt,
				Retries:   attempt - 1,
				Succeeded: false,
			}
		}
		if result == nil {
			result = &TaskResult{
				TaskID: task.ID,
			}
		}
		result.Latency = latency
		result.Worker = w.ID()
		result.Attempt = attempt
		result.Retries = attempt - 1
		result.Succeeded = true
		c.metrics.TasksCompleted.Add(1)
		c.metrics.TotalLatency.Add(int64(latency))
		return result
	}
	c.metrics.TasksFailed.Add(1)
	return &TaskResult{
		TaskID:  task.ID,
		Error:   fmt.Sprintf("all %d attempts failed: %v", maxAttempts, lastErr),
		Worker:  w.ID(),
		Attempt: maxAttempts,
		Retries: maxAttempts - 1,
	}
}

func (c *Coordinator) Results() <-chan *TaskResult {
	return c.resultCh
}

func (c *Coordinator) Stop() {
	close(c.taskCh)
	c.wg.Wait()
	close(c.resultCh)
}

func (c *Coordinator) Drain() {
	c.draining.Store(true)
	c.drainWg.Wait()
	close(c.taskCh)
	c.wg.Wait()
	close(c.resultCh)
}

func (c *Coordinator) WorkerCount() int {
	return len(c.workers)
}

func (c *Coordinator) Metrics() CoordinatorMetricsSnapshot {
	return CoordinatorMetricsSnapshot{
		TasksSubmitted: c.metrics.TasksSubmitted.Load(),
		TasksCompleted: c.metrics.TasksCompleted.Load(),
		TasksFailed:    c.metrics.TasksFailed.Load(),
		TasksRetried:   c.metrics.TasksRetried.Load(),
		TasksTimeout:   c.metrics.TasksTimeout.Load(),
		AvgLatency:     time.Duration(c.metrics.TotalLatency.Load() / max64(c.metrics.TasksCompleted.Load(), 1)),
	}
}

type CoordinatorMetricsSnapshot struct {
	TasksSubmitted int64         `json:"tasks_submitted"`
	TasksCompleted int64         `json:"tasks_completed"`
	TasksFailed    int64         `json:"tasks_failed"`
	TasksRetried   int64         `json:"tasks_retried"`
	TasksTimeout   int64         `json:"tasks_timeout"`
	AvgLatency     time.Duration `json:"avg_latency"`
}

type SimpleWorker struct {
	workerID string
	handler  func(ctx context.Context, task *Task) (map[string]any, error)
}

func NewSimpleWorker(id string, handler func(ctx context.Context, task *Task) (map[string]any, error)) *SimpleWorker {
	return &SimpleWorker{workerID: id, handler: handler}
}

func (w *SimpleWorker) ID() string {
	return w.workerID
}

func (w *SimpleWorker) Execute(ctx context.Context, task *Task) (*TaskResult, error) {
	output, err := w.handler(ctx, task)
	if err != nil {
		return nil, err
	}
	return &TaskResult{
		TaskID: task.ID,
		Output: output,
	}, nil
}

type LoadBalancer struct {
	workers []Worker
	counter uint64
	mu      sync.Mutex
}

func NewLoadBalancer(workers []Worker) *LoadBalancer {
	return &LoadBalancer{workers: workers}
}

func (lb *LoadBalancer) Next() Worker {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if len(lb.workers) == 0 {
		return nil
	}
	w := lb.workers[lb.counter%uint64(len(lb.workers))]
	lb.counter++
	return w
}

func (lb *LoadBalancer) WorkerCount() int {
	return len(lb.workers)
}

type ResultAggregator struct {
	results []TaskResult
	mu      sync.Mutex
}

func NewResultAggregator() *ResultAggregator {
	return &ResultAggregator{}
}

func (ra *ResultAggregator) Add(result TaskResult) {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	ra.results = append(ra.results, result)
}

func (ra *ResultAggregator) All() []TaskResult {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	out := make([]TaskResult, len(ra.results))
	copy(out, ra.results)
	return out
}

func (ra *ResultAggregator) Succeeded() []TaskResult {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	var out []TaskResult
	for _, r := range ra.results {
		if r.Succeeded {
			out = append(out, r)
		}
	}
	return out
}

func (ra *ResultAggregator) Failed() []TaskResult {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	var out []TaskResult
	for _, r := range ra.results {
		if !r.Succeeded || r.Error != "" {
			out = append(out, r)
		}
	}
	return out
}

func (ra *ResultAggregator) Count() int {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	return len(ra.results)
}

func (ra *ResultAggregator) Summary() string {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	failed := 0
	for _, r := range ra.results {
		if !r.Succeeded || r.Error != "" {
			failed++
		}
	}
	return fmt.Sprintf("%d total, %d succeeded, %d failed", len(ra.results), len(ra.results)-failed, failed)
}

type CircuitBreaker struct {
	failureCount int64
	successCount int64
	state        int32
	threshold    int64
	resetTimeout time.Duration
	lastFailure  time.Time
	mu           sync.Mutex
}

const (
	cbClosed int32 = iota
	cbOpen
	cbHalfOpen
)

func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold:    int64(threshold),
		resetTimeout: resetTimeout,
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case cbClosed:
		return true
	case cbOpen:
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.state = cbHalfOpen
			return true
		}
		return false
	case cbHalfOpen:
		return true
	}
	return true
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failureCount = 0
	cb.successCount++
	cb.state = cbClosed
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failureCount++
	cb.lastFailure = time.Now()
	if cb.failureCount >= cb.threshold {
		cb.state = cbOpen
	}
}

func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case cbClosed:
		return "closed"
	case cbOpen:
		return "open"
	case cbHalfOpen:
		return "half-open"
	}
	return "unknown"
}

type CircuitBreakerWorker struct {
	worker Worker
	cb     *CircuitBreaker
}

func NewCircuitBreakerWorker(worker Worker, cb *CircuitBreaker) *CircuitBreakerWorker {
	return &CircuitBreakerWorker{worker: worker, cb: cb}
}

func (w *CircuitBreakerWorker) ID() string {
	return w.worker.ID()
}

func (w *CircuitBreakerWorker) Execute(ctx context.Context, task *Task) (*TaskResult, error) {
	if !w.cb.Allow() {
		return nil, fmt.Errorf("circuit breaker open for worker %s", w.worker.ID())
	}
	result, err := w.worker.Execute(ctx, task)
	if err != nil {
		w.cb.RecordFailure()
		return nil, err
	}
	w.cb.RecordSuccess()
	return result, nil
}

type WorkerHealth struct {
	WorkerID      string        `json:"worker_id"`
	Healthy       bool          `json:"healthy"`
	LastHeartbeat time.Time     `json:"last_heartbeat"`
	ResponseTime  time.Duration `json:"response_time"`
	TasksExecuted int64         `json:"tasks_executed"`
}

type HealthChecker struct {
	workers  []Worker
	interval time.Duration
	health   map[string]*WorkerHealth
	mu       sync.RWMutex
	stopCh   chan struct{}
}

func NewHealthChecker(workers []Worker, interval time.Duration) *HealthChecker {
	hc := &HealthChecker{
		workers:  workers,
		interval: interval,
		health:   make(map[string]*WorkerHealth),
		stopCh:   make(chan struct{}),
	}
	for _, w := range workers {
		hc.health[w.ID()] = &WorkerHealth{
			WorkerID: w.ID(),
			Healthy:  true,
		}
	}
	return hc
}

func (hc *HealthChecker) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(hc.interval)
		defer ticker.Stop()
		for {
			select {
			case <-hc.stopCh:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				hc.checkAll(ctx)
			}
		}
	}()
}

func (hc *HealthChecker) checkAll(ctx context.Context) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	for _, w := range hc.workers {
		h := hc.health[w.ID()]
		start := time.Now()
		_, err := w.Execute(ctx, &Task{
			ID:   "_health_check",
			Type: "_health",
		})
		h.ResponseTime = time.Since(start)
		h.LastHeartbeat = time.Now()
		h.Healthy = err == nil
	}
}

func (hc *HealthChecker) Stop() {
	close(hc.stopCh)
}

func (hc *HealthChecker) Status() map[string]*WorkerHealth {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	out := make(map[string]*WorkerHealth)
	for k, v := range hc.health {
		copy := *v
		out[k] = &copy
	}
	return out
}

func (hc *HealthChecker) HealthyWorkers() []Worker {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	var healthy []Worker
	for _, w := range hc.workers {
		if h, ok := hc.health[w.ID()]; ok && h.Healthy {
			healthy = append(healthy, w)
		}
	}
	return healthy
}

func computeBackoff(attempt int) time.Duration {
	base := time.Duration(1<<uint(attempt-1)) * 100 * time.Millisecond
	if base > 5*time.Second {
		base = 5 * time.Second
	}
	return base
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}


