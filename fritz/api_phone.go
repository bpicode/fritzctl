package fritz

import (
	"github.com/bpicode/fritzctl/httpread"
	"github.com/bpicode/fritzctl/internal/errors"
)

// NewPhone creates a client for interaction with the FB phone sector.
func NewPhone(client *Client) Phone {
	return &phone{client: client}
}

// Phone describes the supported operations.
type Phone interface {
	Calls() ([]Call, error)
}

// Call contains the data for one phone call record.
// codebeat:disable[TOO_MANY_IVARS]
type Call struct {
	Type           string // "1"=incoming, "2"=missed call, "3"=unknown, "4"=outgoing.
	Date           string
	Caller         string
	PhoneNumber    string
	Extension      string
	OwnPhoneNumber string
	Duration       string
}

// codebeat:enable[TOO_MANY_IVARS]

type phone struct {
	client *Client
}

// Calls reads the phone call record list from FB.
func (p *phone) Calls() ([]Call, error) {
	url := p.client.query().path(phoneListURI).query("csv", "").build()
	records, err := httpread.Csv(p.client.getf(url), ';')
	if err != nil {
		return nil, errors.Wrapf(err, "unable read data for phone calls")
	}
	var calls []Call
	for _, r := range records[2:] { // Skip first two lines.
		calls = append(calls, convertRecordToCall(r))
	}
	return calls, err
}

func convertRecordToCall(r []string) Call {
	return Call{
		Type:           r[0],
		Date:           r[1],
		Caller:         r[2],
		PhoneNumber:    r[3],
		Extension:      r[4],
		OwnPhoneNumber: r[5],
		Duration:       r[6],
	}
}
