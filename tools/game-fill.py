#tool that will take this list -> collect data to fill out "game" object -> fill out "game" object -> add to "game" table
#given a list of game titles, collect "publisher" "creator" "release date" "steam link"
#read from html
import bs4
import mechanicalsoup
import urllib.parse
import pudb
import sys

KEY_NAMES=["title","creator","publisher","release","platform","steam_link","status"]

def read_title_list():
	with open("game_list.xml") as f:
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
	if ps[0].text.lower().split(":")[0] == "platform":
		platform = ps[0].text.split(":")[1]
	if ps[1].text.lower().split(":")[0] == "developer(s)":
		creator = ps[1].text.split(":")[1]
	if ps[2].text.lower().split(":")[0] == "publishers(s)":
		publisher = ps[2].text.split(":")[1]
	if ps[3].text.lower().split(":")[0] == "releasedate":
		release = ps[3].text.split(":")[1]
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
		print("release= " + game["release"])
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
			k = input("enter key:0=title,1=creator,2=publisher,3=release,4=platform,5=steam_link,6=status,7=exit")
			if k != "7":
				v = input("enter value for {}, or enter # to exit".format(k))
				if v != "#":
					game[KEY_NAMES[k]] = v
	report_game_info(game)		
	return game

def login(browser):
	browser.open("http://127.0.0.1:8000/admin")
	username="bot"
	password="bot"
	browser.select_form('form[id="game_form"')
	return browser

def enter_game_form(browser, game):
	browser.open("http://127.0.0.1:8000/admin/blog_holdongametime/game/add/")
	i = 0
	browser["id_title"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["id_creator"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["id_publisher"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["id_release"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["id_platform"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["id_steam_link"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["id_status"] = game[KEY_NAMES[i]]
	i = i + 1
	response = browser.submit_selected()
	print(response.text)
	return browser

def main():
	if len(sys.argv) > 1 and sys.argv[1] == "build":
		print("building game list")
	else:
		print("scraping and inserting game info")
		title_list = read_title_list()
		browser = mechanicalsoup.StatefulBrowser()
		#login(browser)
		for title in title_list:
			print(title.text)
			game = research_basic_game(browser, title.text)
			report_game_info(game)
			game = query_continue(game)
			#enter_game_form(browser, game)


if __name__ == '__main__':
	main()