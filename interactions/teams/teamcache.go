package teams

import "tacobot/util"

var cache = util.NewCache[string, []string]("TeamCache")
