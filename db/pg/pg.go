package pg

import (
	"os"
	"os/user"
	"net/url"
	"strings"
	"crypto/tls"
)

import "github.com/go-pg/pg"

// Pg: Global instance will init when ENV:PGURL is set
var Pg *pg.DB

func init() {
	var pgURL = "postgres://localhost:5432"
	if envURL := os.Getenv("PGURL"); envURL != "" {
		pgURL = envURL
	}
	Pg = NewPg(pgURL)
}

// NewPg will create a new pg instance from pg url
// returns nil on error
func NewPg(pgURL string) *pg.DB {
	parsedUrl, err := url.Parse(pgURL)
	if err != nil {
		return nil
	}

	// scheme
	if parsedUrl.Scheme != "postgres" && parsedUrl.Scheme != "postgresql" {
		return nil
	}

	// host
	options := &pg.Options{
		Addr: parsedUrl.Host,
	}

	// port
	if !strings.Contains(options.Addr, ":") {
		options.Addr = options.Addr + ":5432"
	}

	// username and password
	if parsedUrl.User != nil {
		options.User = parsedUrl.User.Username()
		if password, ok := parsedUrl.User.Password(); ok {
			options.Password = password
		}
	}

	// use current user as default
	if options.User == "" {
		if userinfo, err := user.Current(); err != nil {
			return nil
		} else {
			options.User = userinfo.Name
		}
	}

	// database: use postgres as default
	if len(strings.Trim(parsedUrl.Path, "/")) > 0 {
		options.Database = parsedUrl.Path[1:]
	} else {
		options.Database = "postgres"
	}

	// ssl mode
	query, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return nil
	}

	if sslMode, ok := query["sslmode"]; ok && len(sslMode) > 0 {
		switch sslMode[0] {
		case "allow":
			fallthrough
		case "prefer":
			options.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		case "disable":
			options.TLSConfig = nil
		default:
			return nil
		}
	} else {
		// disable tls on default
		options.TLSConfig = nil
	}

	delete(query, "sslmode")
	if len(query) > 0 {
		return nil
	}

	return pg.Connect(options)
}
