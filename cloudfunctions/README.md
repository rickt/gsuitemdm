# gsuitemdm Cloud Functions #

The various GSuiteMDM cloud functions are the core components of the system. They exist to perform various mobile device-related tasks (`approve` a device, `search` for a device, `block` a device, etc), and are used extensively by `mdmtool` and of course can be called via `curl`.

## Design ##
All of the cloud functions are designed to be as simple as possible, and all follow the same general principles, and GCP cloud function design principles/best practices. A cloud function (as used in GSuiteMDM) is an http(s)-triggered mini-webserver. They are deployed to GCP, and scale up/down as needed. A high-level overview of the basic GSuiteMDM cloud functions model is:

1. Basic Checks
  * https listener starts up, listens for requests
  * Verify incoming requests don't have a null body and appear to be valid JSON for our API
  * Retrieve the GSuiteMDM API key from [Secret Manager](https://cloud.google.com/secret-manager/docs/)
  * Verify that a valid API key was sent in the request
  * Verify that a correct action (specific to each cloud function) was sent in the request
  * Perform basic sanity checks on the action-specific data (specific to each cloud function) that was sent in the request
2. GSuiteMDM service starts
  * Retrieve the shared GSuiteMDM configuration from [Secret Manager](https://cloud.google.com/secret-manager/docs/)
  * Retrieve all G Suite domain configurations from [Secret Manager](https://cloud.google.com/secret-manager/docs/)
  * Verify that the domain specified in the request is a valid, configured domain
  * Perform any final (specific to each cloud function) request data validation
  * Verify that confirmation was sent in the request (not all cloud functions require a confirmation)
  * Authenticate with and connect to any necessary GCP services ([Admin SDK](https://developers.google.com/admin-sdk), [Datastore](https://cloud.google.com/datastore), [Google Sheets](https://developers.google.com/sheets/api) etc) using domain-specific service accounts that have been granted G Suite domain-wide delegation
  * Execute the action specific to the cloud function (approve a device, search, wipe a device, etc)
3. Cleanup
  * Update any documents or Datastore entities, as necessary
  * Log appropriate actions/events

## Configuration ##
All of the GSuiteMDM cloud functions have a tiny `.yaml` file that points to actual app configuration stored in Google [Secret Manager](https://cloud.google.com/secret-manager/docs/). For example, the [ApproveDevice](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/approvedevice) cloud function's `.yaml` (example) configuration file looks like:
```
APPNAME: approvedevice
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```
From this `.yaml`, the ApproveDevice cloud function learns it's app name, and the [Google Secret IDs](https://cloud.google.com/secret-manager/docs/managing-secrets) of it's API key, and it's configuration. During app startup, each cloud function retrieves the appropriate secrets from Secret Manager. No other configuration files are necessary. 

All GSuiteMDM cloud functions share the same configuration (i.e., there is a single Secret Manager secret for all cloud functions). To use the example configuration, all cloud functions share the `projects/12334567890/secrets/gsuitemdm_conf` secret ID for their configuration and the `projects/12334567890/secrets/gsuitemdm_apikey` secret ID for their API key.


