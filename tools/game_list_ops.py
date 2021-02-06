#!/bin/python
import bs4

def add_annotation(line, annotated_game_list):
    print("game: {}".format(line))
    annotation = input("please add annotation here: ")
    annotated_game_list.append("{}..{}".format(line, annotation))

def annotate ():
    game_list = open('tools/game_list.txt', 'r').readlines()
    annotation_type = input("what kind of annotation are you going to add?")
    print("you said: {}".format(annotation_type))
    annotated_game_list = []
    print(game_list[0])
    annotated_game_list.append("game_list_title..{}".format(annotation_type))
    print(annotated_game_list)
    separation_count = 0
    for line in game_list:
        print(separation_count)
        if "---" in line:
            separation_count = separation_count + 1
            continue
        if separation_count == 1 or separation_count == 3:
            add_annotation(line, annotated_game_list)
        if separation_count == 2 or separation_count == 4:
            continue
    f = open('tools/annotated_game_list.txt', 'a')
    f.writelines(annotated_game_list)
    f.close()

def fill_xml ():
    f1 = open("game_list.xml", "w")
    f2 = open("game_list.txt", "r")
    final_xml = ["<?xml version='1.0' encoding='UTF-8' ?>", "<games>"]
    for line in f2.readlines():
        if "---" in line:
            continue
        print("<game><title>{}</title></game>".format(line.strip()))
        final_xml.append("<game><title>{}</title></game>".format(line.strip()))
    final_xml.append("</games>")
    f1.writelines(final_xml)
    f1.close()

print('''GAME_STATUS = (
	(0, "Unknown"),
	(1, "Acquired"),
	(2, "Started"),
	(3, "Completed"),
	(4, "Reviewed"),
	(5, "Suggested"),
	(6, "Published"),
    (7, "Discarded")
)''')
fill_xml()