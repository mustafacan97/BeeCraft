package shared

type HALLink struct {
	Href      string `json:"href"`
	Method    string `json:"method,omitempty"`    // GET, POST, PUT, DELETE, vs.
	Type      string `json:"type,omitempty"`      // application/json, etc.
	Title     string `json:"title,omitempty"`     // An explanatory title to the user
	Templated bool   `json:"templated,omitempty"` // If true, it means that the URL contains a parameter
}

type HALResource struct {
	Links HALLinks `json:"_links"`
}

type HALLinks map[string]HALLink
