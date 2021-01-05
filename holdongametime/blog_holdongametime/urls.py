from . import views
from django.urls import path

urlpatterns = [
    path('', views.HomeView.as_view(), name='home'),
    path('about.html', views.AboutView.as_view(), name='about'),
    path('contact.html', views.ContactView.as_view(), name='contact'),
    path('format.html', views.FormatView.as_view(), name='format'),
    path('backlog.html', views.BacklogView.as_view(), name='backlog'),
    path('posts.html', views.PostList.as_view(), name='posts'),
    path('<slug:slug>/', views.PostDetail.as_view(), name='post_detail'),
]