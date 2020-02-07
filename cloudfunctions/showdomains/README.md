# gsuitemdm Cloud Function `showdomains` #

A [cloud Function](https://cloud.google.com/functions/) component of the [`gsuitemdm`](https://github.com/rickt/gsuitemdm) package that outputs a list of all G Suite domains currently configured in the `gsuitemdm` system.

The `showdomains` API is used by the [`mdmtool`](#mdmtool) command line utility (`showdomains` command).

## HOW-TO Configure `showdomains` ##
`showdomains` uses a `.yaml` file containing several environment variables the cloud function reads during app startup. These environment variables point the app to the shared master cloud function configuration and API key that are stored as [Secret Manager secrets](https://cloud.google.com/secret-manager/docs/managing-secrets). An example `.yaml` file for `showdomains`:

```yaml
APPNAME: showdomains
SM_APIKEY_ID: projects/12334567890/secrets/gsuitemdm_apikey
SM_CONFIG_ID: projects/12334567890/secrets/gsuitemdm_conf
```

## HOW-TO Deploy `showdomains` ##
```
$ gcloud functions deploy ShowDomains \
  --runtime go111 \
  --trigger-http \
  --env-vars-file env_showdomains.yaml
```

## HOW-TO Use `showdomains` ##

### API ###
Example expected JSON to show all configured domains:
```json
{
	"action": "showdomains",
	"key": "0123456789"
}
```

Example command line using `curl` and the above JSON to show all configured domains:

```
FIX
```

### `mdmtool` ##
Example command line using `mdmtool` to show all configured domains:
```
FIX
```
