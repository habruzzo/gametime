from django.contrib import admin
from .models import Post
from .models import Game
from .models import Review
from .models import FAQ
from .models import BlogTag


class PostAdmin(admin.ModelAdmin):
    list_display = ('title', 'slug', 'status','created_on')
    list_filter = ("status",)
    search_fields = ['title', 'content']
    prepopulated_fields = {'slug': ('title',)}
  
class GameAdmin(admin.ModelAdmin):
	list_display = ('title', 'status')
	list_filter = ('status',)
	search_field = ['title', 'publisher']

class ReviewAdmin(admin.ModelAdmin):
	list_display = ('title', 'status', 'game_id', 'overall_rating')
	list_filter = ("status",)
	search_field = ['title', 'content']

class FAQAdmin(admin.ModelAdmin):
	list_display = ('id', 'question')
	search_field = ['question', 'answer']

class BlogTagAdmin(admin.ModelAdmin):
	list_display = ('tag_id', 'name')
	search_field = ['name']

admin.site.register(Post, PostAdmin)
admin.site.register(Game, GameAdmin)
admin.site.register(Review, ReviewAdmin)
admin.site.register(FAQ, FAQAdmin)
admin.site.register(BlogTag, BlogTagAdmin)