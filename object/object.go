package object

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"logging-service/utils"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func ListObjects(s3Client *s3.Client, bucketName string) ([]types.Object, error) {
	result, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	var contents []types.Object
	if err != nil {
		log.Printf("Couldn't list objects in bucket %v. %v.\n", bucketName, err)
	} else {
		contents = result.Contents
	}

	return contents, err
}

func PrintObjects(objects []types.Object) {
	fmt.Println("\tKey\t\t\t\t\tLast modified")
	for _, object := range objects {
		fmt.Printf("\t%v\t\t\t\t%v\n", *object.Key, *object.LastModified)
	}
}

func DownloadObject(s3Client *s3.Client, bucketName string, objectKey string, fileName string) error {
	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v: %v. %v.\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()

	file, err := os.Create(filepath.Join(utils.DOWNLOADOBJECTFILEPATH, filepath.Base(fileName)))
	if err != nil {
		log.Printf("Couldn't create file %v. %v.\n", fileName, err)
		return err
	}
	defer file.Close()
	
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. %v.\n", objectKey, err)
	}

	_, err = file.Write(body)

	return err
}

func ReadObject(s3Client *s3.Client, bucketName string, objectKey string) ([][]string, error) {
	readRequestInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(objectKey),
	}

	result, err := s3Client.GetObject(context.TODO(), readRequestInput)
	if err != nil {
		log.Print(err)
	}

	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object %v. %v.\n", objectKey, err)
		return nil, err
	}

	reader := csv.NewReader(bytes.NewBuffer(body))
	rows, err := reader.ReadAll()
	if err != nil {
		log.Printf("Couldn't read object %v. %v.\n", objectKey, err)
		return nil, err
	}

	return rows, nil
}

func ListObjectVersions(s3Client *s3.Client, bucketName string) ([]types.ObjectVersion, error) {
	var err error
	var output *s3.ListObjectVersionsOutput
	var versions []types.ObjectVersion

	input := &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucketName),
	}
	versionPaginator := s3.NewListObjectVersionsPaginator(s3Client, input)

	for versionPaginator.HasMorePages() {
		output, err = versionPaginator.NextPage(context.TODO())

		if err != nil {
			var noSuchBucket *types.NoSuchBucket

			if errors.As(err, &noSuchBucket) {
				log.Printf("Bucket %s does not exist.\n", bucketName)
			}

			break
		} else {
			versions = append(versions, output.Versions...)
		}
	}

	return versions, err
}

func PrintObjectVersions(objectVersions []types.ObjectVersion) {
	fmt.Println("\tKey\t\t\t\t\tLast modified")
	for _, objectVersion := range objectVersions {
		fmt.Printf("\t%v\t\t\t\t%v\n", *objectVersion.Key, *objectVersion.LastModified)
	}
}

func DeleteObjects(s3Client *s3.Client, bucketName string, 
	objects []types.ObjectIdentifier, bypassGovernance bool) error {

	if len(objects) == 0 {
		return nil
	}

	input := s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{
			Objects: objects,
			Quiet: aws.Bool(true),
		},
	}

	if bypassGovernance {
		input.BypassGovernanceRetention = aws.Bool(true)
	}

	delOut, err := s3Client.DeleteObjects(context.TODO(), &input)
	if err != nil || len(delOut.Errors) > 0 {
		log.Printf("Error deleting objects from bucket %s.\n", bucketName)

		if err != nil {
			var noSuchBucket *types.NoSuchBucket

			if errors.As(err, &noSuchBucket) {
				log.Printf("Bucket %s does not exist.\n", bucketName)
				err = noSuchBucket
			}
		} else if len(delOut.Errors) > 0 {
			for _, outErr := range delOut.Errors {
				log.Printf("%s: %s\n", *outErr.Key, *outErr.Message)
			}

			err = fmt.Errorf("%s", *delOut.Errors[0].Message)
		}
	}
	
	return err
}

func DeleteObject(s3Client *s3.Client, bucketName string, objectKey string, versionID string,
	bypassGovernance bool) (bool, error) {

	deleted := false
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(objectKey),
	}

	if versionID != "" {
		input.VersionId = aws.String(versionID)
	}

	if bypassGovernance {
		input.BypassGovernanceRetention = aws.Bool(true)
	}
	
	_, err := s3Client.DeleteObject(context.TODO(), input)
	if err != nil {
		var noSuchObjectKey *types.NoSuchKey
		var apiError error

		if errors.As(err, &noSuchObjectKey) {
			log.Printf("Object %s does not exist in %s.\n", objectKey, bucketName)
			err = noSuchObjectKey
		} else if errors.As(err, &apiError) {
			fmt.Println(err)
			err = nil
		}
	} else {
		deleted = true
	}

	return deleted, err
}

func UploadObject(s3Client *s3.Client, bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(filepath.Join(utils.UPLOADOBJECTFILEPATH, filepath.Base(fileName)))
	if err != nil {
		log.Printf("Couldn't open file %v to upload. %v.\n", fileName, err)
	} else {
		defer file.Close()

		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key: aws.String(objectKey),
			Body: file,
		})

		if err != nil {
			log.Printf("Couldn't upload file %v to %v: %v. %v.\n", fileName, bucketName, objectKey, err)
		}
	}

	return err
}