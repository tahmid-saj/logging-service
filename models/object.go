package models

import (
	"logging-service/object"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Object struct {
	BucketName string `json:"bucketName"`
	ObjectKey  string `json:"objectKey"`
}

type ObjectVersion struct {
	ObjectKey  string `json:"objectKey"`
	VersionID string `json:"versionID"`
}

type ObjectDownloadInput struct {
	FileName string `json:"fileName"`
}

type ObjectCreateInput struct {
	ObjectKey string `json:"objectKey"`
	FileName string `json:"fileName"`
}

type ObjectVersionInput []ObjectVersion

// object
func GetListObjects(bucketName string) (*Response, error) {
	s3Client := InitS3()

	objects, err := object.ListObjects(s3Client, bucketName)
	if err != nil {
		return &Response{
			Ok: false,
			Response: nil,
		}, err
	}

	var resObjects []*Object
	for _, object := range objects {
		resObjects = append(resObjects, &Object{
			BucketName: bucketName,
			ObjectKey: *object.Key,
		})
	}

	return &Response{
		Ok: true,
		Response: resObjects,
	}, nil
}

func PostDownloadObject(bucketName string, objectKey string, fileName string) (*Response, error) {
	s3Client := InitS3()

	err := object.DownloadObject(s3Client, bucketName, objectKey, fileName)
	if err != nil {
		return &Response{
			Ok: false,
			Response: false,
		}, err
	}

	return &Response{
		Ok: true,
		Response: true,
	}, nil
}

func GetListObjectVersions(bucketName string) (*Response, error) {
	s3Client := InitS3()

	objectVersions, err := object.ListObjectVersions(s3Client, bucketName)
	if err != nil {
		return &Response{
			Ok: false,
			Response: false,
		}, err
	}

	resObjectVersions := make(map[string][]string)
	for _, objectVersion := range objectVersions {
		resObjectVersions[*objectVersion.Key] = append(resObjectVersions[*objectVersion.Key], *objectVersion.VersionId)
	}

	return &Response{
		Ok: true,
		Response: resObjectVersions,
	}, nil
}

func PostUploadObject(bucketName string, objectKey string, fileName string) (*Response, error) {
	s3Client := InitS3()

	err := object.UploadObject(s3Client, bucketName, objectKey, fileName)
	if err != nil {
		return &Response{
			Ok: false,
			Response: false,
		}, err
	}

	return &Response{
		Ok: true,
		Response: true,
	}, nil
}

func DeleteObjects(bucketName string, objects []ObjectVersion) (*Response, error) {
	s3Client := InitS3()

	var objectIdentifiers []types.ObjectIdentifier
	for _, object := range objects {
		objectIdentifiers = append(objectIdentifiers, types.ObjectIdentifier{
			Key: &object.ObjectKey,
			VersionId: &object.VersionID,
		})
	}
	
	err := object.DeleteObjects(s3Client, bucketName, objectIdentifiers, true)
	if err != nil {
		return &Response{
			Ok: false,
			Response: false,
		}, err
	}

	return &Response{
		Ok: true,
		Response: true,
	}, nil
}

func DeleteObject(bucketName string, objectKey string, versionID string) (*Response, error) {
	s3Client := InitS3()
	
	deleted, err := object.DeleteObject(s3Client, bucketName, objectKey, versionID, true)
	if err != nil || !deleted {
		return &Response{
			Ok: false,
			Response: false,
		}, err
	}

	return &Response{
		Ok: false,
		Response: true,
	}, nil
}