package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	// Add gopsutil dependency for real metrics
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// Handler now includes dbPool
type Handler struct {
	dbPool *pgxpool.Pool
}

// NewHandler accepts dbPool
func NewHandler(dbPool *pgxpool.Pool) *Handler {
	return &Handler{
		dbPool: dbPool,
	}
}


// SystemMetrics represents CPU and memory usage metrics
type SystemMetrics struct {
	Timestamp     time.Time `json:"timestamp"`
	CPUUsage      float64   `json:"cpu_usage_percent"`
	MemoryUsageMB int       `json:"memory_usage_mb"`
	Status        string    `json:"status"`
}

// WebServerHardeningResponse is the structure for API responses
type WebServerHardeningResponse struct {
	TestID   string        `json:"test_id"`
	Status   string        `json:"status"`
	Metrics  SystemMetrics `json:"metrics"`
}

// WebServerHandler handles web server hardening requests
func (h *Handler) WebServerHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a test ID based on current time
	testID := fmt.Sprintf("test-%d", time.Now().UnixNano())
	
	// Get REAL metrics from the system (minimal implementation)
	metrics := SystemMetrics{
		Timestamp: time.Now(),
		CPUUsage:  getRealCPUUsage(),
		MemoryUsageMB: getRealMemoryUsageMB(),
		Status:    determineStatus(getRealCPUUsage(), getRealMemoryUsageMB()),
	}
	
	response := WebServerHardeningResponse{
		TestID:  testID,
		Status:  "success",
		Metrics: metrics,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getRealCPUUsage gets actual CPU usage percentage
func getRealCPUUsage() float64 {
	// Get CPU usage over 100ms (minimal sampling time)
	percent, _ := cpu.Percent(100*time.Millisecond, false)
	if len(percent) > 0 {
		return percent[0]
	}
	return 0.0
}

// getRealMemoryUsageMB gets actual memory usage in MB
func getRealMemoryUsageMB() int {
	// Get memory information
	v, _ := mem.VirtualMemory()
	
	// Convert bytes to MB
	return int(v.Used / 1024 / 1024)
}

// determineStatus provides a status label based on metrics
func determineStatus(cpu float64, memory int) string {
	if cpu > 85.0 || memory > 1800 {
		return "critical"
	}
	if cpu > 70.0 || memory > 1500 {
		return "warning"
	}
	return "normal"
}