package api

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/parsers"
)

type ModuleController struct{}

func (m ModuleController) DeployModule(c *gin.Context) {
	data := parsers.Data{}
	cmd := exec.Command("r10k", "deploy", "module")
	conf := config.GetConfig()

	err := data.ParseData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error Parsing Webhook", "error": err})
		c.Abort()
		return
	}

	cmd.Args = append(cmd.Args, data.ModuleName)

	if conf.GetBool("verbose") {
		cmd.Args = append(cmd.Args, "-v")
	}

	res, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("cmd.Run() failed with error %s", string(res))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error executing command", "error": string(res)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": string(res)})
	log.Info(fmt.Sprintf("\n%s", string(res)))

}
