package s3utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

//S3Connector is a wrapper for S3 and Session objects
type S3Connector struct {
	Svc  *s3.S3
	Sess *session.Session
}

//GetS3Service return a new service and returna S3Connector
func GetS3Service(server string, region string, id string, secret string) (S3Connector, error) {

	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(id, secret, ""),
		Endpoint:         aws.String(server),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	s, err := session.NewSession(config)
	if err != nil {
		return S3Connector{}, err
	}

	// Create S3 service client
	svc := s3.New(s)
	conn := S3Connector{
		Svc:  svc,
		Sess: s,
	}

	return conn, nil
}

//CreateBucket receives a service and tries to create a bucket
func CreateBucket(conn *S3Connector, name string) (bool, error) {
	svc := conn.Svc
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return false, err
	}

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(name),
	})
	return true, nil
}

//ListBuckets is a wrapper for was  ListBuckets
func ListBuckets(conn *S3Connector) ([]string, error) {
	svc := conn.Svc
	buckets, err := svc.ListBuckets(nil)
	result := make([]string, len(buckets.Buckets), len(buckets.Buckets))
	for i, b := range buckets.Buckets {
		result[i] = aws.StringValue(b.Name)
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

func listObjects(input *s3.ListObjectsV2Input,
	fn func(*s3.ListObjectsV2Input,
		func(*s3.ListObjectsV2Output, bool) bool) error) ([]*s3.Object, error) {
	const n = 500
	input.SetMaxKeys(int64(n))
	summaries := make([]*s3.Object, 0, n)
	err := fn(input,
		func(p *s3.ListObjectsV2Output, lastPage bool) bool {
			summaries = append(summaries, p.Contents...)
			return true
		})
	if err != nil {
		return nil, err
	}
	return summaries, nil
}

//ListObjectsInBucket wrapper for ListObjectsV2
func ListObjectsInBucket(conn *S3Connector, bucket *string, prefix *string, delimiter *string) ([]*s3.Object, error) {
	input := &s3.ListObjectsV2Input{}
	input.SetBucket(*bucket)
	if prefix != nil {
		input.SetPrefix(*prefix)
	}
	if delimiter != nil {
		input.SetDelimiter(*delimiter)
	}
	return listObjects(input, conn.Svc.ListObjectsV2Pages)
}

//PutObject is a wrapper for s3.PutObject
var uploader *s3manager.Uploader

//PutObject s3 uploader wrapper
func PutObject(conn *S3Connector, filename string, bucket string, key string, prefix string, delim string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if uploader == nil {
		uploader = s3manager.NewUploader(conn.Sess)
	}
	r, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	return r.Location, err

}
