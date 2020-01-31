# `gsuitemdm` Setup
For these example setup instructions, we will make the following critical assumptions:
* 3x G Suite domains (`foo.com`, `bar.com`, `xyzzy.com`) are G Suite domains under your control and all have mobile devices protected by [G Suite MDM](https://support.google.com/a/answer/1734200?hl=en)
* We have chosen `foo.com` to be the so-called "master domain", mainly because that is where the [ops team mobile device tracking spreadsheet](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/updatesheet) lives

## Overview of Setup ##
1. Setup a GCP project in your organization for `gsuitemdm`
2. Configure a billing account in that project
3. Enable necessary APIs in that project
4. Create & download [service account](https://cloud.google.com/iam/docs/service-accounts) [JSON credential files](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) for all G Suite domains
5. Grant [domain-wide delegation](https://developers.google.com/admin-sdk/directory/v1/guides/delegation) permissions to service accounts
6. Grant [Directory Admin SDK API scope permissions](https://developers.google.com/admin-sdk/directory/v1/guides/authorizing) to service accounts
7. Create [Secret Manager](https://cloud.google.com/secret-manager/docs/) configuration secrets
8. Setup Google Sheet template for ops team mobile device tracking spreadsheet

## Setup Details ##

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
### 3. Enable necessary APIs in the new projects ###
Now we need to enable some APIs in the new projects. Note that billing *must* be properly setup in the projects before attempting to enable APIs, as some APIs will fail to enable if a legit billing account has not been linked to your GCP project. 
#### 3.1 Enable APIs in the master project ####
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
### 4. Create & download [service account](https://cloud.google.com/iam/docs/service-accounts) [JSON credential files](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) for all G Suite domains ###
Now we must create service accounts within each GCP project in each of our G Suite domains. 

#### 4.1 Create the service accounts in each of the configured domains
Unfortunately, there is no `gcloud`  command or API available to automate these steps. Some pseudo-code might help:
```
foreach DOMAIN in foo bar xyzzy
do
  login to GCP console as user@$DOMAIN.com
  choose `mdm-$DOMAIN` project
  create service account as per [these Google developer docs detail](https://developers.google.com/identity/protocols/OAuth2ServiceAccount#creatinganaccount)
done
```
`foreach DOMAIN in foo bar xyzzy`

`do`

* login to GCP console as user@$DOMAIN.com
* choose `mdm-$DOMAIN` project
* create service account as per [these Google developer docs detail](https://developers.google.com/identity/protocols/OAuth2ServiceAccount#creatinganaccount)

`done`

If greyed out...

#### 4.2 Create the service accounts in additional G Suite domains ####
Now you need to create a service account in each of the additional domains we want to configure (`bar.com`, `xyzzy.com`). 

### 5. Grant [domain-wide delegation](https://developers.google.com/admin-sdk/directory/v1/guides/delegation) permissions to service accounts ###
Docs coming.

### 6. Grant [Directory Admin SDK API scope permissions](https://developers.google.com/admin-sdk/directory/v1/guides/authorizing) to service accounts ###
Docs coming.

### 7. Create [Secret Manager](https://cloud.google.com/secret-manager/docs/) configuration secrets ###
#### 7.1 Create the per-G Suite service account domain credential secrets ####
Using the service account JSON credential files you downloaded earlier, create the secrets: 
```
$ for DOMAIN in foo bar xyzzy
  do
     gcloud beta secrets create credentials_${DOMAIN} \
     --replication-policy automatic \
     --data-file credentials_${DOMAIN}.com.json
  done
```
#### 7.2 Create the shared master configuration secret ####
Use the included [`gsuitemdm_conf_example.json`](https://github.com/rickt/gsuitemdm/blob/master/cloudfunctions/gsuitemdm_conf_example.json) as a template to create your own master configuration, then create the secret: 
```
$ gcloud beta secrets create gsuitemdm_conf \
  --replication-policy automatic \
  --data-file cloudfunctions_conf_new.json
```
#### 7.3 Create the API key secret ####
All calls to any `gsuitemdm` cloud function must be authenticated by sending along the correct API key. Create the API key by use of `echo` and piping into `gcloud` and specifying STDIN (`-`) as the data file:
```
$ echo -n "yourkeygoeshere" | gcloud beta secrets create gsuitemdm_apikey \
  --replication-policy automatic \
  --data-file=-
```
#### 7.4 Create the Slack token secret ####
When Slack calls the `slackdirectory` cloud function API, it will send along a token. This token is checked to verify that it was indeed Slack who made the API call. Create the secret for that token using:
```
$ echo -n "yourslacktokengoeshere" | gcloud beta secrets create gsuitemdm_slacktoken \
  --replication-policy automatic \
  --data-file=-
```
You can configure the token that Slack sends to `slackdirectory` when creating/editing your own `/phone` slash command at [`Yourslack Admin --> Manage Apps --> Custom`](https://YOURSLACK.slack.com/apps/manage/custom-integrations) `--> Slash Commands`
### 8. Setup Google Sheet template for ops team mobile device tracking spreadsheet ###
Docs coming.


