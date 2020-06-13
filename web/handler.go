package web

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"

	"github.com/pabloarizaluna/serverchecker"

	"github.com/buaazp/fasthttprouter"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

const (
	iconPat  = `[^>]*rel="(shortcut )?icon"[^<]*`
	tiPat    = `<title>(.*)</title>`
	cntryPat = `[Cc]ountry:\s+([^\n]+)`
	orgPat   = `OrgName:\s+([^\n]+)`
	urlPat   = `href="([\S]+)"`
	tempPat  = `(OrgName:\s+([^\n]+)|role:\s+([^\n]+))`
)

type Handler struct {
	*fasthttprouter.Router

	store serverchecker.Store
}

func NewHandler(store serverchecker.Store) *Handler {
	h := &Handler{
		Router: fasthttprouter.New(),
		store:  store,
	}
	h.GET("/check/:domain", h.checkDomain)
	h.GET("/domain", h.domains)

	return h
}

func (h *Handler) checkDomain(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("domain").(string)
	var dat map[string]interface{}

	// Getting ssl result
	request := fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s&publish=on", name)
	ready := false
	for !ready {
		resp, err := makeGetRequest(request)
		if err != nil {
			fmt.Fprint(ctx, err.Error())
			return
		}

		if err := json.Unmarshal(resp, &dat); err != nil {
			fmt.Fprint(ctx, err.Error())
			return
		}

		status := dat["status"].(string)
		if status == "READY" {
			ready = true
		}
	}

	// Checking if there any result saved
	d, err := h.store.Domain(name)
	if err != nil && err != sql.ErrNoRows {
		fmt.Fprint(ctx, err.Error())
		return
	}

	if err == sql.ErrNoRows {
		d.ID = uuid.New()
		d.Host = name
	}

	if err := updateHTMLInfo(&d); err != nil {
		fmt.Fprint(ctx, err.Error())
		return
	}

	if err := getServers(&d, dat); err != nil {
		fmt.Fprint(ctx, err.Error())
		return
	}

	grade := lowerGrade(d.Servers)
	if d.SslGrade != grade {
		d.ServersChanged = true
		d.SslGrade = grade
	} else {
		d.ServersChanged = false
	}

	h.store.CreateDomain(&d)

	result, err := json.Marshal(d)
	if err != nil {
		fmt.Fprint(ctx, err.Error())
		return
	}

	fmt.Fprint(ctx, string(result))
}

func makeGetRequest(request string) ([]byte, error) {
	resp, err := http.Get(request)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func updateHTMLInfo(d *serverchecker.Domain) error {
	request := fmt.Sprintf("https://%s", d.Host)
	resp, err := makeGetRequest(request)
	if err != nil {
		return fmt.Errorf("Error getting the HTML: %w", err)
	}

	html := string(resp)
	d.Logo = getLogoURL(html)
	d.Title = getTitle(html)

	return nil
}

func getLogoURL(html string) string {
	iconBlockRgx, _ := regexp.Compile(iconPat)
	iconBlock := iconBlockRgx.FindString(html)

	urlRgx, _ := regexp.Compile(urlPat)
	url := urlRgx.FindStringSubmatch(iconBlock)

	if len(url) < 2 {
		return ""
	}

	return url[1]
}

func getTitle(html string) string {
	tiBlockRgx, _ := regexp.Compile(tiPat)
	tiBlock := tiBlockRgx.FindStringSubmatch(html)

	if len(tiBlock) < 2 {
		return ""
	}

	return tiBlock[1]
}

func getServers(d *serverchecker.Domain, dat map[string]interface{}) error {
	endpoints := dat["endpoints"].([]interface{})
	whoisExec, _ := exec.LookPath("whois")
	var out bytes.Buffer
	cntryRgx, _ := regexp.Compile(cntryPat)
	orgRgx, _ := regexp.Compile(orgPat)
	var ss []serverchecker.Server

	for _, epIntf := range endpoints {
		ep := epIntf.(map[string]interface{})
		addr := ep["ipAddress"].(string)
		sslG := ep["grade"].(string)

		cmdWhoisExec := &exec.Cmd{
			Path:   whoisExec,
			Args:   []string{whoisExec, addr},
			Stdout: &out,
		}
		err := cmdWhoisExec.Run()
		if err != nil {
			fmt.Println("Error with Whois")
			return err
		}
		cmdRes := out.String()

		var s serverchecker.Server

		s.Address = addr
		s.SslGrade = sslG
		s.Country = cntryRgx.FindStringSubmatch(cmdRes)[1]
		s.Owner = orgRgx.FindStringSubmatch(cmdRes)[1]

		ss = append(ss, s)
	}
	d.Servers = ss

	return nil
}

func (h *Handler) domains(ctx *fasthttp.RequestCtx) {
	var history []string

	dd, err := h.store.Domains()
	if err != nil {
		fmt.Fprint(ctx, err.Error())
		return
	}

	for _, d := range dd {
		history = append(history, d.Host)
	}

	enc := json.NewEncoder(ctx)
	resp := map[string][]string{"items": history}
	enc.Encode(resp)
}

func lowerGrade(ss []serverchecker.Server) string {
	lwrgd := float32(1)
	for _, s := range ss {
		if lwrgd > gradeToNum(s.SslGrade) {
			lwrgd = gradeToNum(s.SslGrade)
		}
	}
	return numToGrade(lwrgd)
}

func gradeToNum(grade string) float32 {
	switch grade {
	case "A+":
		return 1
	case "A-":
		return 0.875
	case "A":
		return 0.75
	case "B":
		return 0.625
	case "C":
		return 0.5
	case "D":
		return 0.375
	case "E":
		return 0.25
	case "F":
		return 0.125
	case "T":
		return 0
	default:
		return -1
	}
}

func numToGrade(grade float32) string {
	switch grade {
	case 1:
		return "A+"
	case 0.875:
		return "A-"
	case 0.75:
		return "A"
	case 0.625:
		return "B"
	case 0.5:
		return "C"
	case 0.375:
		return "D"
	case 0.25:
		return "E"
	case 0.125:
		return "F"
	case 0:
		return "T"
	default:
		return "M"
	}
}
