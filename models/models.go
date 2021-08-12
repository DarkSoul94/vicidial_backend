package models

type Lead map[string]interface{}

func (l *Lead) Get(key string, alterValue interface{}) interface{} {
	temp := *l
	if _, ok := temp[key]; !ok {
		return alterValue
	}
	return temp[key]
}
