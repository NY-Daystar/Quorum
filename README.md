<h1 style="text-align:center;" >
<img style="vertical-align: middle;" alt="Title" src="logo.ico" width="48" height="48" /> Quorum
</h1>

![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)
![Codacy Badge](https://app.codacy.com/project/badge/Grade/633862c8e36145e1af52fae32a14c31a)
[![Quorum-CI](https://github.com/NY-Daystar/Quorum/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/NY-Daystar/Quorum/actions/workflows/go.yml)
![License](https://img.shields.io/github/license/ny-daystar/Quorum)
[![Version](https://img.shields.io/github/tag/NY-Daystar/quorum.svg)](https://github.com/NY-Daystar/Quorum/releases)

[![GitHub Releases](https://img.shields.io/github/downloads/ny-daystar/quorum/total)](https://github.com/ny-daystar/quorum/releases)
[![Total views](https://img.shields.io/sourcegraph/rrc/github.com/NY-Daystar/quorum.svg)](https://sourcegraph.com/github.com/NY-Daystar/quorum)

![GitHub watchers](https://img.shields.io/github/watchers/ny-daystar/Quorum) ![GitHub forks](https://img.shields.io/github/forks/ny-daystar/Quorum) ![GitHub Repo stars](https://img.shields.io/github/stars/ny-daystar/Quorum)

![GitHub repo size](https://img.shields.io/github/repo-size/ny-daystar/Quorum) ![GitHub language count](https://img.shields.io/github/languages/count/ny-daystar/Quorum) ![GitHub top language](https://img.shields.io/github/languages/top/ny-daystar/Quorum)
![GitHub issues](https://img.shields.io/github/issues/ny-daystar/Quorum) ![GitHub closed issues](https://img.shields.io/github/issues-closed-raw/ny-daystar/Quorum)

# Summary

- [User Guide](#user-guide)
- [Requirements](#requirements)
- [Get started](#get-started)
- [Explainations](#explainations)
- [Contact](#contact)
- [Credits](#credits)

## User Guide

![Program](/docs/Main.png)

This application allows to backup your mail from gmail

1. Download `Quorum.exe` project from
   [this link](https://github.com/NY-Daystar/Quorum/releases/download/v0.0.1/Quorum.exe)

1. Then setup your [gmail credentials](#create-credentials)

1. Launch executable

1. It gonna ask you to authorize the app to access on your gmail account
   Like below
   ![request access](docs/Request_access.png)
   Authorize this app: <https://accounts.google.com/>
   🔄 Processus
    1. Click on the link
    1. Log into your gmail account
    1. Click on Authorize

1. At this stage you will get in your folder application this
   gmail-backup/
   │
   ├── Quorum.exe
   ├── config.json
   ├── credentials.json
   ├── token.json
   └── backup/

1. At now the application will backup all your mails into `backup/` folder
   ![Process done](docs/Done.png)

### Create credentials

1. Go into [Google Cloud Console](https://console.cloud.google.com/)

1. Create a project
   Click on project selector (top of page)
   New Project
   Name : Quorum
1. Activate `Gmail API`
   Menu → APIs & Services → Library
   Search : Gmail API
   Click on `Enable`
1. Configure OAuth consent screen
   Menu → APIs & Services → OAuth consent screen
   Choose Type : `External`
   Fill: App name → ex: Gmail Backup Tool
   Email → your email address

1. Create credentials
   Menu → Credentials
   Click → Create Credentials → OAuth client ID
   Choose Type : `Desktop App`

1. Finally download the json file le JSON
   It looks like

    ```json
    {
        "installed": {
            "client_id": "...",
            "project_id": "...",
            ...
        }
    }
    ```

    Rename the file into `credentials.json` and put it in quorum/credentials.json

1. Don't forget in `Audience section` add a `Test User`

## Requirements

- [.NET Framework](https://dotnet.microsoft.com/en-us/download/dotnet/7.0) >= 7.0
- For developpment: [VS 2022](https://visualstudio.microsoft.com/fr/vs/) >= 2022

## Get started

1. Download `Doppler` project from [this link](https://github.com/NY-Daystar/Doppler/releases/download/v1.7.0/Doppler-portable)

2. Extract zip on your computer

3. Launch `Doppler.exe`
    - The project ask you to choose a folder in your computer to rename files.
    - It will list folder files and submit several renaming.
    - After choosing one the application rename files automatically

# Explainations

```mermaid
flowchart TD
    A[Oauth Authentication]
    B[Listing mail]

    subgraph msg [For each mail]
        direction TD
        D["Get Contents (Attachments, images, Subject, Date)"]
        E[Decode in EML]
        F[Save EML file]
        G[Decode in HTML]
        H[Save HTML file]
        D-->E
        D-->G
        E-->F
        G-->H
    end

    A-->B
    B-->msg
```

## Contact

- To make a pull request: <https://github.com/NY-Daystar/doppler/pulls>
- To summon an issue: <https://github.com/NY-Daystar/doppler/issues>
- For any specific demand by mail: [luc4snoga@gmail.com](mailto:luc4snoga@gmail.com?subject=[GitHub]%doppler%20Project)

## Credits

Made by Lucas Noga.  
Licensed under GPLv3.
