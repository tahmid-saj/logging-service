# logging-service

Logging service API for log management in S3 from different frontends / APIs / databases. Can also view summaries of logs in dashboard. Developed using Go / Gin, AWS S3.

<br/>
<br/>

## Directory structure

The directory structure is as follows:

- **bucket/**  
  - Handles interactions with AWS S3 buckets for log storage.

- **conf/**  
  - Configuration files for the service.

- **data/**  
  - Contains example or sample log data.

- **models/**  
  - Data models for logging records and configurations.

- **object/**  
  - Manages S3 object operations such as upload and retrieval.

- **routes/**  
  - Defines API routes using the Gin framework.

- **utils/**  
  - Utility functions for log processing and S3 interactions.

- **main.go**  
  - Entry point for the logging service.

<br/>
<br/>

## Overview

### Design

The service is used mainly with a monitoring and notification service. The monitoring service can be found <a href="https://www.sitemonitoring.io/">here</a>. Similar services can be found <a href="https://whimsical.com/web-microservices-6uqvwWZtcBFsNJB2hepGy1">here</a> and below:

#### Similar services

<img width="834" alt="image" src="https://github.com/user-attachments/assets/b54088e7-870c-46dd-9cf6-2e5ec27d9d5c">
