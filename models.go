package main

type TodoItemModel struct {
	Id          int `gorm:"primary_key"`
	Description string
	Completed   bool
}

func (a *App) Migrations() {
	// defer a.DB.Close()

	// a.DB.Debug().DropTableIfExists(&TodoItemModel{})
	a.DB.Debug().AutoMigrate(&TodoItemModel{})
}
