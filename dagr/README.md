This module manages Dagger Engines on a bunch of platforms.

Currently only <Fly.io> is implemented. Kubernetes would be a great second addition - contributions welcome!

```console
dagger functions

Name       Description
get-app    Returns the app name: `dagger call get-app`
on-flyio   Manages Dagger on Fly.io: `dager call on-flyio --token=env:FLY_API_TOKEN`
```

## on-flyio

To see what commands are available, run `dagger functions on-flyio`

The most useful command will most likely be: `dagger call on-flyio --token=env:FLY_API_TOKEN deploy`

This is the quickest (and currently deprecated) way of setting the `FLY_API_TOKEN` environment variable: `export FLY_API_TOKEN=$(flyctl auth token)`

## TL;DR

To see what else is possible, consider running the following command: [`just`](https://github.com/casey/just)

This is what you can expect to see:

https://github.com/user-attachments/assets/c7a857a9-b561-454c-bb54-1a05f4a0cb1d
