package rp

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/missdeer/hannah/config"
	"github.com/missdeer/hannah/util"
)

const (
	chinaIPListURL = `https://cdn.jsdelivr.net/gh/17mon/china_ip_list@master/china_ip_list.txt`
)

type ipField struct {
	mask int
	next ipFieldList
}

type ipFieldList map[byte]*ipField // value - ip field

var (
	ipInChina      ipFieldList
	mutexIPInChina sync.RWMutex
)

func IPv4InChina(ipv4 net.IP) bool {
	mutexIPInChina.RLock()
	defer mutexIPInChina.RUnlock()
	current := &ipInChina
	finalHit := false
	var j byte
	for i := 0; i < 4; i++ {
		mask := 32
		b := ipv4[i]
		hit := false
		for j = 0; j < 7 && hit == false; j++ {
			if field, ok := (*current)[b]; !ok {
				if b == 0 || b == 128 {
					return false
				}
				b &= 0xFF ^ (1<<(j+1) - 1)
			} else {
				mask = field.mask
				current = &(field.next)
				hit = true
			}
		}

		if mask <= (i+1)*8 {
			finalHit = hit
			break
		}
	}

	return finalHit
}

// InChina returns true if the given IP is in China
func InChina(ip string) bool {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil || len(ipv4) < 4 {
		return false
	}
	return IPv4InChina(ipv4)
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// LoadChinaIPList loads china IP list from file
func LoadChinaIPList() error {
	ipListFile := filepath.Join(UserHomeDir(), ".china_ip_list.txt")

	stat, err := os.Stat(ipListFile)
	if os.IsNotExist(err) || stat.ModTime().Add(30*24*time.Hour).Before(time.Now()) {
		// download
		req, err := http.NewRequest("GET", chinaIPListURL, nil)
		if err != nil {
			return err
		}

		req.Header.Set("User-Agent", config.UserAgent)
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Language", "zh-CN,zh-HK;q=0.8,zh-TW;q=0.6,en-US;q=0.4,en;q=0.2")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")

		client := util.GetHttpClient()
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("status code: %d", resp.StatusCode)
		}

		content, err := util.ReadHttpResponseBody(resp)
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(ipListFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer outFile.Close()
		_, err = outFile.Write(content)
		if err != nil {
			return err
		}
	}

	inFile, err := os.OpenFile(ipListFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	mutexIPInChina.Lock()
	defer mutexIPInChina.Unlock()
	ipInChina = make(ipFieldList)
	for scanner.Scan() {
		cidr := scanner.Text()
		_, cidrIPNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}

		mask, _ := cidrIPNet.Mask.Size()
		ipv4 := cidrIPNet.IP
		next := &ipInChina
		for i := 0; i < 4; i++ {
			node, ok := (*next)[ipv4[i]]
			if !ok {
				node = &ipField{mask: mask}
				if i < 3 {
					node.next = make(ipFieldList)
				}
				(*next)[ipv4[i]] = node
			}
			next = &(node.next)
		}
	}
	return nil
}
