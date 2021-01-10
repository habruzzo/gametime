from django.contrib import admin
from .models import Post
from .models import Game
from .models import Review
from .models import FAQ


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

class FAQAdmin(admin.ModelAdmin):
	list_display = ('id', 'question')

admin.site.register(Post, PostAdmin)
admin.site.register(Game, GameAdmin)
admin.site.register(Review, ReviewAdmin)
admin.site.register(FAQ, FAQAdmin)