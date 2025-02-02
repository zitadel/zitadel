---
description: In Actions version 1, you are unable to create a user deletion event.
tags: 
    - FAQ
---

# What are ZITADEL Minimum Self-Hosted Specs?

## Overview

Looking to define minimum specifications needed for Self-Hosting a ZITADEL instance?

## Solution

When self-hosting ZITADEL, it can consume around 512MB RAM and run with less than 1 CPU core. However the database does consume generally 2 CPU under normal operating conditions and 6GB RAM with some caching to it.

In the event of hashing passwords it has the potential to dramatically increase the CPU usage which would foster the recommendation for running 4 CPU cores.

Reference Documentation:

* [ZITADEL Production Resource Usage](https://zitadel.com/docs/self-hosting/manage/production#minimum-system-requirements)
