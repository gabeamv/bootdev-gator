package gatorcommand

const (
	GATOR     = "gator"
	LOGIN     = "login"
	REGISTER  = "register"
	RESET     = "reset"
	USERS     = "users"
	AGG       = "agg"
	ADDFEED   = "addfeed"
	FEEDS     = "feeds"
	FOLLOW    = "follow"
	FOLLOWING = "following"
	UNFOLLOW  = "unfollow"
	BROWSE    = "browse"
	HELP      = "help"
)

const (
	LOGIN_DESCRIBE     = "login [arg0] ... logs in a registered user"
	REGISTER_DESCRIBE  = "register [arg0] ... registers a new unique user i.e. register {new user}"
	RESET_DESCRIBE     = "reset... removes all registered users"
	USERS_DESCRIBE     = "users ... lists all registered users and highlights the current one logged in"
	AGG_DESCRIBE       = "agg [arg0] ... iteratively update the oldest feed by requesting new posts for that feed. user inputs how many seconds they want to update by"
	ADDFEED_DESCRIBE   = "addfeed [arg0, arg1] ... add feed to update to database by name of feed and url"
	FEEDS_DESCRIBE     = "feeds ... list all active feeds"
	FOLLOW_DESCRIBE    = "follow [arg0] ... follows a feed by url"
	FOLLOWING_DESCRIBE = "following ... lists all feeds followed by current user"
	UNFOLLOW_DESCRIBE  = "unfollow [arg0] ... unfollows a feed by url for current user"
	BROWSE_DESCRIBE    = "browse [arg0] ... browses latest published feeds. number of latest feeds given by user (optional)"
	HELP_DESCRIBE      = "help ... lists all possible commands, their descriptions and parameters"
)
