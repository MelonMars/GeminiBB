package main

import (
	"embed"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/pitr/gig"
	db "github.com/userName/otherModule"
)

var Base db.DB
var Generator *db.IDGenerator

//go:embed tmpl/*.gmi
var tmpls embed.FS

type Template struct {
	t *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, _ gig.Context) error {
	return t.t.ExecuteTemplate(w, fmt.Sprintf("%s.gmi", name), data)
}

func main() {
	Base, Generator = db.Init()
	g := gig.Default()
	g.Handle("/section/:section/:thread/newComment", postComment)
	g.Handle("/section/:section/newPost", postThread)
	g.Handle("/section/:section", viewSection)
	g.Handle("/section/:section/:thread", viewThread)
	g.Handle("/", DisplayHome)
	g.Renderer = &Template{t: template.Must(template.ParseFS(tmpls, "tmpl/*.gmi"))}
	g.Run("host.cert", "host.key")
}

func DisplayHome(c gig.Context) error {
	return c.Render("home", struct {
		Section []db.DisplaySection
	}{
		Section: db.GetIndex(&Base),
	})
}

func postThread(c gig.Context) error {
	q := c.URL().RawQuery
	if len(q) == 0 {
		return c.NoContent(gig.StatusInput, "Enter Thread name and title separated by ```:")
	}
	post, err := url.QueryUnescape(q)
	if err != nil {
		return c.NoContent(gig.StatusInput, "Could not parse post in net/url.QueryUnescape(q)")
	}

	chunks := strings.SplitN(post, "```", 2)
	var content, title string
	switch len(chunks) {
	case 0:
		return c.NoContent(gig.StatusInput, "PUT TEXT")
	case 1:
		return c.NoContent(gig.StatusInput, "Enter Thread name and title separated by ```:")
	case 2:
		title = chunks[0]
		content = chunks[1]
	default:
		return c.NoContent(gig.StatusServerUnavailable, "strings.SplitN broke wtf")
	}

	section, err := strconv.Atoi(c.Param("section"))
	_, Base = db.NewThread(uint64(section), &Base, title, content, Generator)
	return viewSection(c)
}

func postComment(c gig.Context) error {
	q := c.URL().RawQuery
	if len(q) == 0 {
		return c.NoContent(gig.StatusInput, "Enter comment")
	}
	post, err := url.QueryUnescape(q)
	if err != nil {
		return c.NoContent(gig.StatusInput, "Could not parse post in net/url.QueryUnescape(q)")
	}

	section, err := strconv.Atoi(c.Param("section"))
	if err != nil {
		return c.NoContent(gig.StatusServerUnavailable, "Slight problem")
	}
	var s uint64 = uint64(section)
	thread, err := strconv.Atoi(c.Param("thread"))
	if err != nil {
		return c.NoContent(gig.StatusServerUnavailable, "Slight Problem")
	}
	var t uint64 = uint64(thread)
	_, Base = db.NewComment(s, t, &Base, post)

	return viewThread(c)
}

func viewSection(c gig.Context) error {
	section, err := strconv.Atoi(c.Param("section"))
	if err != nil {
		return c.NoContent(gig.StatusServerUnavailable, "Slight Problem brwy, srs bsns")
	}
	return c.Render("section", struct {
		Threads []db.DisplayThread
		Section int
		Name    string
	}{
		Threads: db.GetSection(uint64(section), &Base),
		Section: section,
		Name:    db.GetSectionName(uint64(section), &Base),
	})
}

func viewThread(c gig.Context) error {
	section, err := strconv.Atoi(c.Param("section"))
	thread, err := strconv.Atoi(c.Param("thread"))
	if err != nil {
		return c.NoContent(gig.StatusServerUnavailable, "Slight problem mb, this is some srs bsns")
	}
	threadName, threadId, sectionId, sectionName, content, comments := db.GetThread(uint64(section), uint64(thread), &Base)
	return c.Render("thread", struct {
		Thread      string
		ThreadID    int
		SectionID   int
		SectionName string
		Content     string
		Comments    []db.Comment
	}{
		Thread:      threadName,
		ThreadID:    threadId,
		SectionID:   sectionId,
		SectionName: sectionName,
		Content:     content,
		Comments:    comments,
	})
}
