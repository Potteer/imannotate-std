package handlers

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

// CheckScalityCredentials mirrors Imannotate S3, but adds the option to set the Endpoint
func CheckScalityCredentials(c *gin.Context) {
	s3opt := map[string]string{}
	c.Bind(&s3opt)

	conf := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(s3opt["id"], s3opt["secret"], ""),
		Endpoint:         aws.String(s3opt["server"]),
		Region:           aws.String(s3opt["region"]),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(conf)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, err.Error())
		return
	}
	svc := s3.New(sess)

	bl := s3.ListBucketsInput{}
	out, err := svc.ListBuckets(&bl)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, out.Buckets)
}
