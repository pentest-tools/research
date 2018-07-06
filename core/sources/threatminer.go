package sources

import core "github.com/subfinder/research/core"
import "net/http"
import "net"
import "time"
import "bufio"
import "regexp"
import "strings"

type Threatminer struct{}

func (source *Threatminer) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		httpClient := &http.Client{
			Timeout: time.Second * 60,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		}

		domainExtractor, err := regexp.Compile("q=([a-zA-Z0-9\\*_.-]+\\." + domain + ")")
		if err != nil {
			results <- &core.Result{Type: "threatminer", Failure: err}
			return
		}

		resp, err := httpClient.Get("https://www.threatminer.org/getData.php?e=subdomains_container&q=" + domain + "&t=0&rt=10&p=1")
		if err != nil {
			results <- &core.Result{Type: "threatminer", Failure: err}
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				strParts := strings.Split(str, "q=")
				if len(strParts) >= 1 {
					str = strings.Join(strParts[:len(strParts)], "")
					results <- &core.Result{Type: "threatminer", Success: str}
				}
			}
		}
	}(domain, results)
	return results
}
