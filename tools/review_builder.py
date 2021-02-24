#!/bin/python
PATH_PREFIX = "/opt/gametime/reviews/"
FILENAME = "paradise_killer"

class Segment():
    question = ""
    answer = ""
    
    def __init__(self, question="", answer=""):
        self.question = question
        self.answer = answer
    
class Review():
    art = [[],[],[],[]]
    game = [[],[],[],[]]    

def fill_list(lst, line):
    row_list = line.split(',')
    col_ctr = 0
    if len(row_list) != 8:
        print("problem")
        print(len(row_list))
        exit(1)
    for i in lst:
        segment = Segment(row_list[col_ctr], row_list[col_ctr + 1])
        i.append(segment)
        col_ctr = col_ctr + 2
    
def parse_csv ():
    f = open(PATH_PREFIX + FILENAME + ".csv", "r")
    r = Review()
    section = 0
    file_list = f.readlines()
    f.close()
    for line in file_list:
        if section == 0:
            if line.startswith(","):
                section = 1
            if not line.startswith("GRAPHICS") and not line.startswith("ART"):
                fill_list(r.art, line)
                #print("{},{},{},{}".format(len(r.art[0]), len(r.art[1]), len(r.art[2]), len(r.art[3])))
        if section == 1:
            if not line.startswith("MECHANICS") and not line.startswith("GAME"):
                fill_list(r.game, line)
                #print("{},{},{},{}".format(len(r.game[0]), len(r.game[1]), len(r.game[2]), len(r.game[3])))
    return r

def fill_final_list(final_list, review_list, name, title):
    final_list.append("<{}><title>{}</title>\n".format(name, title))
    for i in review_list:
        if not i.answer.startswith("--") and i.answer:
            final_list.append("<question>{}</question>\n".format(i.question.replace('&', ',')))
            final_list.append("<answer>{}</answer>\n".format(i.answer.replace('&', ',').strip()))
    final_list.append("</{}>\n".format(name))

def print_xml(review):
    f = open(PATH_PREFIX + FILENAME + ".xml", "w")
    final_list = ["<review>\n"]
    fill_final_list(final_list, review.game[3], "overall", "Overall")
    final_list.append("<art><title>Art:</title>\n")
    fill_final_list(final_list, review.art[0], "graphics", "Graphics")
    fill_final_list(final_list, review.art[1], "sound", "Sound")
    fill_final_list(final_list, review.art[2], "story", "Story")
    fill_final_list(final_list, review.art[3], "theme", "Themes")
    final_list.append("</art>\n")
    final_list.append("<game><title>Game:</title>\n")
    fill_final_list(final_list, review.game[0], "mechanics", "Mechanics")
    fill_final_list(final_list, review.game[1], "difficulty", "Difficulty")
    fill_final_list(final_list, review.game[2], "experience", "Experience")
    final_list.append("</game>\n")
    final_list.append("</review>")
    f.writelines(final_list)
    f.close()

def remove_trailing_comma(line):
    if line.endswith(",\n"):
        line = line[:-2] + "\n"
        return line

def fill_final_list_json(final_list, review_list, name, title, extra_tab=True):
    tab_str = ["\t", "\t\t", "\t\t\t"]
    if extra_tab:
        tab_str = ["\t\t", "\t\t\t", "\t\t\t\t"]
    final_list.append(tab_str[0] + '"'+ name + '": {\n' + tab_str[1] + '"title":"' + title + '",\n' + tab_str[1] +'"questions": [\n')
    for i in review_list:
        if not i.answer.startswith("--") and len(i.answer) > 1:
            final_list.append(tab_str[2] + '{' + '"question":"{}",\n'.format(i.question.replace('&', ',').replace('"',"'")))
            final_list.append(tab_str[2] + '"answer":"{}"'.format(i.answer.replace('&', ',').replace('"',"'").strip()) + '},\n')
    final_list[-1] = remove_trailing_comma(final_list[-1])
    final_list.append(tab_str[1] + ']\n' + tab_str[0] + '},\n')

def print_json(review):
    f = open(PATH_PREFIX + FILENAME + ".json", "w")
    final_list = ['{\n']
    fill_final_list_json(final_list, review.game[3], "overall", "Overall", False)
    final_list.append('\t"art":{\n\t\t"title":"Art:",\n')
    fill_final_list_json(final_list, review.art[0], "graphics", "Graphics")
    fill_final_list_json(final_list, review.art[1], "sound", "Sound")
    fill_final_list_json(final_list, review.art[2], "story", "Story")
    fill_final_list_json(final_list, review.art[3], "theme", "Themes")
    final_list[-1] = remove_trailing_comma(final_list[-1])
    final_list.append("\t},\n")
    final_list.append('\t"game":{\n\t\t"title":"Game:",\n')
    fill_final_list_json(final_list, review.game[0], "mechanics", "Mechanics")
    fill_final_list_json(final_list, review.game[1], "difficulty", "Difficulty")
    fill_final_list_json(final_list, review.game[2], "experience", "Experience")
    final_list[-1] = remove_trailing_comma(final_list[-1])
    final_list.append("\t}\n")
    final_list.append('\t"pull":""\n')
    final_list.append('}')
    f.writelines(final_list)
    f.close()

print_json(parse_csv())