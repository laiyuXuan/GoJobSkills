package constant

import "time"

//config
var REDIS_SERVER = "127.0.0.1:6379"
var MAX_IDLE = 3
var IDLE_TIMEOUT = 240 * time.Second

//key
var KEY_JIANSHU_ARTICLES_LINKS = "jianshu_article_links"