package conn

type connInfo struct {
	Addr string `toml:"addr"`
	Port int `toml:"port"`
	User string `toml:"user"`
	Password string `toml:"password"`
	Dbname string `toml:"dbname"`
}

type ConnMeta struct {
	Sql struct {
		Postgres map[string]connInfo `toml:"postgres"`
		Sqlite map[string]connInfo `toml:"sqlite"`
	} `toml:"sql"`
}