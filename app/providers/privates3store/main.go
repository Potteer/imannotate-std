package privates3store

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/smileinnovation/imannotate/app/providers/privates3store/s3utils"
)

//"image/jpeg"

type s3res struct {
	Name   string
	Reader io.Reader
}

//PrivateS3ImageProvider has a field for server as S3ImageProvider does not
type PrivateS3ImageProvider struct {
	bucket string
	prefix string
	id     string
	secret string
	server string
	region string
	hit    chan *s3res
}

//NewS3ImageProvider from privateS3store provider a ImageProvider from private s3 such as Scality or MinIO
func NewS3ImageProvider(server, id, secret, region, bucket, prefix string) *PrivateS3ImageProvider {

	log.Println("privates3store news3imageprovider")
	log.Println("id::" + string(id))
	log.Println("secret: " + string(secret))
	log.Println("server: " + string(server))
	log.Println("region: " + string(region))
	log.Println("bucket: " + string(bucket))

	conn, err := s3utils.GetS3Service(server, region, id, secret)
	if err != nil {
		log.Println("PrivateS3Store: ", err)
	}

	sss := &PrivateS3ImageProvider{
		id:     id,
		server: server,
		secret: secret,
		region: region,
		bucket: bucket,
		prefix: prefix,
		hit:    make(chan *s3res),
	}
	go sss.fetch(&conn)
	return sss
}

///*
////GetImage returns a Image from Private S3 store
//func (sss *PrivateS3ImageProvider) GetImage() (string, string, error) {
//	log.Println("start GetImage")
//	im, err := <-sss.hit
//	if err != nil {
//		log.Println(err)
//	}
//	if im == nil {
//		return "", "", errors.New("No new file")
//	}
//	buf := make([]byte, 4)
//	im.Reader.Read(buf)
//	return im.Name, buf, nil
//
//	//log.Println("getimage i ", i)
//	//log.Println("getimage ok ", ok)
//	//was working down
//	/*
//		if i, ok := <-sss.hit; ok {
//			// copy stream
//			im, _, err := image.Decode(i.Reader)
//			if err != nil {
//				log.Println(err)
//			}
//			if im == nil {
//				return "", "", errors.New("No new file")
//			}
//			var buff bytes.Buffer
//
//			if err := png.Encode(&buff, im); err != nil {
//				log.Println("err = ", err)
//				return "", "", errors.New("No new file")
//			}
//			return i.Name, "data:image/png;base64," + base64.StdEncoding.EncodeToString(buff.Bytes()), nil
//		}
//	*/
//	return "", "", errors.New("No new file")
//}
func (sss *PrivateS3ImageProvider) GetImage() (string, string, error) {
	if i, ok := <-sss.hit; ok {
		// copy stream
		fmt.Println(i.Name)
		im, _, err := image.Decode(i.Reader)
		if err != nil {
			log.Println(err)
		}
		var buff bytes.Buffer

		png.Encode(&buff, im)
		return i.Name, "data:image/png;base64," + base64.StdEncoding.EncodeToString(buff.Bytes()), nil
	}
	return "", "", errors.New("No new file")
}

// AddImage puts an image to private s3 store
func (sss *PrivateS3ImageProvider) AddImage(name string, url string) {
	var b []byte

	buff := bytes.NewBuffer(b)
	buff.WriteString(url)

	fmt.Println("fabio: ", "url ", url)

	sss.hit <- &s3res{
		Name:   name,
		Reader: buff,
	}
}

func (sss *PrivateS3ImageProvider) fetch(conn *s3utils.S3Connector) {
	buck := &s3.ListObjectsV2Input{}
	buck.SetBucket(sss.bucket)
	buck.SetPrefix(sss.prefix + "/")
	buck.SetDelimiter("/")

	sss.listThat(conn, buck)
}

func (sss *PrivateS3ImageProvider) listThat(conn *s3utils.S3Connector, buck *s3.ListObjectsV2Input) {
	//func (sss *PrivateS3ImageProvider) listThat(svc *s3.S3, buck *s3.ListObjectsInput) {
	prefixes := []string{}
	log.Println("privates3imageprovider listThat")
	walkFn := func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			if obj.Key == nil {
				continue
			}
			fmt.Println("key:", *obj.Key)

			in := s3.GetObjectInput{
				Bucket: buck.Bucket,
				Key:    obj.Key,
			}
			res, err := conn.Svc.GetObject(&in)
			if err != nil {
				log.Println("Private S3 ERR 1", err)
				continue
			}

			// AddImage...
			sss.hit <- &s3res{
				Name:   *obj.Key,
				Reader: res.Body,
			}
			for _, cp := range page.CommonPrefixes {
				prefixes = append(prefixes, *cp.Prefix)
			}
		}
		return false
	}
	listObjectsInput := &s3.ListObjectsV2Input{
		MaxKeys: aws.Int64(10),
		Bucket:  buck.Bucket,
	}
	conn.Svc.ListObjectsV2Pages(listObjectsInput, walkFn)
	for _, p := range prefixes {
		b := &s3.ListObjectsV2Input{}
		b.SetBucket(*buck.Bucket)
		b.SetPrefix(p)
		b.SetDelimiter("/")
		sss.listThat(conn, b)
	}

	/*
		prefixes := []string{}
		page := 0
			err := svc.ListObjectsV2Pages(buck, func(p *s3.ListObjectsV2Output, lastPage bool) bool {
				//err := svc.ListObjectsPages(buck, func(p *s3.ListObjectsOutput, lastPage bool) bool {
				page++
				for _, cc := range p.Contents {
					isImage := false
					for _, ext := range []string{".jpg", ".jpeg", ".png"} {
						k := strings.ToLower(*cc.Key)
						if strings.HasSuffix(k, ext) {
							isImage = true
						}
					}
					log.Println("privates3imageprovider listThat, cc %v", cc)
					log.Println("privates3imageprovider listThat, isImage %v", isImage)
					if !isImage {
						continue
					}

					in := s3.GetObjectInput{
						Bucket: buck.Bucket,
						Key:    cc.Key,
					}
					res, err := svc.GetObject(&in)
					if err != nil {
						log.Println("Private S3 ERR 1", err)
					}

					// AddImage...
					sss.hit <- &s3res{
						Name:   *cc.Key,
						Reader: res.Body,
					}
				}
				for _, cp := range p.CommonPrefixes {
					prefixes = append(prefixes, *cp.Prefix)
				}
				return lastPage
			})
			if err != nil {
				log.Println("Private S3 failed", err)
			}
			for _, p := range prefixes {
				b := &s3.ListObjectsV2Input{}
				//b := &s3.ListObjectsInput{}
				b.SetBucket(*buck.Bucket)
				b.SetPrefix(p)
				b.SetDelimiter("/")
				sss.listThat(svc, b)
			}
	*/
}
