from django.db import models
from django.contrib.auth.models import User
import uuid

POST_STATUS = (
	(0,"Draft"),
	(1,"Publish")
)

GAME_STATUS = (
	(0, "Unknown"),
	(1, "Acquired"),
	(2, "Started"),
	(3, "Completed"),
	(4, "Reviewed"),
	(5, "Suggested"),
	(6, "Published")
)

PLATFORM = (
	(0, "PC"),
	(1, "GameBoy"),
	(2, "PlayStation"),
	(3, "Xbox"),
	(4, "Nintendo DS"),
	(5, "Wii"),
	(6, "Switch"),
	(7, "Mobile")
)


REVIEW_SEGMENTS = (
	(0, "Overview"),
	(1, "Graphics"),
	(2, "Gameplay"),
	(3, "Music/Sound"),
	(4, "Difficulty"),
	(5, "Story"),
	(6, "Themes,Tropes"),
	(7, "Experience/Bugs"),
	(8, "Overall")
)

class FAQ(models.Model):
	question = models.TextField(max_length=500, unique=True)
	answer = models.TextField(max_length=500, unique=True)

	def __str__(self):
		return str(self.id)

class BlogTag(models.Model):
	tag_id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
	name = models.CharField(max_length=20, unique=True)

	def __str__(self):
		return self.name

class MailingList(models.Model):
	email = models.EmailField(max_length=50)
	name = models.CharField(max_length=50, default="")

class Game(models.Model):
	game_id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
	title = models.TextField(max_length=50, unique=True)
	creator = models.TextField(max_length=50, unique=False)
	publisher = models.TextField(max_length=50, unique=False)
	release_date = models.DateField()
	steam_link = models.URLField(max_length=200)
	status = models.IntegerField(choices=GAME_STATUS, default=0)
	platform = models.IntegerField(choices=PLATFORM, default=0)

	def __str__(self):
		return self.title

	def __id__(self):
		return self.game_id

class Review(models.Model):
	review_id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
	#post_id = models.OneToOneField(Post, on_delete=models.CASCADE)
	title = models.TextField(max_length=50, unique=True)
	status = models.IntegerField(choices=POST_STATUS, default=0)
	game_id = models.ForeignKey(Game, on_delete=models.CASCADE)
	overall_rating = models.IntegerField()
	
	def __str__(self):
		return self.title

	def __id__(self):
		return self.review_id

class Post(models.Model):
	post_id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
	title = models.TextField(max_length=200, unique=True)
	slug = models.SlugField(max_length=200, unique=True)
	author = models.ForeignKey(User, on_delete= models.CASCADE,related_name='blog_posts')
	updated_on = models.DateTimeField(auto_now= True)
	content = models.TextField()
	created_on = models.DateTimeField(auto_now_add=True)
	status = models.IntegerField(choices=POST_STATUS, default=0)
	review_id = models.ForeignKey(Review, on_delete=models.CASCADE)
	game_id = models.ForeignKey(Game, on_delete=models.CASCADE)

	class Meta:
		ordering = ['-created_on']

	def __str__(self):
		return self.title

	def __id__(self):
		return self.post_id



