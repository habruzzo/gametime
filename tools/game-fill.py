#!/bin/python

#tool that will take this list -> collect data to fill out "game" object -> fill out "game" object -> add to "game" table
#given a list of game titles, collect "publisher" "creator" "release date" "steam link"
#read from html

import bs4
import mechanicalsoup
import urllib.parse
import pudb
import sys
import json

PATH_PREFIX = "/opt/gametime/reviews/"
KEY_NAMES = ["title","creator","publisher","release_date","platform","steam_link","status"]
PLATFORMS = ["pc", "gameboy", "playstation", "xbox", "wii", "switch", "mobile", "other"]
IP = "52.37.133.146"

def normalize_platform(platform):
	plat = platform.strip().lower()
	for p in PLATFORMS:
		if p in plat:
			return PLATFORMS.index(p)
	return 8

def read_title_list():
	f = open("game_list.xml")
	soup = bs4.BeautifulSoup(f, "xml")
	print(soup.text)
	title_list = soup.find_all("title")
	return title_list

def sanitize(title):
	return urllib.parse.quote(bytes(title, "utf-8"))

def scrape_basic_data(browser):
	card = browser.page.find_all("div","card-body")[1]
	ps=card.find_all("p")
	platform = ""
	creator = ""
	publisher = ""
	release = ""
	try:
		if ps[0].text.lower().split(":")[0] == "platform":
			platform = ps[0].text.split(":")[1]
		if ps[1].text.lower().split(":")[0] == "developer(s)":
			creator = ps[1].text.split(":")[1]
		if ps[2].text.lower().split(":")[0] == "publishers(s)":
			publisher = ps[2].text.split(":")[1]
		if ps[3].text.lower().split(":")[0] == "releasedate":
			release = ps[3].text.split(":")[1]
	except IndexError:
		pass
	return platform, creator, publisher, release


def research_basic_game(browser, title):
	title_cl=sanitize(title)
	browser.open("https://thegamesdb.net/search.php?name={}".format(title))
	game = {}
	link = browser.find_link('/game.php')
	browser.follow_link(link)
	platform, creator, publisher, release = scrape_basic_data(browser)
	game[KEY_NAMES[0]] = title.strip()
	game[KEY_NAMES[1]] = creator.strip()
	game[KEY_NAMES[2]] = publisher.strip()
	game[KEY_NAMES[3]] = release.strip()
	game[KEY_NAMES[4]] = platform.strip()
	return game

def report_game_info(game):
	if "title" in game.keys():
		print("title= " + game["title"])
	if "creator" in game.keys():
		print("creator= " + game["creator"])
	if "publisher" in game.keys():
		print("publisher= " + game["publisher"])
	if "release" in game.keys():
		print("release= " + game["release_date"])
	if "platform" in game.keys():
		print("platform= " + game["platform"])
	if "steam_link" in game.keys():
		print("steam_link= " + game["steam_link"])
	if "status" in game.keys():
		print("status= " + game["status"])

def query_continue(game):
	info_correct = input("is this game information correct?y/n")
	if info_correct =="y":
		if "status" not in game.keys():
			game["status"] = input("what is the status of this game?0=unknown,1=acquired,2=started,3=completed,4=reviewed,5=suggested,6=published")
	else:
		b = input("would you like to enter manually or skip for now?e=enter,s=skip")
		if b !="e":
			game = {}
		else:
			k = input("enter key:0=title,1=creator,2=publisher,3=release_date,4=platform,5=steam_link,6=status,7=exit")
			if k != "7":
				v = input("enter value for {}, or enter # to exit".format(k))
				if v != "#":
					game[KEY_NAMES[k]] = v
	report_game_info(game)		
	return game

def login(browser):
	url = "http://{}/admin".format(IP)
	browser.open(url)
	username="bot"
	password="bot"
	browser.select_form("#login-form")
	browser["username"] = username
	browser["password"] = password
	response = browser.submit_selected()
	print(response)
	return browser

def enter_game_form(browser, game):
	url = "http://{}/admin/blog_holdongametime/game/add/".format(IP)
	print(browser.open(url))
	form = browser.select_form("#game_form")
	i = 0
	browser["title"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["creator"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["publisher"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["release_date"] = game[KEY_NAMES[i]]
	i = i + 1
	platform = normalize_platform(game[KEY_NAMES[i]])
	print(platform)
	browser["platform"] = platform
	i = i + 1
	browser["steam_link"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["status"] = game[KEY_NAMES[i]]
	i = i + 1
	form.choose_submit("_add_another")
	response = browser.submit_selected()
	print(response.text)
	return browser

def main_nope():
	if len(sys.argv) > 1 and sys.argv[1] == "build":
		print("building game list")
	else:
		print("scraping and inserting game info")
		title_list = read_title_list()
		game_list = []
		browser = mechanicalsoup.StatefulBrowser()
		for title in title_list:
			print(title.text)
			game = research_basic_game(browser, title.text)
			report_game_info(game)
			game = query_continue(game)
			game_list.append(game)
		print(game_list)
		browser = login(browser)
		for game in game_list:
			try:
				enter_game_form(browser, game)
			except KeyError:
				print()
				print(browser.page)
				print(game_list)

def main():
	if len(sys.argv) > 1 and sys.argv[1] == "build":
		print("building game list")
	else:
		print("scraping and inserting game info")
		title_list = read_title_list()
		game_list = []
		browser = mechanicalsoup.StatefulBrowser()
		for title in title_list:
			print(title.text)
			game = research_basic_game(browser, title.text)
			report_game_info(game)
			game = query_continue(game)
			game_list.append(game)
		#print(game_list)
		f = open(PATH_PREFIX + "game_list.json", "w")
		json.dump(game_list,f)
		f.close()


if __name__ == '__main__':
	main()