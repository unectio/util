package k8s

type DepDesc struct {
	proj		string
	lang		string
	class		string
}

func (dd *DepDesc)Proj() string { return dd.proj }
func (dd *DepDesc)Lang() string { return dd.lang }
func (dd *DepDesc)Class() string { return dd.class }

func MkDepDesc(proj, lang, class string) *DepDesc {
	return &DepDesc{
		proj:	proj,
		lang:	lang,
		class:	class,
	}
}

func (dd *DepDesc)Image(pfx string) string {
	return pfx + "/" + dd.lang
}

func (dd *DepDesc)Name() string {
	return "dep-" + dd.Key()
}

func (dd *DepDesc)Key() string {
	/*
	 * We keep deployment for each project's language
	 */
	return dd.proj + "-" + dd.lang + "-" + dd.class
}
