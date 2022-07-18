//controllers/gorm.go
package controllers

import (
	"blog_holdongametime/app/models"
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/revel/revel"
)

const dsn string = "host=localhost user=postgres password=postgres dbname=postgres port=6001 sslmode=disable"

// type: revel controller with `*gorm.DB`
type GormController struct {
	*revel.Controller
	DB *gorm.DB
}

// it can be used for jobs
var Gdb *gorm.DB

func (c *GormController) SetDB() revel.Result {
	c.DB, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return nil
}

func (c *GormController) GetPost(slug string) models.GormPost {
	var post models.GormPost
	c.DB.Where("slug=?", slug).First(&post)

	return post
}

func (c *GormController) GetPostList() []models.GormPost {
	var p []models.GormPost
	result := c.DB.Order("created_at desc").Find(&p)
	if result.Error != nil {
		fmt.Println(result.Error)
		panic(result.Error)
	}
	return p
}

func (c *GormController) BuildReviewSkeleton(slug string) *models.ReviewSkeleton {
	var r models.Review
	c.DB.Where("id=?", models.GenerateId(slug, "review")).First(&r)
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	var rs models.ReviewSkeleton
	err := json.Unmarshal(r.JsonBytes, &rs)
	if err != nil {
		fmt.Println(string(r.JsonBytes))
		panic(err)
	}
	return &rs
}

func LoadGormSlugs() []JsonMapping {
	data, err := ioutil.ReadFile(jsonSlugPath)
	if err != nil {
		panic(err)
	}
	var j []JsonMapping
	err = json.Unmarshal(data, &j)
	if err != nil {
		panic(err)
	}
	return j
}

func LoadGormGames() []*models.GormGame {
	data, err := ioutil.ReadFile(jsonGamePath)
	if err != nil {
		panic(err)
	}
	var g []models.GormGame
	err = json.Unmarshal(data, &g)
	if err != nil {
		panic(err)
	}
	var gClean []*models.GormGame
	for _, v := range g {
		gClean = append(gClean, models.NewGormGame(v.Title, v.Slug, v.Platform, v.Publisher, v.Creator, v.ReleaseDate, v.GamesDbLink, v.Status, v.Images))
	}
	return gClean
}

func LoadGormPosts(gdb *gorm.DB) []*models.GormPost {
	data, err := ioutil.ReadFile(jsonPostPath)

	if err != nil {
		panic(err)
	}
	var p []models.PostJson
	err = json.Unmarshal(data, &p)
	if err != nil {
		panic(err)
	}
	var pClean []*models.GormPost
	for _, v := range p {
		var g models.GormGame
		result := gdb.Where("slug=?", v.Slug).First(&g)
		if result.Error != nil {
			fmt.Println("FATAL", result.Error)
			panic(result.Error)
		}
		path := fmt.Sprintf("%s.json", v.Slug)
		pClean = append(pClean, models.NewGormPost(v.Title, &g, v.Slug, path, v.PublishDate))
	}
	return pClean
}

func LoadReviews(posts []*models.GormPost) []*models.Review {
	var rClean []*models.Review
	for _, v := range posts {
		rClean = append(rClean, models.NewReview(v))
	}
	return rClean
}

func InitGormDB() {
	Gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	//Gdb.LogMode(true) // Print SQL statements
	if err != nil {
		fmt.Println("FATAL", err)
		panic(err)
	}
	Gdb.AutoMigrate(&models.GormGame{})
	Gdb.AutoMigrate(&models.GormPost{})
	Gdb.AutoMigrate(&models.Review{})

	// slugs := LoadGormSlugs()
	// for _, v := range slugs {
	// 	if err := Gdb.Where("slug=?", v.Slug).First(&models.JsonMapping{}).Error; err != nil {
	// 		if err := Gdb.Create(v).Error; err != nil {
	// 			fmt.Println("FATAL", err)
	// 			panic(err)
	// 		}
	// 	}
	// }

	games := LoadGormGames()
	for _, v := range games {
		if err := Gdb.Where("slug=?", v.Slug).First(&models.GormGame{}).Error; err != nil {
			if err := Gdb.Create(v).Error; err != nil {
				fmt.Println("FATAL", err)
				panic(err)
			}
		}
	}

	posts := LoadGormPosts(Gdb)
	for _, v := range posts {
		if err := Gdb.Where("slug=?", v.Slug).First(&models.GormPost{}).Error; err != nil {
			if err := Gdb.Create(v).Error; err != nil {
				fmt.Println("FATAL", err)
				panic(err)
			}
		}
	}

	revs := LoadReviews(posts)
	for _, v := range revs {
		if err := Gdb.Where("path=?", v.Path).First(&models.Review{}).Error; err != nil {
			if err := Gdb.Create(v).Error; err != nil {
				fmt.Println("FATAL", err)
				panic(err)
			}
		}
	}

}
