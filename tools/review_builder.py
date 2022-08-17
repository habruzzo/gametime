#!/bin/python
import re, os, string
CSV_PREFIX = "/opt/gametime/reviews/csv/"
JSON_PREFIX = "/opt/gametime/reviews/json/"
XML_PREFIX = ""

class Segment():
    question = ""
    answer = ""
    
    def __init__(self, question="", answer=""):
        self.question = question
        self.answer = answer
    
    def print_shit(self):
        if self.question == "" or self.answer == "":
            return
        s = "{} {}".format(self.question, self.answer)
        print(s.strip())

class Review():
    art = [[],[],[],[]]
    game = [[],[],[],[]]
    title = ""
    slug = ""
    author = ""

    def __init__(self):
        art = [[],[],[],[]]
        game = [[],[],[],[]]
        title = ""
        slug = ""
        author = ""

    def print_shit(self):
        for i in self.art:
            for j in i:
                j.print_shit()    
        for i in self.game:
            for j in i:
                j.print_shit()

def fill_list(lst, line):
    row_list = line.split(',')
    col_ctr = 0
    if len(row_list) != 8:
        print("problem")
        print(line)
        print(len(row_list))
        exit(1)
    for i in lst:
        segment = Segment(row_list[col_ctr], row_list[col_ctr + 1])
        i.append(segment)
        col_ctr = col_ctr + 2
    
def parse_csv (r, filename):
    f = open(CSV_PREFIX + filename + ".csv", "r")
    section = 0
    file_list = []
    f.seek(0)
    file_list = f.readlines()
    f.seek(0)
    f.close()
    comma = re.compile(",,,,,,,(.*)")
    r.art = [[],[],[],[]]
    r.game = [[],[],[],[]]
    for line in file_list:
        if section <= 0:
            if comma.match(line):
                section += 1
            if line.startswith("game title"):
                parse_title_row(line, r)
                section -= 1
                continue
            if not line.startswith("GRAPHICS") and not line.startswith("ART"):
                fill_list(r.art, line)
                #print("{},{},{},{}".format(len(r.art[0]), len(r.art[1]), len(r.art[2]), len(r.art[3])))
        if section == 1:
            if not line.startswith("MECHANICS") and not line.startswith("GAME"):
                fill_list(r.game, line)
                #print("{},{},{},{}".format(len(r.game[0]), len(r.game[1]), len(r.game[2]), len(r.game[3])))
    del file_list
    return

def parse_title_row(line, review):
    item = line.split(",")
    review.title = item[1]
    review.slug = item[3]
    review.author = item[5]

def parse_csv_ho (r, filename):
    f = open(CSV_PREFIX + filename + ".csv", "r")
    section = 0
    file_list = []
    file_list.clear()
    file_list = f.readlines()
    f.close()
    r.art = [[],[],[],[]]
    r.game = [[],[],[],[]]
    comma = re.compile(",,,,,,,(.*)")
    for line in file_list:
        if section <= 0:
            if comma.match(line):
                #print("end of part 1")
                section += 1
            if line.startswith("game title"):
                parse_title_row(line, r)
                section -= 1
                continue
            if not line.startswith("GRAPHICS") and not line.startswith("ART"):
                fill_list(r.art, line)
                #print("{},{},{},{}".format(r.art[0], r.art[1], r.art[2], r.art[3]))
        if section == 1:
            if not line.startswith("MECHANICS") and not line.startswith("GAME"):
                fill_list(r.game, line)
                #print("{},{},{},{}".format(r.game[0], r.game[1], r.game[2], r.game[3]))
    #r.print_shit()
    return

# def fill_final_list(final_list, review_list, name, title):
#     final_list.append("<{}><title>{}</title>\n".format(name, title))
#     for i in review_list:
#         if not i.answer.startswith("--") and i.answer:
#             final_list.append("<question>{}</question>\n".format(i.question.replace('&', ',')))
#             final_list.append("<answer>{}</answer>\n".format(i.answer.replace('&', ',').strip()))
#     final_list.append("</{}>\n".format(name))

# def print_xml(review, slug):
#     f = open(XML_PREFIX + slug + ".xml", "w")
#     final_list = ["<review>\n"]
#     fill_final_list(final_list, review.game[3], "overall", "Overall")
#     final_list.append("<art><title>Art:</title>\n")
#     fill_final_list(final_list, review.art[0], "graphics", "Graphics")
#     fill_final_list(final_list, review.art[1], "sound", "Sound")
#     fill_final_list(final_list, review.art[2], "story", "Story")
#     fill_final_list(final_list, review.art[3], "theme", "Themes")
#     final_list.append("</art>\n")
#     final_list.append("<game><title>Game:</title>\n")
#     fill_final_list(final_list, review.game[0], "mechanics", "Mechanics")
#     fill_final_list(final_list, review.game[1], "difficulty", "Difficulty")
#     fill_final_list(final_list, review.game[2], "experience", "Experience")
#     final_list.append("</game>\n")
#     final_list.append("</review>")
#     f.writelines(final_list)
#     f.close()

def remove_trailing_comma(line):
    if line.endswith(",\n"):
        line = line[:-2] + "\n"
        return line

def fill_final_list_json(final_list, review_list, name, title, extra_tab=True):
    tab_str = ["", "\t", "\t\t"]
    if extra_tab:
        tab_str = ["\t", "\t\t", "\t\t\t"]
    final_list.append(tab_str[1] +'"'+title+'": [\n')
    for i in review_list:
        if not i.answer.startswith("--") and len(i.answer) > 1:
            final_list.append(tab_str[2] + '{' + '"question":"{}",\n'.format(i.question.replace('&', ',').replace('"',"'")))
            final_list.append(tab_str[2] + '"answer":"{}"'.format(i.answer.replace('&', ',').replace('"',"'").strip()) + '},\n')
    final_list[-1] = remove_trailing_comma(final_list[-1])
    final_list.append(tab_str[1] + '],\n')

def print_json(review):
    final_list = ['{\n']
    fill_final_list_json(final_list, review.game[3], "overall", "Overall", False)
    final_list.append('\t"Art": {\n')
    fill_final_list_json(final_list, review.art[0], "graphics", "Graphics")
    fill_final_list_json(final_list, review.art[1], "sound", "Sound")
    fill_final_list_json(final_list, review.art[2], "story", "Story")
    fill_final_list_json(final_list, review.art[3], "theme", "Themes")
    final_list[-1] = remove_trailing_comma(final_list[-1])
    final_list.append("\t},\n")
    final_list.append('\t"Gameplay": {\n')
    fill_final_list_json(final_list, review.game[0], "mechanics", "Mechanics")
    fill_final_list_json(final_list, review.game[1], "difficulty", "Difficulty")
    fill_final_list_json(final_list, review.game[2], "experience", "Experience")
    final_list[-1] = remove_trailing_comma(final_list[-1])
    final_list.append("\t},\n")
    final_list.append('\t"pull":"SET_QUOTE",\n')
    final_list.append('\t"game":{"title":"' + review.title + '"},\n')
    final_list.append('\t"author":{"name":"' + review.author + '"},\n')
    final_list.append('\t"slug":"' + review.slug + '",\n')
    final_list.append('\t"imgs":["","",""]\n')
    final_list.append('}')

    f = open(JSON_PREFIX + review.slug.lower() + ".json", "w")
    f.writelines(final_list)
    f.close()

def fill_final_list_json_ho(final_list, review_list, name, title, another_coming=True, extra_tab=True):
    tab_str = ["", "\t", "\t\t"]
    if extra_tab:
        tab_str = ["\t", "\t\t", "\t\t\t"]
    final_list.append(tab_str[1] +'"'+ title+'": [\n')
    for i in review_list:
        if not i.answer.startswith("--") and i.answer.strip() != "":
            final_list.append(tab_str[2] + '{' + '"question":"{}",\n'.format(i.question.replace('&', ',').replace('"',"'")))
            final_list.append(tab_str[2] + '"answer":"{}"'.format(i.answer.replace('&', ',').replace('"',"'").strip()) + '},\n')
    final_list[-1] = remove_trailing_comma(final_list[-1])
    if another_coming:
        final_list.append(tab_str[1] + '],\n')
    else:
        final_list.append(tab_str[1] + ']\n')


def print_json_ho(review):
    final_list = ['{\n']
    fill_final_list_json_ho(final_list, review.game[3], "overall", "Overall", True, False)
    final_list.append('\t"Art": {\n')
    fill_final_list_json_ho(final_list, review.art[0], "graphics", "Graphics")
    fill_final_list_json_ho(final_list, review.art[1], "sound", "Sound")
    fill_final_list_json_ho(final_list, review.art[2], "story", "Story", False)
    final_list.append("\t},\n")
    final_list.append('\t"Gameplay": {\n')
    fill_final_list_json_ho(final_list, review.game[1], "difficulty", "Difficulty")
    fill_final_list_json_ho(final_list, review.game[2], "experience", "Experience", False)
    final_list.append("\t},\n")
    final_list.append('\t"pull":"SET_QUOTE",\n')
    final_list.append('\t"game":{"title":"' + review.title + '"},\n')
    final_list.append('\t"author":{"name":"' + review.author + '"},\n')
    final_list.append('\t"slug":"' + review.slug + '",\n')
    final_list.append('\t"imgs":["","",""]\n')
    final_list.append('}')

    f = open(JSON_PREFIX + review.slug.lower() + ".json", "w")
    f.writelines(final_list)
    f.close()

def run_all():
    file_set = os.getenv("FILE_SET")
    print(file_set)
    if file_set is None:
        print("hey!! fill out FILE_SET!!")
        print('''ex: FILE_SET=`ls reviews/csv | cut -d '.' -f1 | tr "\n" " "` ''')
        # os.system('for i in `find reviews/csv -type f -regex ".* .*"| wc -l | xargs seq`; do NAME=`find reviews/csv -type f -regex ".* .*" | tr "\n" "\t" | cut -f $i`; NEWNAME=`echo $NAME | tr " " "_"`; echo $NAME $NEWNAME; mv $NAME $NEWNAME; done')
        return
    fileset = file_set.split()
    review_list = []
    for i in range(len(fileset)):
        review_list.append(Review())
    print(review_list)
    for index, filename in enumerate(fileset):
        print("##########################",filename)
        if "_h_o" in filename:
            print(index)
            parse_csv_ho(review_list[index], filename)
            print_json_ho(review_list[index])
        else:
            print(index)
            parse_csv(review_list[index], filename)
            print_json(review_list[index], filename)

def main():
    ra = input("run all reviews?(y/n, default n): ")
    if ra == "y":
        run_all()
        exit(0)
    filename = input("file name: ")
    ho = input("hidden object(y/n, default n): ")
    if ho == "y":
        r = Review()
        parse_csv_ho(r, filename)
        print_json_ho(r)
        exit(0)
    else:
        parse_csv(r, filename)
        print_json(r)
        exit(0)
    

if __name__ == "__main__":
    main()