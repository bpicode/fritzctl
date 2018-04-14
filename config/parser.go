package config

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	errors2 "github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Parser defines one method, Parse, which reads from 3rd party source(s).
type Parser interface {
	Parse() (*Config, error)
}

type parser struct {
	sources []source
}

type source struct {
	openCloser
	Decode
}

type openCloser interface {
	Open() (io.Reader, error)
	io.Closer
}

// Decode uses an io.Reader, whose read bytes shall be marshaled into the passed interface. An error is returned if the
// operation did not work.
type Decode func(io.Reader, interface{}) error

// Place is a configuration option for a parser.
type Place func(*parser)

// NewParser constructs a Parser that looks for configuration in the passed Places. If multiple Places would spawn
// correct configurations, the first-without-error is used.
func NewParser(places ...Place) Parser {
	s := &parser{}
	for _, p := range places {
		p(s)
	}
	return s
}

type homeDirOpenCloser struct {
	path string
	io.Closer
	user func() (*user.User, error)
}

func (h *homeDirOpenCloser) Open() (io.Reader, error) {
	hd := homeDirOf(h.user)
	path, err := hd(h.path)
	if err != nil {
		return nil, errors2.Wrapf(err, "cannot open '%s' in current user's home directory", h.path)
	}
	f, err := os.Open(path)
	h.Closer = f
	return f, err
}

// InHomeDir looks for a file inside a user's home directory.
func InHomeDir(u func() (*user.User, error), path string, decoder Decode) Place {
	return func(p *parser) {
		s := source{}
		s.Decode = decoder
		s.openCloser = &homeDirOpenCloser{path: path, user: u}
		p.sources = append(p.sources, s)
	}
}

// InDir looks for a file inside a given directory.
func InDir(dir string, file string, decoder Decode) Place {
	return func(p *parser) {
		s := source{}
		s.Decode = decoder
		s.openCloser = &homeDirOpenCloser{
			path: file,
			user: func() (*user.User, error) {
				return &user.User{HomeDir: dir}, nil
			},
		}
		p.sources = append(p.sources, s)
	}
}

// JSON returns a Decode function that uses JSON format.
func JSON() Decode {
	return func(r io.Reader, v interface{}) error {
		return json.NewDecoder(r).Decode(v)
	}
}

// YAML returns a Decode function that uses YML format.
func YAML() Decode {
	return func(r io.Reader, v interface{}) error {
		bytes, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(bytes, v)
	}
}

// Parse runs though all places and tries to parse the configuration. If none of the places yields a successful config,
// an error is returned.
func (p *parser) Parse() (*Config, error) {
	var errs []error
	for _, s := range p.sources {
		r, err := s.Open()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		c, err := p.decode(s, r)
		s.Close()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		return c, nil
	}
	err := p.joinErrors(errs)
	return nil, errors2.Wrapf(err, "unable to find a usable config source")
}

func (p *parser) joinErrors(errs []error) error {
	var msgs = make([]string, 0, len(errs))
	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}
	return errors.New("no valid config found in the following locations:\n  " + strings.Join(msgs, "\n  "))
}

func (p *parser) decode(s source, r io.Reader) (*Config, error) {
	cfg := Config{}
	net := Net{}
	pki := Pki{}
	login := Login{}
	err := s.Decode(r, &struct {
		*Net
		*Login
		*Pki
	}{&net, &login, &pki})
	cfg.Pki = &pki
	cfg.Login = &login
	cfg.Net = &net
	return &cfg, err
}
