package fritz

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bpicode/fritzctl/internal/stringutils"
)

// BoxData contains runtime information of the FRITZ!Box.
// codebeat:disable[TOO_MANY_IVARS]
type BoxData struct {
	Model           Model
	FirmwareVersion FirmwareVersion
	Runtime         Runtime
	Hash            string
	Status          string
}

// Model contains information about the type of the box.
type Model struct {
	Name     string
	Annex    string
	Branding string
}

// FirmwareVersion represent the FRITZ!OS version.
type FirmwareVersion struct {
	Image             string
	OsVersionMajor    string
	OsVersionMinor    string
	OsVersionRevision string
}

// Runtime contains data on how long the FRITZ!Box has been running.
type Runtime struct {
	Hours   uint64
	Days    uint64
	Months  uint64
	Years   uint64
	Reboots uint64
}

// codebeat:enable[TOO_MANY_IVARS]

// String returns a textual representation of Model.
func (m Model) String() string {
	return fmt.Sprintf("%s, ADSL standard '%s', branded as '%s'", m.Name, m.Annex,
		stringutils.DefaultIfEmpty(m.Branding, "unknown"))
}

// String returns a textual representation of FirmwareVersion.
func (v FirmwareVersion) String() string {
	return fmt.Sprintf("FRITZ!OS %s.%s (%s.%s.%s revision %s)", v.OsVersionMajor, v.OsVersionMinor,
		v.Image, v.OsVersionMajor, v.OsVersionMinor, stringutils.DefaultIfEmpty(v.OsVersionRevision, "unknown"))
}

// String returns a textual representation of Runtime.
func (r Runtime) String() string {
	return fmt.Sprintf("%d years, %d months, %d days, %d hours, %d reboots", r.Years, r.Months, r.Days, r.Hours, r.Reboots)
}

type boxDataParser struct {
}

func (p boxDataParser) parse(raw string) *BoxData {
	trimmed := strings.TrimSpace(raw)
	tokens := strings.Split(trimmed, "-")
	var tenTokens [10]string
	copy(tenTokens[:], tokens)
	if len(tokens) > 9 {
		tenTokens[9] = strings.Join(tokens[9:], "-")
	}
	return p.parseTokens(tenTokens)
}

func (p boxDataParser) parseTokens(tokens [10]string) *BoxData {
	var data BoxData
	data.Model = p.model(tokens[0], tokens[1], tokens[9])
	data.Runtime = p.runningFor(tokens[2], tokens[3])
	data.Hash = fmt.Sprintf("%s-%s", tokens[4], tokens[5])
	data.Status = tokens[6]
	data.FirmwareVersion = p.firmwareVersion(tokens[7], tokens[8])
	return &data
}

func (p boxDataParser) runningFor(firstToken, secondToken string) Runtime {
	var r Runtime
	r.Hours, _ = strconv.ParseUint(firstToken[:2], 10, 64)
	r.Days, _ = strconv.ParseUint(firstToken[2:4], 10, 64)
	r.Months, _ = strconv.ParseUint(firstToken[4:], 10, 64)
	r.Years, _ = strconv.ParseUint(secondToken[:2], 10, 64)
	r.Reboots, _ = strconv.ParseUint(secondToken[2:], 10, 64)
	return r

}

func (p boxDataParser) firmwareVersion(version, revision string) FirmwareVersion {
	return FirmwareVersion{
		Image:             version[0:3],
		OsVersionMajor:    version[3:5],
		OsVersionMinor:    version[5:],
		OsVersionRevision: revision,
	}
}
func (p boxDataParser) model(name, annex, branding string) Model {
	return Model{
		Name:     name,
		Annex:    annex,
		Branding: branding,
	}
}
