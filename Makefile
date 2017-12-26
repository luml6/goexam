OUTDIR=.

APP_NAME=notepad-api

all:$(APP_NAME)

$(APP_NAME): 
	go build -o $@

run:
	nohup ./$(APP_NAME) 2>&1 > $(APP_NAME).log &
	
kill:
	killall $(APP_NAME)
	
clean:
	rm -f $(APP_NAME)
