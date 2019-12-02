package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smileinnovation/imannotate/api/annotation"
	"github.com/smileinnovation/imannotate/api/annotation/exporter"
	"github.com/smileinnovation/imannotate/api/auth"
	"github.com/smileinnovation/imannotate/api/project"
	"github.com/smileinnovation/imannotate/app/registry"
)

func GetProjects(c *gin.Context) {
	u := auth.GetCurrentUser(c.Request)
	if u == nil {
		c.JSON(http.StatusUnauthorized, errors.New("No user found"))
		return
	}
	projects := project.GetAll(u)
	if len(projects) == 0 {
		c.JSON(http.StatusNotFound, projects)
		return
	}
	c.JSON(http.StatusOK, projects)
}

func GetProject(c *gin.Context) {
	name := c.Param("name")
	status := http.StatusOK
	p := project.Get(name)

	if p == nil {
		status = http.StatusNotFound
	}
	c.JSON(status, p)
}

func NewProject(c *gin.Context) {
	p := &project.Project{}
	err := c.Bind(p)
	if err != nil {
	}

	err = project.New(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, p)
}

func UpdateProject(c *gin.Context) {
	p := &project.Project{}
	c.Bind(p)
	id := p.Id

	registry.RemoveProvider(p)

	if err := project.Update(p); err != nil {
		c.JSON(http.StatusNotAcceptable, err.Error())
		return
	}

	// get updated project
	p = project.Get(id)
	c.JSON(http.StatusOK, p)
}

func GetNextImage(c *gin.Context) {
	p := project.Get(c.Param("name"))
	name, image, _ := project.NextImage(p)

	c.JSON(http.StatusOK, map[string]string{
		"name": name,
		"url":  image,
	})
}

func SaveAnnotation(c *gin.Context) {
	ann := &annotation.Annotation{}
	u := auth.GetCurrentUser(c.Request)

	c.Bind(ann)

	// TODO: use id
	pid := c.Param("name")
	prj := project.Get(pid)
	if !project.CanAnnotate(u, prj) {
		c.JSON(http.StatusUnauthorized, errors.New("You're not allowed to annotate image on that project"))
		return
	}
	if err := annotation.Save(prj, ann); err != nil {
		c.JSON(http.StatusNotAcceptable, err.Error())
		return
	}

	c.JSON(http.StatusAccepted, ann)
}

func GetContributors(c *gin.Context) {
	pid := c.Param("name")
	prj := project.Get(pid)

	users := project.GetContributors(prj)
	status := http.StatusOK
	if len(users) == 0 {
		status = http.StatusNotFound
	}
	c.JSON(status, users)
}

func AddContributor(c *gin.Context) {
	pid := c.Param("name")
	uid := c.Param("user")
	prj := project.Get(pid)
	user, err := auth.Get(uid)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := project.AddContributor(user, prj); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else {
		c.JSON(http.StatusCreated, "injected")
	}

}

func RemoveContributor(c *gin.Context) {
	pid := c.Param("name")
	uid := c.Param("user")
	prj := project.Get(pid)
	user, err := auth.Get(uid)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := project.RemoveContributor(user, prj); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else {
		c.JSON(http.StatusCreated, "removed")
	}
}

func ExportProject(c *gin.Context) {
	typ := c.Param("format")
	pid := c.Param("name")
	prj := project.Get(pid)

	var exp exporter.Exporter
	switch typ {
	case "csv":
		exp = &exporter.CSV{}
	}

	ann := annotation.Get(prj)
	reader := exp.Export(ann)

	c.Status(200)
	io.Copy(c.Writer, reader)
}

func DeleteProject(c *gin.Context) {
	pname := c.Param("name")
	prj := project.Get(pname)
	if err := project.Delete(prj); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, pname+" deleted")

}

func ProjectInfo(c *gin.Context) {
	pname := c.Param("name")
	prj := project.Get(pname)

	if prj == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	count := project.GetStats().CountAnnotations(prj)
	ctrb := len(project.GetContributors(prj))
	c.JSON(http.StatusOK, map[string]interface{}{
		"annotations":  count,
		"contributors": ctrb,
	})
}
