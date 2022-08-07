package db

import (
	"context"
	"encoding/json"
	"fmt"
	"gametime"
	"gametime/config"
	"os"
	"testing"

	"github.com/dgraph-io/dgo/v200/protos/api"

	"github.com/google/uuid"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T) {
	initSchema(t)
	t.Run("test insert post success", testInsertPostSuccess)
	t.Run("test get author success", testGetAuthorSuccess)
	t.Run("test get game success", testGetGameSuccess)
	t.Run("test get review success", testGetReviewSuccess)

	//t.Run("test get posts success", testGetPostsSuccess)
}

func setupTest(t *testing.T) (*Dgraph, gametime.Review) {
	fn := "../config/rubric_-_80_days.json"
	f, err := os.Open(fn)
	assert.NoError(t, err)
	var buf [20000]byte
	n, err := f.Read(buf[:])
	f.Close()
	assert.NoError(t, err)
	var p gametime.Review
	p.Text = string(buf[:n])
	assert.NoError(t, json.Unmarshal(buf[:n], &p))
	log := logrus.New()
	d := NewDgraph(log, config.Config{
		Dgraph: config.Db{
			Url: "localhost:9080",
		},
	})
	return d, p
}

func testInsertPostSuccess(t *testing.T) {
	d, p := setupTest(t)
	fmt.Println(p.Slug, p.Author.Name, p.Game.Title)
	assert.NoError(t, d.InsertPost(p))
}

func testGetPostSuccess(t *testing.T) {
	testInsertPostSuccess(t)
	d, p := setupTest(t)
	rev, err := d.GetReviewBySlug("80_days")
	assert.NoError(t, err)

	assert.Equal(t, p, rev)

	post, err := d.GetPostByReview(rev)
	assert.NoError(t, err)
	assert.Equal(t, post.ReviewId, rev.Id)
	assert.NotEqual(t, post.Id, "")
}

func testGetAuthorSuccess(t *testing.T) {
	d, _ := setupTest(t)
	name := fmt.Sprintf("test_author_%s", uuid.New().String())
	da := d.getAuthorByName(name)
	assert.Equal(t, name, da.AuthorName)
	assert.Equal(t, "", da.Uid)
	assert.NoError(t, d.insertAuthor(da))
	da = d.getAuthorByName(name)
	assert.Equal(t, da.AuthorName, name)
	assert.NotEqual(t, "", da.Uid)
}

func testGetGameSuccess(t *testing.T) {
	d, _ := setupTest(t)
	title := fmt.Sprintf("test_game_%s", uuid.New().String())
	da := d.getGameByTitle(title)
	assert.Equal(t, title, da.GameTitle)
	assert.Equal(t, "", da.Uid)
	assert.NoError(t, d.insertGame(da))
	da = d.getGameByTitle(title)
	assert.Equal(t, da.GameTitle, title)
	assert.NotEqual(t, "", da.Uid)
}

func testGetReviewSuccess(t *testing.T) {
	d, _ := setupTest(t)
	slug := fmt.Sprintf("test_review_%s", uuid.New().String())
	da, err := d.getReviewBySlug(slug)
	assert.NoError(t, err)
	assert.Equal(t, slug, da.Slug)
	assert.Equal(t, "", da.Uid)
	assert.NoError(t, d.insertReview(da))
	da, err = d.getReviewBySlug(slug)
	assert.NoError(t, err)
	assert.Equal(t, da.Slug, slug)
	assert.NotEqual(t, "", da.Uid)
}

func testGetPostRecentSuccess(t *testing.T) {
	testInsertPostSuccess(t)
	d, _ := setupTest(t)
	d, _ = setupTest(t)
	dp, err := d.GetPostsByMostRecent()
	assert.NoError(t, err)

	da, err := d.getReviewBySlug("80_days")
	assert.NoError(t, err)

	assert.Equal(t, da.Slug, "80_days")
	assert.NotEqual(t, dp[0].ReviewId, da.Uid)
}

func initSchema(t *testing.T) {
	d, _ := setupTest(t)
	d.newClient(d.cfg.Dgraph.Url)
	fn := "../conf/schema.dql"
	f, err := os.Open(fn)
	assert.NoError(t, err)
	var buf [20000]byte
	n, err := f.Read(buf[:])
	f.Close()
	assert.NoError(t, err)
	schema := string(buf[:n])

	op := api.Operation{Schema: schema}
	err = d.Alter(context.TODO(), &op)
	assert.NoError(t, err)
}
