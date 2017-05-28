package plugins

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/arachnist/repost-plugin-server"
	"github.com/arachnist/repost-plugin-server/util"
)

func bonjour(ctx context.Context, sender, recipient string, message []string) (rmsg []string) {
	t, _ := time.Parse("2006-01-02", "2015-12-01")
	max := int(time.Now().Sub(t).Hours())/24 + 1

	img, err := util.HTTPGetXpath("http://ditesbonjouralamadame.tumblr.com/page/"+fmt.Sprintf("%d", rand.Intn(max)+1), "//div[@class='photo post']//a/@href")
	if err != nil {
		rmsg = []string{"error:", err.Error()}
	} else {
		rmsg = []string{"bonjour", "(nsfw):", img}
	}

	return rmsg
}

func init() {
	rps.Register("bonjour", bonjour)
}
