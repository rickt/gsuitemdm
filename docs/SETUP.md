# `gsuitemdm` Setup
For these example setup instructions, we will make the following critical assumptions:
* 3x G Suite domains (`foo.com`, `bar.com`, `xyzzy.com`) are G Suite domains under your control and all have mobile devices protected by [G Suite MDM](https://support.google.com/a/answer/1734200?hl=en)

## Overview of Setup ##
1. Setup a GCP project in your organization for `gsuitemdm`
2. Enable necessary APIs in that project
3. Create & download [service account](https://cloud.google.com/iam/docs/service-accounts) [JSON credential files](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) for all G Suite domains
4. Grant [domain-wide delegation](https://developers.google.com/admin-sdk/directory/v1/guides/delegation) permissions to service accounts
5. Grant [Directory Admin SDK API scope permissions](https://developers.google.com/admin-sdk/directory/v1/guides/authorizing) to service accounts
6. Create [Secret Manager](https://cloud.google.com/secret-manager/docs/) configuration secrets
7. Setup Google Sheet template for ops team mobile device tracking spreadsheet

## Setup Details ##

### 1. Setup a GCP project in your organization for `gsuitemdm` ###
```
$ cloud projects create PROJECTNAME
```

### 2. Enable necessary APIs in that project ###
