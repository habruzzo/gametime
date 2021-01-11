from django.views import generic

from .models import Post
from .models import Game
from .models import Review
from .models import Bug

class PostList(generic.ListView):
    queryset = Post.objects.filter(status=1).order_by('-created_on')
    template_name = 'posts.html'

class PostDetail(generic.DetailView):
    model = Post
    template_name = 'single.html'
    def get_context_data(self, **kwargs):
    	context = super(PostDetail, self).get_context_data(**kwargs)
    	context['need_tags'] = True
    	return context


class SidebarView(generic.TemplateView):
    def get_context_data(self, **kwargs):
        context = super(SidebarView, self).get_context_data(**kwargs)
        context['post_list'] = Post.objects.all()[:3]
        return context

class BugAndSidebarView(SidebarView):
	def get_context_data(self, **kwargs):
		context = super(BugAndSidebarView, self).get_context_data(**kwargs)
		context['bug_list'] = Bug.objects.all()
		return context

class AboutView(SidebarView):
	model = Post
	template_name = 'about.html'

class BacklogView(BugAndSidebarView):
	model = Game
	queryset = Game.objects.order_by('status')
	template_name = 'backlog.html'

class ContactView(SidebarView):
	model = Post
	template_name = 'contact.html'

class FormatView(SidebarView):
	model = Review
	template_name = 'format.html'

class HomeView(SidebarView):
	model = Post
	template_name = 'index.html'

