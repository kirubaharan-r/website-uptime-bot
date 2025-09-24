package monitoring

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
	"uptime-monitor/pkg/models"
	"uptime-monitor/pkg/storage"
)

// StartMonitoring starts the monitoring loop.
func StartMonitoring(s storage.Store, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			fmt.Println("--- New Check Cycle ---")
			websites, err := s.GetWebsites()
			if err != nil {
				fmt.Printf("Error getting websites: %s\n", err)
				continue
			}
			for _, website := range websites {
				go CheckWebsite(s, website)
			}
		}
	}
}

// PerformCheck performs a single check of a website.
func PerformCheck(url string) models.Check {
	check := models.Check{
		Timestamp: time.Now(),
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	start := time.Now()
	resp, err := client.Get(url)
	check.ResponseTime = time.Since(start)

	if err != nil {
		check.Status = "down"
		fmt.Printf("Error checking %s: %s\n", url, err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			check.Status = "up"
		} else {
			check.Status = "down"
		}
		check.Headers = resp.Header

		if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
			cert := resp.TLS.PeerCertificates[0]
			check.SSLInfo = &models.SSLInfo{
				Subject:    cert.Subject.String(),
				Issuer:     cert.Issuer.String(),
				NotBefore:  cert.NotBefore,
				NotAfter:   cert.NotAfter,
				IsValid:    time.Now().Before(cert.NotAfter),
			}
		}

		fmt.Printf("Checked %s: %s (%d)\n", url, check.Status, resp.StatusCode)
	}
	return check
}

// CheckWebsite performs a single check of a website and stores the result.
func CheckWebsite(s storage.Store, website models.Website) {
	check := PerformCheck(website.URL)
	check.WebsiteID = website.ID
	if err := s.AddCheck(check); err != nil {
		fmt.Printf("Error adding check for %s: %s\n", website.URL, err)
	}
}
