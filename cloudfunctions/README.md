# gsuitemdm Cloud Functions #

The core components of the `gsuitemdm` system are deployed as [GCP cloud functions](https://cloud.google.com/functions/docs/). Each cloud function exists to perform a single mobile device-related task (`approve` a device, `search` for a device, `block` a device, etc), and are used extensively by the included command line tool [`mdmtool`](https://github.com/rickt/gsuitemdm/tree/master/mdmtool) as well as being able to be called via `curl` or whatever http(s) library you prefer.

## List of `gsuitemdm` Cloud Functions ##
(where `$CFPREFIX` is the URL prefix of your GCP project, such as `https://us-central1-yourproject.cloudfunctions.net`.

Cloud Function | What the Cloud Function Does | API Endpoint URL
:--- | :--- | :---
 `ApproveDevice`	 | Approves a mobile device 	 | `$CFPREFIX/ApproveDevice`
 `BlockDevice` 	 | Blocks a mobile device	 | `$CFPREFIX/BlockDevice`
 `DeleteDevice`	 | Deletes a mobile device from company MDM	 | `$CFPREFIX/DeleteDevice`
 `Directory`	 | Company phone directory	 | `$CFPREFIX/Directory`
 `SearchDatastore` 	 | Searches Google Datastore for a mobile device	 | `$CFPREFIX/SearchDatastore`
 `SlackDirectory`	 | Company phone directory specifically for Slack	 | `$CFPREFIX/SlackDirectory`
 `UpdateDatastore`	 | Updates a mobile device in Google Datastore with fresh data from the Google Admin SDK	 | `$CFPREFIX/UpdateDatastore`
 `UpdateSheet`	 | Updates the Google Sheet	 | `$CFPREFIX/UpdateSheet`
 `WipeDevice`	 | Wipes a mobile device	 | `$CFPREFIX/WipeDevice`

## Design ##
All of the `gsuitemdm` cloud functions are designed to be as simple as possible, all use the same general design principles and follow [recommended GCP cloud function design principles/best practices](https://cloud.google.com/functions/docs/bestpractices/tips). All `gsuitemdm` cloud functions are super lightweight [http(s)-triggered](https://cloud.google.com/functions/docs/writing/http#writing_http_helloworld-go) mini-webservers. They are deployed to GCP using `gcloud`, and scale up/down as needed. Each cloud function deployment consists of a single `.go` source file and a `.yaml` file containing several environment variables pointing to a shared configuration. 

The basic `gsuitemdm` cloud function model consists of 3 steps:

1. **Basic Checks**
2. **`gsuitemdm` service startup & execution of requested action**
3. **Cleanup**

Let's look at each step in detail:

1. **Basic Checks**
  * https listener starts up, listens for requests
  * Verify incoming requests don't have a null body and appear to be valid JSON for our API
  * Retrieve the GSuiteMDM API key from [Secret Manager](https://cloud.google.com/secret-manager/docs/)
  * Verify that a valid API key was sent in the request
  * Verify that a correct action (specific to each cloud function) was sent in the request
  * Perform basic sanity checks on the action-specific data (specific to each cloud function) that was sent in the request
2. **`gsuitemdm` service startup & execution of requested action**
  * Retrieve the shared GSuiteMDM configuration from [Secret Manager](https://cloud.google.com/secret-manager/docs/)
  * Retrieve all G Suite domain configurations from [Secret Manager](https://cloud.google.com/secret-manager/docs/)
  * Verify that the domain specified in the request is a valid, configured domain
  * Perform any final (specific to each cloud function) request data validation
  * Verify that confirmation was sent in the request (not all cloud functions require a confirmation)
  * Authenticate with and connect to any necessary GCP services ([Admin SDK](https://developers.google.com/admin-sdk), [Datastore](https://cloud.google.com/datastore), [Google Sheets](https://developers.google.com/sheets/api) etc) using domain-specific service accounts that have been granted [G Suite domain-wide delegation](https://developers.google.com/admin-sdk/directory/v1/guides/delegation)
  * Execute the action specific to the cloud function (approve a device, search, wipe a device, etc)
3. **Cleanup**
  * Update any [Datastore](https://cloud.google.com/datastore/) entities, as necessary
  * Log appropriate actions/events

## Configuration ##
All `gsuitemdm` cloud functions share a single master JSON configuration stored in Google [Secret Manager](https://cloud.google.com/secret-manager/docs/).  In order to know where to download this shared master configuration from (and because we do not hardcode such things), each cloud function has a tiny `.yaml` file that is deployed along with the source code to GCP. This `.yaml` file specifies the 2x environment variables pointing to the [Google Secret IDs](https://cloud.google.com/secret-manager/docs/managing-secrets) of the shared master configuration, and is used during cloud function app startup to download the shared master configuration. For example, the [ApproveDevice](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/approvedevice) cloud function's `.yaml` (example) configuration file looks like:
```
APPNAME: approvedevice
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```
From this `.yaml`, the ApproveDevice cloud function learns it's app name (`APPNAME` aka `approvedevice`), and the [Google Secret IDs](https://cloud.google.com/secret-manager/docs/managing-secrets) of it's API key (`SM_APIKEY_ID` aka Secret Manager secret `projects/12334567890/secrets/gsuitemdm_apikey`), and the shared master configuration (`SM_CONFIG_ID` aka Secret Manager secret `projects/12334567890/secrets/gsuitemdm_conf`). During app startup, each cloud function retrieves the appropriate secrets from Secret Manager. No other configuration files are necessary. 

See the `HOW-To Configure` section of each cloud function's `README.md` for full details.

### Configuration Secrets ###
Aside from each cloud function's `.yaml`, all configuration data, API keys and credentials are stored as secrets within Secret Manager. These must be created. 

The `gsuitemdm` system requires the following Secret Manager secrets:

**Secret Name** | **Purpose**
:--- | :---
`gsuitemdm_apikey` | Key used to authenticate API requests
`gsuitemdm_conf` | Shared cloud function master configuration
`gsuitemdm_slacktoken` | Token used to authenticate `slackdirectory` API requests from Slack
`credentials_DOMAINNAME` | Service account credentials JSON for each G Suite DOMAINNAME

#### Creating Configuration Secrets ####

##### Creating the shared cloud function master configuration secret #####
Use the included `cloudfunctions_conf_example.json` to create your own master configuration.
```
$ gcloud beta secrets create gsuitemdm_conf \
  --replication-policy automatic \
  --data-file cloudfunctions_conf_master.json
```

##### Creating the API key secret #####
All calls to any `gsuitemdm` cloud function must be authenticated by sending along the correct API key. Create the API key by use of `echo` and piping into `gcloud` and specifying STDIN (`-`) as the data file:
```
$ echo -n "yourkeygoeshere" | gcloud beta secrets create gsuitemdm_conf \
  --replication-policy automatic \
  --data-file=-
```

##### Creating the per-G Suite domain credentials secrets #####
Assuming we want to configure `foo.com`, `bar.com` and `xyzzy.com` in the `gsuitemdm` system, and we have downloaded the relevant G Suite domain-specific service account JSON credentials files for `foo.com`, `bar.com` and `xyzzy.com` and named them appropriately:
```
$ for DOMAIN in foo bar xyzzy
  do
     gcloud beta secrets create credentials_${DOMAIN} \
     --replication-policy automatic \
     --data-file credentials_${DOMAIN}.com.json
  done
```

#### Creating the Slack token secret ####
When Slack calls the `slackdirectory` cloud function API, it will send along a token. This token is checked to verify that it was indeed Slack who made the API call. Create the secret using:
```
$ echo -n "yourslacktokengoeshere" | gcloud beta secrets create gsuitemdm_slacktoken \
  --replication-policy automatic \
  --data-file=-
```

#### Updating Configuration Secrets ####

##### Updating configuration secrets in Secret Manager #####
```
$ gcloud beta secrets versions add gsuitemdm_conf \
  --data-file cloudfunctions_conf_new.json
```


## Deployment ##
All `gsuitemdm` cloud functions are deployed to GCP in the same manner:

```
$ gcloud functions deploy <FUNCTION_NAME> \
  --runtime go111 \
  --trigger-http \
  --env-vars-file env_<FUNCTION>.yaml
```

See the `HOW-To Deploy` section of each cloud function's `README.md` for full details. 

You may find the [`deploy_all_cloudfunctions.sh`](https://github.com/rickt/gsuitemdm/blob/master/cloudfunctions/deploy_all_cloudfunctions.sh) shell script useful.
