#tool that will take this list -> collect data to fill out "game" object -> fill out "game" object -> add to "game" table
#given a list of game titles, collect "publisher" "creator" "release date" "steam link"
#read from html
import mechanicalsoup

def read_title_list(file_name):
	soup = bs4.BeautifulSoup("game_list.xml"):
	title_list = soup.find_all("title").contents

def login(browser):
	browser.open("http://127.0.0.1:8000/admin")
	username="bot"
	password="bot"
	browser.select_form('form[id="game_form"')
	return browser

def enter_game_form(browser):
	browser.open("http://127.0.0.1:8000/admin/blog_holdongametime/game/add/")
	browser["id_title"] = 
	browser["id_creator"] =
	browser["id_publisher"] =
	browser["id_steam_link"] =
	browser["id_status"] = 
	browser[""]
	return

def main():

	browser = mechanicalsoup.StatefulBrowser()
	enter_game_form(login(browser))
