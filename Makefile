include .env

define check_arg
	@[ ${1} ] || ( echo ">> ${2} is not set, use: ${2}=value"; exit 1 )
endef

CMD_CD_APP_PATH='cd ${APP_PATH}'
CMD_BUILD_IMAGE='${CMD_CD_APP_PATH} && docker build --force-rm --no-cache --build-arg appname=${APP_NAME} -t ${APP_NAME} .'
CMD_RUN_CONTAINER='${CMD_CD_APP_PATH} && docker run --restart always -v $${PWD}/config:/app/config -v $${PWD}/data:/app/data --name ${APP_NAME} -d ${APP_NAME}'
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
	@scp ./config/config.yaml ${USER}@${HOST}:${APP_PATH}/config

## Upload tube config
upload.tube.config:
	$(call check_arg, ${tube}, tube)
	@scp ./config/${tube}.config.yaml ${USER}@${HOST}:${APP_PATH}/config

# Create data, config folders
create.data.and.config:
	@ssh ${USER}@${HOST} mkdir -p ${APP_PATH}/data
	@ssh ${USER}@${HOST} mkdir -p ${APP_PATH}/config

## Upload all app needed files
upload.all: mkdir.apppath
	@make upload.bin
	@make upload.dockerfile
	@make create.data.and.config
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
