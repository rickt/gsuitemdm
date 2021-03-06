# `gsuitemdm` Setup
For these example setup instructions, we will make the following critical assumptions:
* 3x G Suite domains (`foo.com`, `bar.com`, `xyzzy.com`) are G Suite domains under your control and all have mobile devices protected by [G Suite MDM](https://support.google.com/a/answer/1734200?hl=en)
* We have chosen `foo.com` to be the so-called "master domain", mainly because that is where the [ops team mobile device tracking spreadsheet](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/updatesheet) lives

## Overview of Setup ##
0. Clone the repo
1. Setup GCP projects 
2. Enable necessary APIs in those projects
3. Create & download [service account](https://cloud.google.com/iam/docs/service-accounts) [JSON credential files](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) for all G Suite domains
4. Grant API scope permissions to service accounts 
5. Create [Secret Manager](https://cloud.google.com/secret-manager/docs/) configuration secrets
6. Setup Google Sheet template for ops team mobile device tracking spreadsheet
7. Configure the `.yamls`
8. Deploy *all the things*

## Setup Details ##

### 0. Clone the `gsuitemdm` repo ###
```
$ git clone https://github.com/rickt/gsuitemdm`
```
### 1. Setup GCP projects ###
GCP best practices dictate that `gsuitemdm` requires a project in each G Suite domain that will be configured. 
#### 1.1 Setup a GCP project in the 'master' domain for `gsuitemdm` ####
Use `gcloud` to authenticate as an user in the `foo.com` G Suite 'master' domain, create a new project, setup billing:
```
$ gcloud auth login user@foo.com
$ gcloud projects create mdm-foo
$ gcloud beta billing accounts list
ACCOUNT_ID            NAME                        OPEN  MASTER_ACCOUNT_ID
000000-111111-222222  Foo Main Billing Account    True
$ gcloud beta billing projects link mdm-foo \
  --billing-account 000000-111111-222222
```
#### 1.2 Setup GCP projects in other domains ####
```
$ gcloud auth login user@bar.com
$ gcloud projects create mdm-bar
$ gcloud beta billing projects link mdm-bar \
  --billing-account blah-blah-blah

$ gcloud auth login user@xyzzy.com
$ gcloud projects create mdm-xyzzy
$ gcloud beta billing projects link mdm-xyzzy \
  --billing-account gude-gude-tama
```
### 2. Enable necessary APIs in the new projects ###
Now we need to enable some APIs in the new projects. Note that billing *must* be properly setup in the projects before attempting to enable APIs, as some APIs will fail to enable if a legit billing account has not been linked to your GCP project. 
#### 2.1 Enable APIs in the master project ####
Since the core components of the `gsuitemdm` system will run within the `mdm-foo` GCP project in the 'master domain' `foo.com`, this project needs more APIs enabled than other domains/projects. So, we must enable the following APIs: [Admin SDK/Directory API](https://developers.google.com/admin-sdk), [Cloud Functions](https://cloud.google.com/functions/docs/reference/rest), [Cloud Scheduler](https://cloud.google.com/scheduler/docs/reference/rest/), [Datastore](https://cloud.google.com/datastore/docs/reference/data/rest/), [Stackdriver Logging](https://cloud.google.com/logging/docs/reference/v2/rest), [Secret Manager](https://cloud.google.com/secret-manager/docs/accessing-the-api), [Sheets](https://developers.google.com/sheets/api):
```
$ gcloud auth login user@foo.com
$ gcloud config set project mdm-foo
$ for API in admin cloudfunctions cloudscheduler datastore logging secretmanager sheets
do
   gcloud services enable ${API}.googleapis.com
done
```
The remaining domain's GCP projects only require the [Admin SDK/Directory API](https://developers.google.com/admin-sdk) to be enabled: 
```
$ gcloud auth login user@bar.com
$ gcloud config set project mdm-bar
$ gcloud services enable admin.googleapis.com

$ gcloud auth login user@xyzzy.com
$ gcloud config set project mdm-xyzzy
$ gcloud services enable admin.googleapis.com
```
### 3. Create & download [service account](https://cloud.google.com/iam/docs/service-accounts) [JSON credential files](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) for all G Suite domains ###
Now we must create service accounts within each GCP project in each of our G Suite domains. 

#### 3.1 Create the service accounts in each of the configured domains
Unfortunately, there is no `gcloud`  command or API available to automate the following steps, they must be performed via web. Some pseudo-code might help:

`foreach DOMAIN in foo bar xyzzy`

`do`

* Login to GCP console as `user@$DOMAIN.com`
* Select `mdm-$DOMAIN` project
* Configure OAuth Consent Screen at [GCP Console](https://console.cloud.google.com/apis/credentials) `--> APIs & Services --> OAuth Consent Screen`
  * Type: `External`, App name: `mdm-$DOMAIN`, everything else is default, click `Save`
  * More details [available here](https://support.google.com/cloud/answer/6158849?hl=en) `--> User Consent`
* Create service account at [GCP Console](https://console.cloud.google.com/iam-admin/serviceaccounts) `--> IAM & Admin --> Service Accounts -- > Create Service Account`
  * Account Name: `G Suite MDM Service Account`, Account ID: `gsuitemdm`
  * Skip roles in screen 2
  * Create & download JSON key, naming convention: `credentials_$DOMAIN.com.json`
  * More details [available here](https://developers.google.com/identity/protocols/OAuth2ServiceAccount#creatinganaccount)
* Enable Domain-Wide Delegation at [GCP Console](https://console.cloud.google.com/apis/credentials) `--> APIs & Services --> Credentials`
  * Edit service account `mdm-$DOMAIN` and ENABLE `G Suite Domain-wide Delegation`

`done`

Note that it is *absolutely essential* that Domain-Wide Delegation is enabled for all service accounts!!! If you find that the Domain-Wide Delegation check box is not selectable, just [configure & save the OAuth Consent Screen](https://support.google.com/cloud/answer/6158849?hl=en). The Domain-Wide Delegation check box will become selectable after this.

At this point in our example setup, we have the following domains, projects & service accounts:

G Suite Domain | GCP Project | Service Account | Credentials JSON
:--- | :--- | :--- | :---
`foo.com` | `mdm-foo` | `gsuitemdm@mdm-foo.iam.gserviceaccount.com` | `credentials_foo.com.json`
`bar.com` | `mdm-bar` | `gsuitemdm@mdm-bar.iam.gserviceaccount.com` | `credentials_bar.com.json`
`xyzzy.com` | `mdm-xyzzy` | `gsuitemdm@mdm-xyzzy.iam.gserviceaccount.com` | `credentials_xyzzy.com.json`

### 4. Grant API scope permissions to service accounts ###
Now that we have created the service accounts, they need to be access to some Google API scopes. Following our example setup, these steps must be performed by a G Suite Super Administrator user in each of the `foo.com`, `bar.com` and `xyzzy.com` domains as per [these instructions](https://developers.google.com/identity/protocols/OAuth2ServiceAccount#delegatingauthority), starting from the `"Then, an administrator of the G Suite domain must complete [...]"` section. 

Starting with the 'master' domain, within the [Admin Console](https://admin.google.com) for `foo.com`, the Client ID of the `gsuitemdm@mdm-foo.iam.gserviceaccount.com` service account must be granted the following scopes:
```
https://www.googleapis.com/auth/admin.directory.device.mobile.action
https://www.googleapis.com/auth/admin.directory.device.mobile.readonly
https://www.googleapis.com/auth/spreadsheets
```
And within the [Admin Console](https://admin.google.com) for `bar.com`, the Client ID of the `gsuitemdm@mdm-bar.iam.gserviceaccount.com` service account must be granted the following scopes:
```
https://www.googleapis.com/auth/admin.directory.device.mobile.action
https://www.googleapis.com/auth/admin.directory.device.mobile.readonly
```
And finally, within the [Admin Console](https://admin.google.com) of `xyzzy.com`, the Client ID of the `gsuitemdm@mdm-xyzzy.iam.gserviceaccount.com` service account must be granted the following scopes:
```
https://www.googleapis.com/auth/admin.directory.device.mobile.action
https://www.googleapis.com/auth/admin.directory.device.mobile.readonly
```
At this point, our service accounts have been granted the necessary authority to use and query the Admin SDK APIs for their respective G Suite domains, and the service account in the master domain `foo.com` has additionally been granted access to the Sheets API (because the mobile device tracking sheet lives inside Google Drive within the master domain).

### 5. Create [Secret Manager](https://cloud.google.com/secret-manager/docs/) configuration secrets ###
#### 5.1 Create the per-G Suite service account domain credential secrets ####
Using the service account JSON credential files you [downloaded in step 3.1](https://github.com/rickt/gsuitemdm/blob/master/docs/SETUP.md#31-create-the-service-accounts-in-each-of-the-configured-domains), create the secrets in the master GCP project:
```
$ gcloud auth login user@foo.com
$ gcloud config set project mdm-foo
$ for DOMAIN in foo bar xyzzy
  do
     gcloud secrets create credentials_${DOMAIN} \
     --replication-policy automatic \
     --data-file credentials_${DOMAIN}.com.json
  done
```
#### 5.2 Create the shared master configuration secret ####
Also within the master project, use the included [`gsuitemdm_conf_example.json`](https://github.com/rickt/gsuitemdm/blob/master/cloudfunctions/gsuitemdm_conf_example.json) as a template to create your own master configuration, then create the secret: 
```
$ gcloud config set project mdm-foo
$ gcloud secrets create gsuitemdm_conf \
  --replication-policy automatic \
  --data-file cloudfunctions_conf_new.json
```
#### 5.3 Create the API key secret ####
All calls to any `gsuitemdm` cloud function must be authenticated by sending along the correct API key. Create the API key by use of `echo` and piping into `gcloud` and specifying STDIN (`-`) as the data file:
```
$ gcloud config set project mdm-foo
$ echo -n "yourapikeygoeshere" | gcloud secrets create gsuitemdm_apikey \
  --replication-policy automatic \
  --data-file=-
```
#### 5.4 Create the Slack token secret ####
When Slack calls the `slackdirectory` cloud function API, it will send along a token. This token is checked to verify that it was indeed Slack who made the API call. Create the secret for that token using:
```
$ gcloud config set project mdm-foo
$ echo -n "yourslacktokengoeshere" | gcloud secrets create gsuitemdm_slacktoken \
  --replication-policy automatic \
  --data-file=-
```
You can configure the token that Slack sends to `slackdirectory` when creating/editing your own `/phone` slash command at [`Yourslack Admin --> Manage Apps --> Custom`](https://YOURSLACK.slack.com/apps/manage/custom-integrations) `--> Slash Commands`

At this point, we have the following secrets:
```
$ gcloud config set project mdm-foo
$ gcloud secrets list 
NAME                           CREATED              REPLICATION_POLICY  LOCATIONS
credentials_bar                2020-01-24T16:09:20  automatic           -
credentials_foo                2020-01-24T16:09:22  automatic           -
credentials_xyzzy              2020-01-24T16:09:23  automatic           -
gsuitemdm_apikey               2020-01-24T22:59:47  automatic           -
gsuitemdm_conf                 2020-01-27T15:08:50  automatic           -
gsuitemdm_slacktoken           2020-01-27T22:25:29  automatic           -
```
### 6. Setup Google Sheet template for ops team mobile device tracking spreadsheet ###
1. Make a copy of [this Google Sheet](https://update.url) and save it in Google Drive. Now get the ID of your sheet; this is the part after `https://docs.google.com/spreadsheets/d/` in the sheet's URL but before `/edit`. Add that sheet ID to the main JSON configuration file, `"sheetid": "yourgooglesheetidgoeshere"`
2. Add the email address of the G Suite user who you wish to update the Google sheet as, to the main JSON configuration file, `"sheetwho": "username@yourgsuitedomain.com"`

### 7. Configure the `.yamls` ###
Copy and modify the `FUNCTION_env_example.yaml` files in each cloud function folder to create your own `FUNCTION_env.yaml` files that you will use to deploy the cloud functions. 

Be sure to modify the `SM_APIKEY_ID` and `SM_CONFIG_ID` entries to point to the IDs of the secrets you created in [step 5](https://github.com/rickt/gsuitemdm/blob/master/docs/SETUP.md#5-create-secret-manager-configuration-secrets).

### 8. Deploy & test *all the things* ###
Now that everything is in place, update the [`deploy_all_cloudfunctions.sh`](https://github.com/rickt/gsuitemdm/blob/master/cloudfunctions/deploy_all_cloudfunctions.sh) script to suit your environment/project/folders/etc, and then run it! 

Once all the cloud functions are deployed, you can test using `mdmtool` or use `curl` directly.

### 9. Adding new Domains ###
To add new G Suite domains to `gsuitemdm`:

#### 9.1 Add domain-specific credentials as Secret Manager secrets ####
Assuming your new domain is `foobarbaz.com`, add the JSON credentials file for the service account in the MDM project for your new domain into Secret Manager:
```
$ gcloud secrets create credentials_foobarbaz \
  --replication-policy automatic \
  --data-file credentials_foobarbaz.com.json
```

#### 9.2 Modify core `gsuitemdm_conf.json` and update in Secret Manager ####
Add the new domain into the domains array in your local copy of the core `gsuitemdm` configuration file and add it as a new version into Secret Manager:
```
$ cloud secrets versions add gsuitemdm_conf --data-file gsuitemdm_conf.json
```

#### 9.3 Add 'accessors' role for secret reading ####
Add the `Secret Manager Secret Accessor` role to the master project service account (in the example case, `gsuitemdm@mdm-foo.iam.gserviceaccount.com`).
