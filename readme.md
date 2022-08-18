"holdon Game Time"

holdongametime.com

logo?
draw something on krita nbd

twitter: holdongametime ~
instagram: holdongametime

to-do:



slow

TO-DO:
home finish --find images
contact finish
caddy https
backlog finish
mailing list
ads
write 4 more reviews
tags
comments
social media

V1:
home finish
contact finish
backlog finish
https
data dump and reload tool (done)
ci/cd automated deployments 

V2
mailing list
ads
write 4 more reviews
tags

V3
comments
social media bot

backlog finish:
	-get game list from big fish?
	-fill out games db (migration tool)
	

combination of "personal" blog posts and "reviews"

get dgraph:
curl https://get.dgraph.io -sSf | bash

tool that will get names of all games in my steam library ->
https://store.steampowered.com/account/licenses/
tool that will take this list -> collect data to fill out "game" object -> fill out "game" object -> add to "game" table 

tool to scrape post entry from google sheets rubric
post plainly formatted series of questions and answers, like an interview. so 
SECTION(overall, art, game)
	HEADING(graphics, audio, etc)
	  question:
		my answer

i dont think i need the  database, actually. i need to think about storage/serving content. do i want things to be in the database? or do i want them to be stored on disk/build them each time or only once? 
build once
store online, in private github? and then download them onto memory
tool to build page from rubric (done)
tool to save page into some storage
tool to load articles()

review process:
fill out google sheet
appscript:
trigger on creation of new sheet? trigger on edit of "index sheet" where ill add finished review after its done
appscript will hit public endpoint ping
public endpoint ping will use google workspace api to authorize, grab the sheet, replace the stupid commas, download as csv
replace commas with &
download as csv
parse csv to json
add post to list
load json into db

either add game to json or run game tool

in aws, one machine
caddy docker?
dgraph dockers
app on box

public ip address

"install": 
amazon linux 2
ssh with creds
scp creds and some git creds
install git
git pull everything down
install docker
run the docker containers
load data backup if needed
build and run the application
start caddy
verify everything works

"publish review":
ping endpoint to trigger new review is available, some sort of os flag is set, or a job is run or something
oh yeah i could do golang exec and a make command and the tool will do the import
google workspace api: account auth, replace stupid commas, pull down new sheet csv
run review builder
git push review

github trigger to service
git pull review on box
restart application
application loads all games and reviews
new review is added to db

"deploy":
push to github
aws lambda? or ssh exec
stop service
pull from github
rebuild and re run application

"backups":
every 2 weeks?
ssh exec cron or aws lambda
dump db
git push (size could depend as content grows? ehhh ill be fine for a while)



improve tooling!

Make a UI for editing data?



