// Copyright (c) 2013 Erik St. Martin, Brian Ketelsen. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package msg

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type Service struct {
	UUID        string
	Name        string
	Version     string
	Environment string
	Region      string
	Host        string
	Port        uint16
	TTL         uint32 // Seconds
	Expires     time.Time
	Callback    map[string]Callback `json:"-"` // Callbacks are found by UUID
}

// Returns the amount of time remaining before expiration
func (s *Service) RemainingTTL() uint32 {
	d := s.Expires.Sub(time.Now())
	ttl := uint32(d.Seconds())

	if ttl < 1 {
		return 0
	}
	return ttl
}

// Updates TTL property to the RemainingTTL
func (s *Service) UpdateTTL() {
	s.TTL = s.RemainingTTL()
}

type Callback struct {
	UUID string

	// Name of the service
	Name        string
	Version     string
	Environment string
	Region      string
	Host        string

	Reply string
	Port  uint16
}

// Call calls the callback and performs the HTTP request.
func (c Callback) Call(s Service) {
		log.Println("Performing callback to:", c.Reply, c.Port)
		// TODO(miek): Use DELETE with some data
		resp, err := http.Get("http://" + c.Reply + ":" + strconv.Itoa(int(c.Port)) + "/skydns/callbacks/" + c.UUID)
		if err != nil {
			return
		}
		resp.Body.Close()
		return
}
