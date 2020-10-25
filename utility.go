package main

import "github.com/aws/aws-sdk-go/service/s3"

func ListBuckets() (*s3.ListBucketsOutput, error) {
	return s3session.ListBuckets(&s3.ListBucketsInput{})
}
