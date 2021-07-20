package parsers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/bitbucket"
)

func (d *Data) ParseBitbucket(c *gin.Context) error {
	bh, err := bitbucket.New()
	if err != nil {
		return err
	}

	var event bitbucket.Event

	payload, err := bh.Parse(c.Request, event)
	if err != nil {
		return err
	}

	switch p := payload.(type) {
	case bitbucket.RepoPushPayload:
		d.Deleted = d.BitbucketDeleted(p)

		if d.Deleted {
			d.Branch = ""
		} else {
			d.Branch = p.Push.Changes[0].New.Name
		}

		d.ModuleName = p.Repository.Name
		d.RepoName = p.Repository.FullName
		d.RepoUser = p.Repository.Owner.NickName
	default:
		return fmt.Errorf("unknown event type %s", payload)
	}

	return nil
}

func (d *Data) BitbucketDeleted(b bitbucket.RepoPushPayload) bool {
	return b.Push.Changes[0].Closed
}
