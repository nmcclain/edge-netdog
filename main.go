package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var version = "1.0.0"

var usage = `edge-netdog
Usage:
  edge-netdog [options] <config>
  edge-netdog --version
  edge-netdog --help
`

func main() {
	cfg, err := getcfg()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	if err := netMonitor(cfg); err != nil {
		log.Fatalf("Monitor error: %v", err)
	}
	log.Printf("Exiting")
}

func netMonitor(cfg Config) error {
	failedAttempts := 0
	remediationActionIdx := 0
	for {
		if err := netCheck(cfg.Global.TargetUrl, cfg.Global.TargetMatch); err != nil {
			failedAttempts++
			if failedAttempts <= cfg.Global.TargetAttempts {
				log.Printf("Check failed, trying again: %s", err)
			} else {
				var action string
				if remediationActionIdx >= len(cfg.Actions) {
					action = cfg.Actions[len(cfg.Actions)-1]
					log.Errorf("Check failed, re-trying final remediation action: %s [error: %s]", cfg.Actions[len(cfg.Actions)-1], err)
				} else {
					action = cfg.Actions[remediationActionIdx]
					log.Errorf("Check failed, trying remediation action %d: %s [error: %s]", remediationActionIdx+1, cfg.Actions[remediationActionIdx], err)
				}
				if err := remediate(action); err != nil {
					log.Errorf("Remediation error [action %s]: %s", err)
				}
				remediationActionIdx++
				log.Debugf("Waiting %s for remediation before trying again.", cfg.Global.ActionDelayRaw)
				time.Sleep(cfg.Global.ActionDelay)
				continue
			}
		} else {
			log.Debugf("Check OK")
			failedAttempts = 0
			remediationActionIdx = 0
		}
		time.Sleep(cfg.Global.MonitorInterval)
	}
	return nil
}

func netCheck(target, targetMatch string) error {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Invalid response code: %v", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if !strings.Contains(string(body), targetMatch) {
		return fmt.Errorf("target_match not found: %s", targetMatch)
	}
	return nil
}

func remediate(action string) error {
	out, err := exec.Command("bash", "-c", action).Output()
	if err != nil {
		return fmt.Errorf("Error running command [%s]: %v", action, err)
	}
	log.Debugf("Command output: %s", out)
	return nil
}

var client = &http.Client{
	Transport: &http.Transport{
		ResponseHeaderTimeout: 20 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   20 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 1,
	},
}
