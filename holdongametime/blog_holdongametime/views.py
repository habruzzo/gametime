from django.views import generic

from .models import Post
from .models import Game
from .models import Review

class PostList(generic.ListView):
    queryset = Post.objects.filter(status=1).order_by('-created_on')
    template_name = 'posts.html'

class PostDetail(generic.DetailView):
    model = Post
    template_name = 'single.html'

class MultipleModelView(generic.TemplateView):
    def get_context_data(self, **kwargs):
        context = super(MultipleModelView, self).get_context_data(**kwargs)
        context['post_list'] = Post.objects.all()[:3]
        return context

class AboutView(MultipleModelView):
	model = Post
	template_name = 'about.html'

class BacklogView(MultipleModelView):
	model = Game
	queryset = Game.objects.order_by('status')
	template_name = 'backlog.html'

class ContactView(MultipleModelView):
	model = Post
	template_name = 'contact.html'

class FormatView(MultipleModelView):
	model = Review
	template_name = 'format.html'

class HomeView(MultipleModelView):
	model = Post
	template_name = 'index.html'

