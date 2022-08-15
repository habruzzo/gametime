package db

import (
	"context"
	"encoding/json"
	"fmt"
	"gametime"
	"github.com/dgraph-io/dgo/v200/protos/api"
)

func (d *Dgraph) GetAuthorByName(name string) gametime.Author {
	d.log.Info("get author by name", "(app)")
	da := d.getAuthorByName(name)
	return d.dAuthorToAuthor(da)
}

func (d *Dgraph) getAuthorByName(name string) Dauthor {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("get author by name")
	queryFmt := `{authorByName(func: eq(authorName,"%s")) {
		uid
		authorName
		dgraph.type
	}}`
	if name == "" {
		return d.authorToDAuthor(gametime.Author{Name: name})
	}
	query := fmt.Sprintf(queryFmt, name)
	// fmt.Println(query)
	resp, err := d.Dgraph.NewReadOnlyTxn().Query(context.Background(), query)
	if err != nil {
		d.log.Error("error getting author", "error", err)
		return d.authorToDAuthor(gametime.Author{Name: name})
	}
	type authors struct {
		Dauthor []Dauthor `json:"authorByName"`
	}
	var dAuthor authors
	err = json.Unmarshal(resp.Json, &dAuthor)
	if err != nil {
		d.log.Error("error unmarshalling author", "error", err)
		return d.authorToDAuthor(gametime.Author{Name: name})
	}
	if len(dAuthor.Dauthor) < 1 {
		d.log.Error("error no author results in db", "name", name)
		return d.authorToDAuthor(gametime.Author{Name: name})
	}
	return dAuthor.Dauthor[0]
}

func (d *Dgraph) insertAuthor(dAuthor Dauthor) error {
	d.newClient(d.cfg.Dgraph.Url)
	d.log.Info("insert author")
	b, err := json.Marshal(dAuthor)
	mu := api.Mutation{
		SetJson:   b,
		CommitNow: true,
	}
	//d.log.Info(string(mu.SetJson))
	resp, err := d.Dgraph.NewTxn().Mutate(context.Background(), &mu)
	if err != nil {
		d.log.Error("error inserting author", "author", dAuthor.AuthorName)
		return err
	}
	d.log.Info("", "uids", resp.Uids)
	return nil
}

func (d *Dgraph) dAuthorToAuthor(da Dauthor) gametime.Author {
	return gametime.Author{
		Id:   da.Uid,
		Name: da.AuthorName,
	}
}

func (d *Dgraph) authorToDAuthor(a gametime.Author) Dauthor {
	return Dauthor{
		Uid:        a.Id,
		AuthorName: a.Name,
		DgraphType: authorType,
	}
}
