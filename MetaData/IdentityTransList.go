package MetaData

type IdentityTransListPair struct {
	IdentityTransList IdentityTransformation
	Height            int
}

type IdentityTransList struct {
	IdentityTransList []IdentityTransListPair
}

func (id *IdentityTransList) SetIdentityTransList(m map[IdentityTransformation]int) {
	identityTransList := make([]IdentityTransListPair, 0)

	for k, v := range m {
		var pair IdentityTransListPair
		pair.IdentityTransList = k
		pair.Height = v
		identityTransList = append(identityTransList, pair)
	}
	id.IdentityTransList = identityTransList
}

func (id *IdentityTransList) GetIdentityTransList() map[IdentityTransformation]int {
	identityList := make(map[IdentityTransformation]int)

	for _, pair := range id.IdentityTransList {
		identityList[pair.IdentityTransList] = pair.Height
	}
	return identityList
}
