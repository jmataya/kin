package kin

// func TestRowResult(t *testing.T) {
// 	getCurrentFile()
// 	t.Error("Error")
// }

// func getCurrentFile() {
// 	_, filename, _, _ := runtime.Caller(1)
// 	fmt.Printf("Filename: %s\n", filename)

// 	fset := token.NewFileSet()
// 	astFile, err := parser.ParseFile(fset, filename, nil, parser.PackageClauseOnly)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("AST: %+v\n", astFile)
// }

// func packageName(file string) (string, error) {
// 	fset := token.NewFileSet()

// 	// parse the go soure file, but only the package clause
// 	astFile, err := parser.ParseFile(fset, l.path, nil, parser.PackageClauseOnly)
// 	if err != nil {
// 		return "", err
// 	}

// 	if astFile.Name == nil {
// 		return "", fmt.Errorf("no package name found")
// 	}

// 	return "", astFile.Name.Name
// }

// const (
// 	createBuilderTable = `
// 		create table builders (
// 			id serial primary key,
// 			name text not null,
// 			attributes jsonb not null default '{}',
// 			is_active boolean not null default false,
// 			created_at timestamp without time zone
// 		);
// 	`

// 	sqlInsertBuilder = `
// 		INSERT INTO builders (name, attributes, is_active, created_at)
// 		VALUES ($1, $2, $3, $4)
// 		RETURNING *
// 	`
// )

// type builderModel struct {
// 	ID         int
// 	Name       string
// 	Attributes map[string]interface{}
// 	IsActive   bool
// 	CreatedAt  time.Time
// }

// func (b *builderModel) Columns() []FieldBuilder {
// 	return []FieldBuilder{
// 		IntField("id", &b.ID),
// 		StringField("name", &b.Name),
// 		JSONField("attributes", &b.Attributes),
// 		BoolField("is_active", &b.IsActive),
// 		TimeField("created_at", &b.CreatedAt),
// 	}
// }

// func TestBuild(t *testing.T) {
// 	migrationPath := "./sql"
// 	cleanupMigrationDir(migrationPath)
// 	if err := setupMigrationDir(migrationPath); err != nil {
// 		t.Errorf("setupMigrationDir = %v", err)
// 		return
// 	}

// 	if err := createFile(migrationPath, "1__create_builders.sql", createBuilderTable); err != nil {
// 		t.Errorf("createFile = %v", err)
// 		cleanupMigrationDir(migrationPath)
// 		return
// 	}

// 	connStr := os.Getenv("POSTGRES_URL")
// 	migrator, _ := NewMigratorConnection(connStr)
// 	defer migrator.Close()

// 	if err := migrator.Migrate(migrationPath); err != nil {
// 		t.Errorf("migrator.Migrate(%s) = %v, want nil", migrationPath, err)
// 		cleanupMigrationDir(migrationPath)
// 		return
// 	}

// 	name := "Builder Test"
// 	attributes := `{ "lang": "en" }`
// 	isActive := true
// 	createdAt := time.Now().UTC()

// 	builder := new(builderModel)

// 	db, _ := NewConnection(connStr)
// 	defer db.Close()
// 	err := db.Query(sqlInsertBuilder, name, attributes, isActive, createdAt).OneAndExtract(builder)
// 	if err != nil {
// 		t.Errorf("db.Query(...).OneAndExtract(...) = %v, want <nil>", err)
// 		return
// 	}

// 	if builder.ID == 0 {
// 		t.Error("builder.ID = 0, want > 0")
// 	}

// 	if builder.Name != name {
// 		t.Errorf("builder.Name = %v, want %v", builder.Name, name)
// 	}

// 	if builder.IsActive != isActive {
// 		t.Errorf("builder.IsActive = %v, want %v", builder.IsActive, isActive)
// 	}

// 	if builder.CreatedAt.Equal(createdAt) {
// 		t.Errorf("builder.CreatedAt = %v, want %v", builder.CreatedAt, createdAt)
// 	}

// 	if builder.Attributes["lang"] != "en" {
// 		t.Error(`"builder.Attributes != { "lang": "en" }`)
// 	}

// 	cleanupMigrationDir(migrationPath)
// }
