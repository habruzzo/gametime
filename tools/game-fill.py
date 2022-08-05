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
JSON_PREFIX = PATH_PREFIX + "/json/"
SLUG_FILE = "slug_list.txt"
KEY_NAMES = ["title","creator","publisher","release_date","platform","gamesdb_link","status", "imgs", "slug"]
PLATFORMS = ["pc", "gameboy", "playstation", "xbox", "wii", "switch", "mobile", "other"]
IP = "52.37.133.146"

def normalize_platform(platform):
	plat = platform.strip().lower()
	for p in PLATFORMS:
		if p in plat:
			return PLATFORMS.index(p)
	return 8

# def read_title_list():
# 	f = open("game_list.xml")
# 	soup = bs4.BeautifulSoup(f, "xml")
# 	print(soup.text)
# 	title_list = soup.find_all("title")
# 	f.close()
# 	return title_list
def read_title_list():
	f = open("game_list.txt")
	title_list = f.readlines()
	f.close()
	return title_list

def return_difference(full, subset):
	diff = []
	for item in full:
		if not item in subset:
			diff.append(item)
	return diff

def rewrite_title_list(tl):
	f = open("game_list_1.txt", "w")
	f.writelines(tl)
	f.close()

def sanitize(title):
	return urllib.parse.quote(bytes(title, "utf-8"))

def temp_slug(title):
	return title.lower().replace(" ", "_")
	
def read_slug_list():
	f = open(SLUG_FILE)
	slug_list = f.readlines()
	sl = []
	for slug in slug_list:
		sl.append(slug.strip())
	f.close()
	return sl

def rewrite_slug_list(sl):
	f = open(SLUG_FILE, "w")
	f.writelines(sl)
	f.close()

def grab_images(browser):
	il = []
	img_links=browser.page.find_all(class_="fancybox-thumb")
	for img in img_links:
		print(img)
		print(img["href"])
		il.append(img["href"])
	return il

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
	browser.open("https://thegamesdb.net/search.php?name={}".format(title_cl))
	game = {}
	links = browser.links('/game.php')
	pc_link = links[0]
	release2 = ""
	for link in links:
		plat = link.find(class_="bg-secondary").contents[5].text
		print(plat)
		if "PC" in plat:
			pc_link = link
			release2 = link.find(class_="bg-secondary").contents[3].text
			print(release2)
			break
		# if input("is {}this the correct game?y/n".format(link) ) == "n":
		# 	continue
		# else:
		# 	pc_link = link
		# 	release2 = link.find(class_="bg-secondary").contents[3]
		# 	print(release2)
		# 	break

	browser.follow_link(pc_link)
	gamesdb_link = browser.url
	imgs = grab_images(browser)
	platform, creator, publisher, release = scrape_basic_data(browser)
	if not "-" in release.strip():
		release = release2
	game[KEY_NAMES[0]] = title.strip()
	game[KEY_NAMES[1]] = creator.strip()
	game[KEY_NAMES[2]] = publisher.strip()
	game[KEY_NAMES[3]] = release.strip()
	game[KEY_NAMES[4]] = platform.strip()
	game[KEY_NAMES[5]] = gamesdb_link
	game[KEY_NAMES[7]] = imgs
	print(game)
	return game

def report_game_info(game):
	if "title" in game.keys():
		print("title= " + game["title"])
	if "creator" in game.keys():
		print("creator= " + game["creator"])
	if "publisher" in game.keys():
		print("publisher= " + game["publisher"])
	if "release_date" in game.keys():
		print("release= " + game["release_date"])
	if "platform" in game.keys():
		print("platform= " + game["platform"])
	if "steam_link" in game.keys():
		print("steam_link= " + game["steam_link"])
	if "status" in game.keys():
		print("status= " + game["status"])
	if "gamesdb_link" in game.keys():
		print("link= " + game["gamesdb_link"])
	if "imgs" in game.keys():
		print("imgs= {}".format(game["imgs"]))

def query_continue(game, sl):
	info_correct = input("is this game information correct?y/n")
	if info_correct =="y":
		if "status" not in game.keys():
			game["status"] = input("what is the status of this game?0=unknown,1=acquired,2=started,3=completed,4=reviewed,5=suggested,6=published")
		if "slug" not in game.keys():
			s = temp_slug(game["title"])
			if not s in sl:
				if input("is {} a good slug? we think it is not in the slug list already. y/n".format(s)) == "y":
					game["slug"] = s
					sl.append(s)
				else:
					game["slug"] = input("fine then enter a slug now:")

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

# def login(browser):
# 	url = "http://{}/admin".format(IP)
# 	browser.open(url)
# 	username="bot"
# 	password="bot"
# 	browser.select_form("#login-form")
# 	browser["username"] = username
# 	browser["password"] = password
# 	response = browser.submit_selected()
# 	print(response)
# 	return browser

# def enter_game_form(browser, game):
# 	url = "http://{}/admin/blog_holdongametime/game/add/".format(IP)
# 	print(browser.open(url))
# 	form = browser.select_form("#game_form")
# 	i = 0
# 	browser["title"] = game[KEY_NAMES[i]]
# 	i = i + 1
# 	browser["creator"] = game[KEY_NAMES[i]]
# 	i = i + 1
# 	browser["publisher"] = game[KEY_NAMES[i]]
# 	i = i + 1
# 	browser["release_date"] = game[KEY_NAMES[i]]
# 	i = i + 1
# 	platform = normalize_platform(game[KEY_NAMES[i]])
# 	print(platform)
# 	browser["platform"] = platform
# 	i = i + 1
# 	browser["steam_link"] = game[KEY_NAMES[i]]
# 	i = i + 1
# 	browser["status"] = game[KEY_NAMES[i]]
# 	i = i + 1
# 	form.choose_submit("_add_another")
# 	response = browser.submit_selected()
# 	print(response.text)
# 	return browser

# def main_nope():
# 	if len(sys.argv) > 1 and sys.argv[1] == "build":
# 		print("building game list")
# 	else:
# 		print("scraping and inserting game info")
# 		title_list = read_title_list()
# 		game_list = []
# 		browser = mechanicalsoup.StatefulBrowser()
# 		for title in title_list:
# 			print(title.text)
# 			game = research_basic_game(browser, title.text)
# 			report_game_info(game)
# 			game = query_continue(game)
# 			game_list.append(game)
# 		print(game_list)
# 		browser = login(browser)
# 		for game in game_list:
# 			try:
# 				enter_game_form(browser, game)
# 			except KeyError:
# 				print()
# 				print(browser.page)
# 				print(game_list)

def main():
	if len(sys.argv) > 1 and sys.argv[1] == "build":
		print("building game list")
	else:
		print("scraping and inserting game info")
		title_list = read_title_list()
		done_titles = []
		game_list = []
		sl = read_slug_list()
		browser = mechanicalsoup.StatefulBrowser()
		count = 0
		for title in title_list:
			if title.strip() == "---":
				continue
			done_titles.append(title)
			if count == 1:
				break
			print(title.strip())
			game = research_basic_game(browser, title.strip())
			report_game_info(game)
			game = query_continue(game, sl)
			game_list.append(game)
			count = count + 1
		#print(game_list)
		f = open(JSON_PREFIX + "game_list_1.json", "w")
		json.dump(game_list,f)
		f.close()
		rewrite_slug_list(sl)
		rewrite_title_list(return_difference(title_list, done_titles))


if __name__ == '__main__':
	main()