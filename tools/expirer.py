#!/usr/bin/python
# -*- coding: utf-8 -*-

import boto3
from datetime import datetime, timedelta
import os
import pytz

S3_ENDPOINT_URL = os.environ.get("S3_ENDPOINT_URL")
S3_BUCKET = os.environ.get("S3_BUCKET", "test")
AWS_ACCESS_KEY_ID = os.environ.get("AWS_ACCESS_KEY_ID", "accessKey1")
AWS_SECRET_ACCESS_KEY = os.environ.get("AWS_SECRET_ACCESS_KEY", "verySecretKey1")
S3_OBJECT_RETENTION_MINUTES = int(os.environ.get("S3_OBJECT_RETENTION_MINUTES", "60"))
ISO_PREFIX = os.environ.get("ISO_PREFIX", "discovery-image")

max_timestamp = pytz.UTC.localize(datetime.now()) - timedelta(
    minutes=S3_OBJECT_RETENTION_MINUTES
)

client = boto3.client(
    "s3",
    use_ssl=False,
    endpoint_url=S3_ENDPOINT_URL,
    aws_access_key_id=AWS_ACCESS_KEY_ID,
    aws_secret_access_key=AWS_SECRET_ACCESS_KEY,
)
paginator = client.get_paginator("list_objects")
operation_parameters = {"Bucket": S3_BUCKET, "Prefix": ISO_PREFIX}
page_iterator = paginator.paginate(**operation_parameters)
for page in page_iterator:
    if not page.get("Contents"):
        continue
    for obj in page["Contents"]:
        if obj["LastModified"] < max_timestamp:
            print("Deleting expired: " + obj["Key"])
            client.delete_object(Bucket=S3_BUCKET, Key=obj["Key"])
