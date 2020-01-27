# gsuitemdm Cloud Functions #

The various GSuiteMDM cloud functions are the core components of the system. They exist to perform various mobile device-related tasks (`approve` a device, `search` for a device, `block` a device, etc), and are used extensively by `mdmtool` and of course can be called via `curl`.

## Configuration ##
All of the GSuiteMDM cloud functions have a tiny `.yaml` file that points to actual app configuration stored in Google [Secret Manager](https://cloud.google.com/secret-manager/docs/). For example, the [ApproveDevice](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/approvedevice) cloud function's `.yaml` (example) configuration file looks like:
```
APPNAME: approvedevice
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```
From this `.yaml`, the ApproveDevice cloud function learns it's app name, and the [Google Secret IDs](https://cloud.google.com/secret-manager/docs/managing-secrets) of it's API key, and it's configuration. During app startup, each cloud function retrieves the appropriate secrets from Secret Manager. No other configuration files are necessary. 



