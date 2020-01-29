# gsuitemdm Cloud Functions #

The various `gsuitemdm` cloud functions are the core components of the `gsuitemdm` system. Each cloud function exists to perform a single mobile device-related task (`approve` a device, `search` for a device, `block` a device, etc). The cloud functions are deployed to GCP to make the `gsuitemdm` API endpoints available via https, and are used extensively by the included command line tool [`mdmtool`](https://github.com/rickt/gsuitemdm/tree/master/mdmtool) as well as being able to be called via `curl` or whatever http(s) library you prefer.

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

## Deployment ##
All `gsuitemdm` cloud functions are deployed to GCP in the same manner:

```
$ gcloud functions deploy <FUNCTION_NAME> --runtime go111 --trigger-http \
  --env-vars-file env_<FUNCTION>.yaml
```

See the `HOW-To Deploy` section of each cloud function's `README.md` for full details. 

You may find the [`deploy_all_cloudfunctions.sh`](https://github.com/rickt/gsuitemdm/blob/master/cloudfunctions/deploy_all_cloudfunctions.sh) shell script useful.
