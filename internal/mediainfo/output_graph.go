package mediainfo

func RenderGraphSVG(reports []Report) string {
	_ = reports
	return "<svg xmlns=\"http://www.w3.org/2000/svg\"></svg>"
}

func RenderGraphDOT(reports []Report) string {
	_ = reports
	return "digraph mediainfo {}"
}
