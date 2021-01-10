#load faq.txt questions and answers into db
import mechanicalsoup

KEY_NAMES = ["question", "answer"]
FAQ_PATH = "../holdongametime/static/content/faq.xml"

def read_faq_list():
	with open(FAQ_PATH) as f:
		soup = bs4.BeautifulSoup(f, "xml")
		print(soup.text)
		faq_list = soup.find_all("question")
		return faq_list

def enter_faq_form(browser, game):
	browser.open("http://127.0.0.1:8000/admin/blog_holdongametime/faq/add/")
	browser.select_form("faq_form")
	i = 0
	browser["id_question"] = game[KEY_NAMES[i]]
	i = i + 1
	browser["id_answer"] = game[KEY_NAMES[i]]
	i = i + 1
	response = browser.submit_selected()
	print(response.text)
	return browser

def main():
	browser = mechanicalsoup.StatefulBrowser()
	faq_list = read_faq_list(FAQ_PATH)
	for faq in faq_list:
		browser = enter_faq_form(browser, faq)

if __name__ == '__main__':
	main()