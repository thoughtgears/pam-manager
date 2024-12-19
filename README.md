# Pam manager

## Description

A cloud run service that manages pam access request automatically based on requester and other parameters.
It will give the requester access to the requested resource if the requester is allowed to access the resource.

It uses Google Auth to authenticate the requester and you can use the JWT token to request PAM access.
