package sql

type Class struct {
	MapWhere    map[string]string
	MapSet      map[string]string
	StringTable string
}

func New() Class {
	this := Class{}
	this.MapWhere = map[string]string{}
	this.MapSet = map[string]string{}
	return this
}

func (this *Class) Table(stringTable string) *Class {
	this.StringTable = stringTable
	return this
}

func (this *Class) Where(mapWhere map[string]string) *Class {
	this.MapWhere = mapWhere
	return this
}

func (this *Class) WhereAnd(name string, value string) *Class {
	this.MapWhere[name] = value
	return this
}

func (this *Class) Set(name string, value string) *Class {
	this.MapSet[name] = value
	return this
}

func (this *Class) Save(MapSet map[string]string) error {
	this.MapSet = MapSet
	return nil
}
