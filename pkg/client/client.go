package client

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

type Record struct {
	ID     int
	Name   string
	FQDN   string
	Target string
}

type Domain struct {
	ID      int
	Name    string
	Records []Record
}
type Domains struct {
	apiKey       string
	linodeClient *linodego.Client
	debug        bool
	Domains      []Domain
}

func (d *Domains) getLinodeClient() {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: d.apiKey})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	client := linodego.NewClient(oauth2Client)
	d.linodeClient = &client
	d.linodeClient.SetDebug(d.debug)
}

func (d *Domains) getDomainsList(ctx context.Context) error {
	domains, err := d.linodeClient.ListDomains(ctx, nil)
	if err != nil {
		return err
	}

	for _, domain := range domains {
		r, err := d.linodeClient.ListDomainRecords(ctx, domain.ID, nil)
		if err != nil {
			return err
		}
		records := make([]Record, len(r))
		for _, subdomain := range r {
			if subdomain.Name == "" ||
				subdomain.Name == "*" {
				continue
			}
			if subdomain.Type == linodego.RecordTypeA {
				fqdn := fmt.Sprintf("%v.%v", subdomain.Name, domain.Domain)
				records = append(records, Record{
					ID:     subdomain.ID,
					Name:   subdomain.Name,
					FQDN:   fqdn,
					Target: subdomain.Target,
				})
			}
		}

		d.Domains = append(d.Domains, Domain{
			ID:      domain.ID,
			Name:    domain.Domain,
			Records: records,
		})
	}

	return nil
}

func (d *Domains) getDomainIDbyRecordID(recordID int) (int, error) {
	for _, domain := range d.Domains {
		for _, record := range domain.Records {
			if record.ID == recordID {
				return domain.ID, nil
			}
		}
	}

	return 0, fmt.Errorf("record not found: %v", recordID)
}

func (d *Domains) getFQDNbyRecordID(recordID int) string {
	for _, domain := range d.Domains {
		for _, record := range domain.Records {
			if record.ID == recordID {
				return record.FQDN
			}
		}
	}

	return ""
}

func (d *Domains) updateDomain(ctx context.Context, domainID, recordID int, newIP string) error {
	updateOpts := linodego.DomainRecordUpdateOptions{
		Target: newIP,
	}

	record, err := d.linodeClient.UpdateDomainRecord(ctx, domainID, recordID, updateOpts)
	if err != nil {
		return err
	}

	fmt.Printf("Domain %v update to %v\n", d.getFQDNbyRecordID(recordID), record.Target)

	return nil
}

func (d *Domains) printDomains() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Domain ID", "Main Domain", "Record ID", "Sub Domain", "Current IP"})

	var rows []table.Row
	for _, domain := range d.Domains {
		for _, record := range domain.Records {
			rows = append(rows, table.Row{
				domain.ID,
				domain.Name,
				record.ID,
				record.Name,
				record.Target,
			})
		}
	}
	t.AppendRows(rows)

	t.Render()
}

// Client provides the method to get the domains and update a record with a new IP
func Client(ctx context.Context, apiKey string, debug bool, recordID int, newIP string) error {
	domains := Domains{
		apiKey: apiKey,
		debug:  debug,
	}

	// Get Client for the session
	domains.getLinodeClient()

	// Always get the list of Domains even if we don't update any domain
	err := domains.getDomainsList(ctx)
	if err != nil {
		return err
	}

	if recordID == 0 && newIP == "" {
		domains.printDomains()
		return nil
	}

	// The user will pass the recordID we also need the ID
	domainID, err := domains.getDomainIDbyRecordID(recordID)
	if err != nil {
		return err
	}

	// Update the domain
	err = domains.updateDomain(ctx, domainID, recordID, newIP)

	return err
}
