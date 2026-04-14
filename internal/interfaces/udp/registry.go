package udp

import (
	"log"
	"net"
	"sync"
	"time"
)

type Peer struct {
	Addr     *net.UDPAddr
	MangaID  int
	LastSeen time.Time
}

type Registry struct {
	// peers stores string(addr) -> *Peer
	peers      sync.Map
	ttl        time.Duration
	gcInterval time.Duration
}

func NewRegistry(ttl time.Duration, gcInterval time.Duration) *Registry {
	r := &Registry{
		ttl:        ttl,
		gcInterval: gcInterval,
	}
	go r.startGC()
	return r
}

func (r *Registry) Register(mangaID int, addr *net.UDPAddr) {
	r.peers.Store(addr.String(), &Peer{
		Addr:     addr,
		MangaID:  mangaID,
		LastSeen: time.Now(),
	})
	log.Printf("📡 UDP Registry: Registered peer %s for Manga %d", addr.String(), mangaID)
}

func (r *Registry) KeepAlive(addr *net.UDPAddr) {
	if val, ok := r.peers.Load(addr.String()); ok {
		peer := val.(*Peer)
		peer.LastSeen = time.Now()
		r.peers.Store(addr.String(), peer)
	}
}

func (r *Registry) startGC() {
	ticker := time.NewTicker(r.gcInterval)
	for range ticker.C {
		now := time.Now()
		r.peers.Range(func(key, value interface{}) bool {
			peer := value.(*Peer)
			if now.Sub(peer.LastSeen) > r.ttl {
				r.peers.Delete(key)
				log.Printf("🧹 UDP Registry: Pruned inactive peer %s (TTL Exceeded)", key)
			}
			return true
		})
	}
}

func (r *Registry) GetPeers(mangaID int) []*net.UDPAddr {
	var addrs []*net.UDPAddr
	r.peers.Range(func(key, value interface{}) bool {
		peer := value.(*Peer)
		// 0 can be used for global notifications if needed, 
		// but per Mẹ Architect, we filter by manga_id.
		if peer.MangaID == mangaID || mangaID == 0 {
			addrs = append(addrs, peer.Addr)
		}
		return true
	})
	return addrs
}
