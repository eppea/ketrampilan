module main

go 1.20
// Migrate tabel murid
db.AutoMigrate(&Murid{})

return db, nil
