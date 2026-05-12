package e2e

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Server handles the lifecycle of a MangaHub server instance for E2E testing.
type Server struct {
	BinaryPath string
	DBPath     string
	Port       int
	TCPPort    int
	UDPPort    int
	GRPCPort   int
	Cmd        *exec.Cmd
}

func StartServer(binaryPath string) (*Server, error) {
	// 1. Get random ports
	p1, _ := getFreePort()
	p2, _ := getFreePort()
	p3, _ := getFreePort()
	p4, _ := getFreePort()

	dbPath := fmt.Sprintf("e2e_test_%d.db", p1)
	
	s := &Server{
		BinaryPath: binaryPath,
		DBPath:     dbPath,
		Port:       p1,
		TCPPort:    p2,
		UDPPort:    p3,
		GRPCPort:   p4,
	}

	// 2. Start process with env overrides
	s.Cmd = exec.Command(binaryPath)
	s.Cmd.Env = append(os.Environ(),
		"PORT="+strconv.Itoa(p1),
		"TCP_PORT="+strconv.Itoa(p2),
		"UDP_PORT="+strconv.Itoa(p3),
		"GRPC_PORT="+strconv.Itoa(p4),
		"DB_PATH="+dbPath,
		"JWT_SECRET=e2e-secret-key",
	)
	
	logFile, err := os.Create(fmt.Sprintf("server_%d.log", p1))
	if err == nil {
		s.Cmd.Stdout = logFile
		s.Cmd.Stderr = logFile
	}

	if err := s.Cmd.Start(); err != nil {
		return nil, err
	}

	// 3. Wait for healthy
	healthy := false
	for i := 0; i < 20; i++ {
		time.Sleep(100 * time.Millisecond)
		conn, err := net.DialTimeout("tcp", "localhost:"+strconv.Itoa(p1), 50*time.Millisecond)
		if err == nil {
			conn.Close()
			healthy = true
			break
		}
	}

	if !healthy {
		s.Stop()
		return nil, fmt.Errorf("server failed to start within timeout")
	}

	return s, nil
}

func (s *Server) Stop() {
	if s == nil {
		return
	}
	if s.Cmd != nil && s.Cmd.Process != nil {
		s.Cmd.Process.Signal(os.Interrupt)
		s.Cmd.Wait()
	}
	if s.DBPath != "" {
		os.Remove(s.DBPath)
	}
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
