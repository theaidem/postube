# Postube

Parse youtube channels, and post videos to vk.com public page, and telegram channel. 

Example vk public pages: [dnbtube](https://vk.com/dnbtube), [loungetube](https://vk.com/loungetube)

## Get started

### Clone the repo

```git clone git@github.com:theaidem/postube.git project_name```

`project_name` is the name of the project without spaces ex. `onlycats`, `funy_videos`, etc

### Copy and Define project env variables

```cp .env.example .env```

then edit `.env` file

### Copy and Define configuration file for the application

```cp config.example.yaml config.yaml```

 Firstly, let's create a [Standalone app](https://vk.com/editapp?act=create)

 ### Generate access token:

**Opening OAuth Authorization Dialog:**

 ```
 https://oauth.vk.com/authorize?
    client_id=APP_ID
    &redirect_uri=https://vk.com
    &scope=73744
    &response_type=code
    &v=5.92
 ```

 > The **APP_ID** is your standalone application ID

 After successful application authorization, user's browser will be redirected to ``redirect_uri`` URL specified when the authorization dialog box was opened. With that, code to receive code access key will be passed in GET parameter to the specified URL:

 > http://redirect_uri?code=7a6fa4dff77a228eeda56603b8f53806c883f011c40b72630bb50df056f6479e52a

 **Receiving "access_token":**

 ```
 https://oauth.vk.com/access_token? 
    client_id=APP_ID
    &client_secret=APP_SECRET
    &code=CODE
    &redirect_uri=https://vk.com
```

 > The **APP_ID** is your standalone application ID

 > The **APP_SECRET** secret key of the standalone application
 
 > The **CODE** code parameter from the previous step

Example of server response:

```json
{
    "access_token":"533bacf01e11f55b536a565b57531ac114461ae8736d6506a3", "expires_in":43200, 
    "user_id":6492
}
```

Copy ``access_token`` to config.yaml

```yaml
...
    # Token, for access VK API
    accessToken: "access_token"
...
```

More about authorization see [there](https://vk.com/dev/auth_sites)

### Create vk public page

https://vk.com/groups?w=groups_create

Paste your group id to config.yaml

```yaml
...
    # Your public ID
    groupID: "group_id"
...
```

 ### Create telegram bot account and API token

Follow https://telegram.me/BotFather 

Paste your API token to config.yaml

```yaml
...
    # Token your telegram bot
    token: "telegram_api_token"
...
```

 ### Create telegram channel

Paste channel name to config.yaml

```yaml
...
    # Name your chanel, example: "@postube"
    channel: "@name"
...
```

 ### Define your youtube channels ids collection

 Example:

```yaml
# Youtube channels ID's
channels:
  # SpaceX (example!)
  - UCtI0Hodo5o5dUb67FeUjDeA
  # BostonDynamics (example!)
  - UC7vVhkEfw4nOGp8TyDk7RcQ
```

## Start parser, and publish content

```shell
make run
```

## Deployment to server

> This **optional** step, you can use your own approach for that

Edit `.env` vars, `USER` `HOST` `APP_PATH`

See comments for details

### Create path and upload app files on remote server

```shell
make upload.all
```

### Build and run docker container for app

Make sure Docker has been installed on your server 

```shell
make remote.build.docker.image
```

and run it

```shell
make remote.run.docker.conatainer
```

check container logs

```shell
make remote.logs.docker.image
```

### Configure application

You can edit `config.yaml` locally, and upload to server, the application will re-reading the config file automatically

```shell
make upload.config
```

