package db

import (
	"context"
	"encoding/json"
	"fmt"
	"gametime"
	"github.com/dgraph-io/dgo/v200/protos/api"
)

func (d *Dgraph) InsertStatus(s gametime.Status) {
	ds := d.getStatusByName(s.Name)
	if ds.Uid == "" {
		d.log.Info("status doesnt exist in db", "name", s.Name)
		d.insertStatus(ds)
	}
}

func (d *Dgraph) insertStatus(ds DgameStatus) {
	b, err := json.Marshal(ds)
	if err != nil {
		d.log.Error("error marshalling status", "error", err, "name", ds.Name)
		return
	}
	mu := api.Mutation{
		SetJson:   b,
		CommitNow: true,
	}
	//d.log.Info(string(mu.SetJson))
	resp, err := d.Dgraph.NewTxn().Mutate(context.Background(), &mu)
	if err != nil {
		d.log.Error("error inserting status", "name", ds.Name)
		return
	}
	d.log.Info("", "uids", resp.Uids)
	return
}

func (d *Dgraph) GetStatusByName(name string) gametime.Status {
	d.log.Info("get status by name", "(app)")
	da := d.getStatusByName(name)
	return d.dstatusToStatus(da)
}

func (d *Dgraph) getStatusByName(name string) DgameStatus {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("get status by name")
	queryFmt := `{statusByName(func: eq(statusName,"%s")) {
		uid
		statusName
		dgraph.type
	}}`
	if name == "" {
		return d.statusToDstatus(gametime.Status{Name: name})
	}
	query := fmt.Sprintf(queryFmt, name)
	// fmt.Println(query)
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting status", "error", err)
		return d.statusToDstatus(gametime.Status{Name: name})
	}
	type statuss struct {
		Dstatus []DgameStatus `json:"statusByName"`
	}
	var dstatus statuss
	err = json.Unmarshal(resp.Json, &dstatus)
	if err != nil {
		d.log.Error("error unmarshalling status", "error", err)
		return d.statusToDstatus(gametime.Status{Name: name})
	}
	if len(dstatus.Dstatus) < 1 {
		d.log.Error("error no status results in db", "name", name)
		return d.statusToDstatus(gametime.Status{Name: name})
	}
	return dstatus.Dstatus[0]
}

func (d *Dgraph) dstatusToStatus(ds DgameStatus) gametime.Status {
	switch ds.Name {
	case gametime.Unknown.Name:
		s := gametime.Unknown
		s.Id = ds.Uid
		return s
	case gametime.Wishlist.Name:
		s := gametime.Wishlist
		s.Id = ds.Uid
		return s
	case gametime.Installed.Name:
		s := gametime.Installed
		s.Id = ds.Uid
		return s
	case gametime.PlayedSome.Name:
		s := gametime.PlayedSome
		s.Id = ds.Uid
		return s
	case gametime.PlayedMost.Name:
		s := gametime.PlayedMost
		s.Id = ds.Uid
		return s
	case gametime.Completed.Name:
		s := gametime.Completed
		s.Id = ds.Uid
		return s
	case gametime.Reviewed.Name:
		s := gametime.Reviewed
		s.Id = ds.Uid
		return s
	case gametime.WontReview.Name:
		s := gametime.WontReview
		s.Id = ds.Uid
		return s
	default:
		s := gametime.Unknown
		s.Id = ds.Uid
		return s
	}
}

func (d *Dgraph) statusToDstatus(s gametime.Status) DgameStatus {
	return DgameStatus{
		Name:       s.Name,
		Uid:        s.Id,
		DgraphType: statusType,
	}
}
