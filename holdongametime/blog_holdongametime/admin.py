from django.contrib import admin
from .models import Post
from .models import Game
from .models import Review

class PostAdmin(admin.ModelAdmin):
    list_display = ('title', 'slug', 'status','created_on')
    list_filter = ("status",)
    search_fields = ['title', 'content']
    prepopulated_fields = {'slug': ('title',)}
  
class GameAdmin(admin.ModelAdmin):
	list_display = ('title', 'status')
	list_filter = ('status',)

class ReviewAdmin(admin.ModelAdmin):
	list_display = ('title', 'post_id', 'game_id', 'overall_rating')

admin.site.register(Post, PostAdmin)
admin.site.register(Game, GameAdmin)
admin.site.register(Review, ReviewAdmin)