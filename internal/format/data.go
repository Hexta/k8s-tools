package format

type Data struct {
	columns []interface{}
	rows    [][]interface{}
}

func NewData() *Data {
	return &Data{
		columns: make([]interface{}, 0),
		rows:    make([][]interface{}, 0),
	}
}

func (d *Data) AddColumn(column interface{}) {
	d.columns = append(d.columns, column)
}

func (d *Data) AddRow(row []interface{}) {
	d.rows = append(d.rows, row)
}

func (d *Data) GetColumns() []interface{} {
	return d.columns
}

func (d *Data) GetRows() [][]interface{} {
	return d.rows
}
