package moni

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Site is basically a website wich includes API interfaces
type Site struct {
	URL    string
	IP     string
	Health bool
	Pagemap

	// Crawl job info
	lastCrawled time.Time
	nextCrawl   time.Time

	*time.Timer
}

type Sitemap map[string]*Site

var (
	Sites Sitemap = make(Sitemap, 10)
)

func FetchSites() Sitemap {
	if Sites == nil || len(Sites) < 1 {
		st := GetStorage()
		if _, err := st.FetchObject("sites", Sites); err != nil {
			log.Errorf(" failed to read saved 'sites' %v", err)
			Sites = make(Sitemap)
		} else {
			log.Infof("retrived Sites from the filesystem %+v", Sites)
		}
	}
	return Sites
}

// StoreSites will attempt to store our memory version of
// the Sitemap to a file.  Hope it all works out, we will get
// a log message if there is a problem
func StoreSites() {
	if st := GetStorage(); st != nil {
		if obj, err := st.StoreObject("sites", Sites); err != nil {
			log.Errorf("failed StoreObject Sites %v", err)
		} else {
			log.Infof("Sites stored object: %+v\n", obj)
		}
	} else {
		log.Error("sites is nil")
	}
}

// AddNewSite will create a New Site from the url, including
// verify and sanitize the url and so on.
func AddNewSite(url string) *Site {

	sched.URLQ <- url

	s := &Site{
		URL:     url,
		Pagemap: make(Pagemap),
	}
	Sites[url] = s

	// Schedule a new crawl
	// Store the site
	log.Infof("Added new site %s ~ calling StoreSites()", url)

	// This should not cause any problems, que no?
	go StoreSites()
	return s
}

func (s *Site) ScheduleCrawl() {
	timer := time.AfterFunc(time.Minute*5, func() {

	})
	defer timer.Stop()
}

// RemoveSite represented by the URL from the list of sites to manage
func RemoveSite(url string) {

	// Unschedule the site from the crawler
	log.Infoln("Deleting URL ", url)
	Sites.Delete(url)
}

func (s Sitemap) Find(url string) (site *Site, ex bool) {
	site, ex = s[url]
	return site, ex
}

func (s Sitemap) Get(url string) (site *Site) {
	if site, ex := s.Find(url); ex {
		return site
	}
	return nil
}

func (s Sitemap) Exists(url string) bool {
	_, ex := s.Find(url)
	return ex
}

func (s Sitemap) Delete(url string) {
	if _, ex := s[url]; ex {
		delete(s, url)
	}
}

func (s Sitemap) Store() {
	st := GetStorage()
	if _, err := st.StoreObject("sites", s); err != nil {
		log.Errorf("failed saving Sites %v", err)
	}
}

func (s *Sitemap) Fetch() {
	st := GetStorage()
	if _, err := st.FetchObject("sites", s); err != nil {
		log.Errorf("failed to fetch Sites %v", err)
	}
}
