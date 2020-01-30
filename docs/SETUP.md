# `gsuitemdm` Setup
For these example setup instructions, we will make the following critical assumptions:
* 3x G Suite domains (`foo.com`, `bar.com`, `xyzzy.com`) are G Suite domains under your control and all have mobile devices protected by [G Suite MDM](https://support.google.com/a/answer/1734200?hl=en)

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
$ gcloud config set project gsuitemdm
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

