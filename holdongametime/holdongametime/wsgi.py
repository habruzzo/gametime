"""
WSGI config for holdongametime project.

It exposes the WSGI callable as a module-level variable named ``application``.

For more information on this file, see
https://docs.djangoproject.com/en/3.1/howto/deployment/wsgi/
"""

import os
import sys

sys.path.append('/opt/holdongametime')
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'holdongametime.settings')
os.environ.setdefault('PYTHONPATH', '/opt/holdongametime')


import django

from django.core.wsgi import get_wsgi_application
#from django.contrib.auth.handlers.modwsgi import check_password


application = get_wsgi_application()
