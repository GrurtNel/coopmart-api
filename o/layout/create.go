package layout

func (l *Layout) Create() error {
	return layoutTable.Create(l)
}

// func GetLayouts()  {
// 	var layouts []*Layout
// 	layoutTable.FindWhere(query bson.M, result interface{})
// 	return
// }
