# Pam manager

## Description

A cloud run service that manages pam access request automatically based on requester and other parameters.
It will give the requester access to the requested resource if the requester is allowed to access the resource.


## Setup

- Create a pub/sub topic for the alert, `pam-grant-request`
- Grant `service-558467044633@gcp-sa-monitoring-notification.iam.gserviceaccount.com` `pubsub.publisher` permission on the selected topic (project level)
- Create a notification channel for Pub/Sub and select the topic `pam-grant-request`
- Create a log alert that sends a message notification to the created notification channel
```shell
resource.type="audited_resource" 
protoPayload.methodName="google.cloud.privilegedaccessmanager.v1alpha.PrivilegedAccessManager.CreateGrant" 
protoPayload.@type="type.googleapis.com/google.cloud.audit.AuditLog"
```
