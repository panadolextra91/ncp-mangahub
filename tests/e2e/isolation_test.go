package e2e

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsolation_SlowTCPConsumer(t *testing.T) {
	s, err := StartServer("../../mangahub_server")
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer s.Stop()

	// 1. Register Admin and Login via HTTP to get token
	token := getAuthToken(t, s.Port)

	// 2. Connect a "Good" TCP Client
	goodConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", s.TCPPort))
	assert.NoError(t, err)
	defer goodConn.Close()
	fmt.Fprintf(goodConn, "AUTH %s\n", token)
	goodReader := bufio.NewReader(goodConn)
	goodReader.ReadString('\n') // Consume 200 OK

	// 3. Connect a "Bad" TCP Client (Slow Consumer)
	badConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", s.TCPPort))
	assert.NoError(t, err)
	defer badConn.Close()
	fmt.Fprintf(badConn, "AUTH %s\n", token)
	// We do NOT read from badConn, letting its TCP buffer fill up.

	// 4. Trigger an event via HTTP (Create Manga)
	// We create multiple mangas to ensure the bus is active
	for i := 1; i <= 5; i++ {
		createManga(t, s.Port, token, fmt.Sprintf("Speed-Manga-%d", i))
	}

	// 5. Verify the Good Client receives all events promptly
	// If the bad client blocked the bus, the good client would be stalled.
	for i := 1; i <= 5; i++ {
		goodConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		line, err := goodReader.ReadString('\n')
		assert.NoError(t, err, "Good client should receive event %d", i)
		assert.Contains(t, line, "Speed-Manga")
	}

	// 6. Conclusion: The bad client's overflow didn't stall the good client.
}

func getAuthToken(t *testing.T, port int) string {
	url := fmt.Sprintf("http://localhost:%d/api/auth/register", port)
	body := `{"username":"admin","password":"password","role":"admin"}`
	http.Post(url, "application/json", bytes.NewBufferString(body))

	url = fmt.Sprintf("http://localhost:%d/api/auth/login", port)
	body = `{"username":"admin","password":"password"}`
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("Login POST failed: %v", err)
	}
	defer resp.Body.Close()

	var res struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil || res.Token == "" {
		t.Fatalf("Failed to decode token at port %d: %v", port, err)
	}
	return res.Token
}

func createManga(t *testing.T, port int, token string, title string) {
	url := fmt.Sprintf("http://localhost:%d/api/manga", port)
	body := fmt.Sprintf(`{"title":"%s","author":"E2E","description":"Test"}`, title)
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()
}
