# `gsuitemdm` Setup
For these example setup instructions, we will make the following critical assumptions:
* 3x G Suite domains (`foo.com`, `bar.com`, `xyzzy.com`) are G Suite domains under your control and all have mobile devices protected by [G Suite MDM](https://support.google.com/a/answer/1734200?hl=en)
* We have chosen `foo.com` to be the "master" domain, mainly because that is where the [ops team mobile device tracking spreadsheet](https://github.com/rickt/gsuitemdm/tree/master/cloudfunctions/updatesheet) lives

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

### 1. Setup a GCP project in your organization for `gsuitemdm` ###
Create a new project:
```
$ gcloud projects create PROJECTNAME
```
Set the new project as your current/configured project:
```
$ gcloud config set project PROJECTNAME
```
### 2. Configure a billing account in that project
List your existing billing accounts:
```
$ gcloud beta billing accounts list
ACCOUNT_ID            NAME                        OPEN  MASTER_ACCOUNT_ID
000000-111111-222222  Main Billing Account        True
111111-222222-333333  Secondary Billing Account   True
```
Link the new GCP project to a billing account:
```
$ gcloud beta billing projects link PROJECTNAME --billing-account 000000-111111-222222
```
### 3. Enable necessary APIs in that project ###
```
$ for API in admin cloudfunctions cloudscheduler datastore logging secretmanager sheets
do
   gcloud services enable ${API}.googleapis.com
done
```
### 4. Create & download [service account](https://cloud.google.com/iam/docs/service-accounts) [JSON credential files](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) for all G Suite domains ###

### 5. Grant [domain-wide delegation](https://developers.google.com/admin-sdk/directory/v1/guides/delegation) permissions to service accounts ###

### 6. Grant [Directory Admin SDK API scope permissions](https://developers.google.com/admin-sdk/directory/v1/guides/authorizing) to service accounts ###

### 7. Create [Secret Manager](https://cloud.google.com/secret-manager/docs/) configuration secrets ###
#### 7.1 Create the shared master configuration secret ####
Use the included [`gsuitemdm_conf_example.json`](https://github.com/rickt/gsuitemdm/blob/master/cloudfunctions/gsuitemdm_conf_example.json) as a template to create your own master configuration, then create the secret: 
```
$ gcloud beta secrets create gsuitemdm_conf \
  --replication-policy automatic \
  --data-file cloudfunctions_conf_new.json
```
#### 7.2 Create the API key secret ####
All calls to any `gsuitemdm` cloud function must be authenticated by sending along the correct API key. Create the API key by use of `echo` and piping into `gcloud` and specifying STDIN (`-`) as the data file:
```
$ echo -n "yourkeygoeshere" | gcloud beta secrets create gsuitemdm_apikey \
  --replication-policy automatic \
  --data-file=-
```
#### 7.3 Create the per-G Suite service account domain credential secrets ####
Using the service account JSON credential files you downloaded earlier, create the secrets: 
```
$ for DOMAIN in foo bar xyzzy
  do
     gcloud beta secrets create credentials_${DOMAIN} \
     --replication-policy automatic \
     --data-file credentials_${DOMAIN}.com.json
  done
```

### 8. Setup Google Sheet template for ops team mobile device tracking spreadsheet ###

