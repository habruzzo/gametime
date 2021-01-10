# Generated by Django 3.1.5 on 2021-01-10 07:51

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('blog_holdongametime', '0002_auto_20210110_0727'),
    ]

    operations = [
        migrations.AddField(
            model_name='game',
            name='platform',
            field=models.IntegerField(choices=[(0, 'PC'), (1, 'GameBoy'), (2, 'PlayStation'), (3, 'Xbox'), (4, 'Nintendo DS'), (5, 'Wii'), (6, 'Switch'), (7, 'Mobile')], default=0),
        ),
    ]
