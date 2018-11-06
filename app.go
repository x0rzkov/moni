package moni

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// All global variables
var (
	config *Configuration

	app   *App
	acl   *AccessList
	sites Sitemap
	pages Pagemap

	urlQ   *URLQ
	crawlQ *CrawlQ
	saveQ  *SaveQ

	server *http.Server
)

// After main is called and args are parsed
func (a *App) Init() {
	app = NewApp(config)
	acl = initACL()
	sites = initSites()
	pages = initPages()

	urlQ = NewURLQ()
	crawlQ = NewCrawlQ()
	saveQ = NewSaveQ()

	server = NewServer(config.Addrport)
}

func (app *App) Start() {

	go urlQ.Watch()
	go crawlQ.Watch()
	go saveQ.Watch()

	server.ListenAndServe()
}

// ====================================================================
//                           App
// ====================================================================

// App is the One TRUE App! All Hail the App!  It is the global context
// of everything.  It contains some information for the web page to be
// displayed, it also maintains configurations and managers the server
// and the scheduler.
type App struct {

	// Basic meta stuff for App web page and content
	Title string // name of the page (url title)
	Name  string // name for fun and profit
	Tmpl  string // base or frame template
	Frag  string // request.URL.Fragment

	// Tmplates to handle html and text formatting
	AppTemplates

	*log.Entry
}

// NewApp will produce a new App
func NewApp(cfg *Configuration) (app *App) {
	app = &App{
		Name:  "ClowOpsApp",
		Title: "Clowd ~ Operations",
		Tmpl:  "index.html",
	}
	app.Title = app.Name
	SetConfig(cfg)

	// Setup the logger
	app.Entry = log.WithFields(log.Fields{
		"app":  app.Title,
		"tmpl": app.Tmpl,
	})
	return app
}

// NewApp will produce a new App
func NewTestApp(config *Configuration) (app *App) {
	app = &App{
		Name:  "ClowOpsApp",
		Title: "Clowd ~ Operations",
		Tmpl:  "index.html",
	}
	app.Title = app.Name
	return app
}

// ====================================================================
//                      App Templates
// ====================================================================

// Contains various pointers to Go templates and the compiled
// version of the templates.
type AppTemplates struct {
	TmplBasedir string
	TmplName    string
	*template.Template
}

// Acculmulate the data needed for the template
type Appdata struct {
	Sites []*Site
	*Configuration
}

// Builder constructs (and sends) the response back to the
// user.  It determines with template pieces to put together,
// assembles them and off they go
func (app *App) PrepareTemplates(tmpldir string) {
	pattern := filepath.Join(tmpldir, "*.html")
	app.Template = template.Must(template.ParseGlob(pattern))
}

func (app *App) DumpTemplates() {
	fmt.Println("Templates: ", app.Template.Name())
	fmt.Println(app.DefinedTemplates())
}

// ====================================================================
//                      App Assembler
// ====================================================================

// Assemble traverses our local representation of the outgoing documents,
// occaisionally run stuff through a template, writing out successful
// stuff as required.
func (app *App) Assemble(w http.ResponseWriter, tmplname string) {
	// Here we go, create our html for our site.  Building the page happens
	// in two parts.  1. A semi-generic frame is created with designated areas
	// can be overwritten with application specific information.
	if app.Template == nil {
		app.PrepareTemplates(config.Tmpldir)
	}

	d := &Appdata{
		Configuration: config,
	}

	if err := app.ExecuteTemplate(w, "index.html", d); err != nil {
		app.Fatalln(err)
	}
}
