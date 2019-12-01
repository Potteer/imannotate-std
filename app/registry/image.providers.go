package registry

import (
	"log"
	"fmt"

	"github.com/smileinnovation/imannotate/api/project"
	"github.com/smileinnovation/imannotate/api/providers"
	"github.com/smileinnovation/imannotate/app/providers/memorygc"
	"github.com/smileinnovation/imannotate/app/providers/privates3store"
	"github.com/smileinnovation/imannotate/app/providers/qwant"
	"github.com/smileinnovation/imannotate/app/providers/s3store"
)

var imageProviders = make(map[string]providers.ImageProvider)

func GetProvider(prj *project.Project) providers.ImageProvider {
	fmt.Println("getprovider")

	if prj.ImageProvider == "" {
		log.Println("No image provider defined")
		return nil
	}
	fmt.Println("ip", imageProviders[prj.Id])
	if p, ok := imageProviders[prj.Id]; !ok {
		createImageProvider(prj)
		return GetProvider(prj)
	} else {
		gc := memorygc.NewMemoryGC()
		gc.SetImageProvider(p)
		SetGC(prj, gc)
		fmt.Println("return p")
		return p
	}
}

func RemoveProvider(prj *project.Project) {
	delete(imageProviders, prj.Id)
}

func createImageProvider(prj *project.Project) {
	opt := prj.ImageProviderOptions

	switch prj.ImageProvider {
	case "qwant":
		provider := qwant.NewQwant(opt["qwantQuery"])
		imageProviders[prj.Id] = provider
	case "s3":
		provider := s3store.NewS3ImageProvider(opt["id"], opt["secret"], opt["region"], opt["bucket"], opt["prefix"])
		imageProviders[prj.Id] = provider
	case "scality":
		provider := privates3store.NewS3ImageProvider(opt["server"], opt["id"], opt["secret"], opt["region"], opt["bucket"], opt["prefix"])
		imageProviders[prj.Id] = provider
	}
}
