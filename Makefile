include .env

CMD_CD_APP_PATH='cd ${APP_PATH}'
CMD_BUILD_IMAGE='${CMD_CD_APP_PATH} && docker build --build-arg appname=${APP_NAME} -t ${APP_NAME} .'
CMD_RUN_CONTAINER='${CMD_CD_APP_PATH} && docker run -v $${PWD}:/app --name ${APP_NAME} -d ${APP_NAME}'
CMD_RESTART_CONTAINER='docker restart ${APP_NAME}'
CMD_STOP_CONTAINER='docker stop ${APP_NAME}'
CMD_LOGS_CONTAINER='docker logs -f --tail=100 ${APP_NAME}'

appname:
	@echo Application named: ${APP_NAME}

run: appname build.local
	@./${APP_NAME}

clean:
	@rm -rf ${APP_NAME}
	@echo clean: OK!

clean.db:
	@rm data data.idx
	@echo clean.db: OK!

build.local: clean
	@go build -o ${APP_NAME}
	@echo build.local: OK!

build.linux: clean
	@GOOS=linux GOARCH=amd64 go build -o ${APP_NAME}
	@echo build.linux: OK!

# Remote rules below

## Create application directory on remote production server
mkdir.apppath:
	@ssh ${USER}@${HOST} mkdir -p ${APP_PATH} 

## Upload app binary file
upload.bin: build.linux
	@scp ./${APP_NAME} ${USER}@${HOST}:${APP_PATH} 

## Upload dockerfile
upload.dockerfile:
	@scp ./Dockerfile ${USER}@${HOST}:${APP_PATH} 

## Upload app config
upload.config:
	@scp ./config.yaml ${USER}@${HOST}:${APP_PATH}

## Upload all app needed files
upload.all: mkdir.apppath
	@make upload.bin
	@make upload.dockerfile
	@make upload.config
	@echo upload: OK!

## Build docker image remotely
remote.build.docker.image:
	@ssh ${USER}@${HOST} ${CMD_BUILD_IMAGE}

## Run container remotely
remote.run.docker.conatainer:
	@ssh ${USER}@${HOST} ${CMD_RUN_CONTAINER}

## Stop container remotely
remote.stop.docker.conatainer:
	@ssh ${USER}@${HOST} ${CMD_STOP_CONTAINER}

## Restart container remotely
remote.restart.docker.conatainer:
	@ssh ${USER}@${HOST} ${CMD_RESTART_CONTAINER}

## View logs for the container
remote.logs.docker.conatainer:
	@ssh ${USER}@${HOST} ${CMD_LOGS_CONTAINER}
