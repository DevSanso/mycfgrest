package handle

type handleMetaData struct {
	Url string `toml:"url"`
	
	Request struct {
		Method      string   `toml:"method"`
		ContentType []string `toml:"content_type"`

		QueryString map[string]struct {
			Type   string `toml:"type"`
			Symbol string `toml:"symbol"`
		} `toml:"query_string"`

		Body map[string]struct {
			Type   string `toml:"type"`
			Key    string `toml:"key"`
			Symbol string `toml:"symbol"`
		}
	} `toml:"request"`

	Load map[string]struct {
		Type     string `toml:"type"`
		LoadName string `toml:"load_name"`
		Command  string `toml:"command"`

		GetData map[string]struct {
			Type   string `toml:"type"`
			Symbol string `toml:"symbol"`
		} `toml:"get_data"`
	} `toml:"load"`

	Response struct {
		ContentType string `toml:"content_type"`
		Template    string `toml:"template"`
	} `toml:"response"`
}

type HandleMeta struct {
	Data handleMetaData `toml:"handle"`
}
