from django.views import generic
from .models import Post

class PostList(generic.ListView):
    queryset = Post.objects.filter(status=1).order_by('-created_on')
    template_name = 'posts.html'

class PostDetail(generic.DetailView):
    model = Post
    template_name = 'single.html'

class AboutView(generic.ListView):
	model = Post
	template_name = 'about.html'

class BacklogView(generic.ListView):
	model = Game
	queryset = Game.objects.order_by('status')
	template_name = 'backlog.html'

class ContactView(generic.ListView):
	model = Post
	template_name = 'contact.html'

class FormatView(generic.ListView):
	model = Review
	template_name = 'format.html'

class HomeView(generic.ListView):
	model = Post
	template_name = 'index.html'