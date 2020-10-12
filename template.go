// DO NOT EDIT.

package main

import "text/template"

var tmpl = template.Must(template.New("template").Delims(`@@`, `@@`).Parse("-- Dockerfile --\nFROM golang:1.14.2-alpine3.11\n\nENV GO111MODULE=on\n\nWORKDIR /app\nCOPY go.mod .\n\nRUN go mod tidy\nCOPY . .\n-- README.md --\n# QUIC START\n\n- this app was made by github.com/hagarihayato/sprint\n-- application/README.md --\n### Application layer\n\n### OverView\n\n- this layer is like UseCase . This layer receives information from the repository layer.\n\n### How to use\n\n#### Create Application Handler\n\n- At first: Create a function that belongs to a structure\n```\nex...\n\nfunc (e *exampleApplication) GetExample() (*model.Example, error) {\n    // abridgement\n}\n```\n\n- And then: Fill the interface\n```\nex...\n\ntype ExampleApplication interface {\n    GetExample() (*model.Example, error) // additional codes\n}\n```\n-- application/sample_application.go --\npackage application\n\nimport (\n\t\"@@.ImportPath@@/domain/model\"\n\t\"@@.ImportPath@@/domain/repository\"\n)\n\ntype (\n\tsampleApplication struct {\n\t\trepository.SampleRepository\n\t}\n\tSampleApplication interface {\n\t\tGetSamples() ([]*model.Sample, error)\n\t}\n)\n\nfunc NewSampleApplication(rs repository.SampleRepository) SampleApplication {\n\treturn &sampleApplication{rs}\n}\n\nfunc (s *sampleApplication) GetSamples() ([]*model.Sample, error) {\n\treturn s.SampleRepository.Fetch()\n}\n-- build.sh --\n#!/bin/bash\n\ngo build -o app && ./app\n-- docker-compose.yml --\nversion: \"3.5\"\nservices:\n  @@.ImportPath@@:\n    container_name: @@.ImportPath@@\n    build: .\n    tty: true\n    restart: always\n    volumes:\n      - .:/app\n    depends_on:\n      - @@.ImportPath@@db\n    ports:\n      - 8080:8080\n    command: sh ./build.sh\n  @@.ImportPath@@db:\n    image: postgres:10-alpine\n    container_name: @@.ImportPath@@db\n    ports:\n      - \"5432:5432\"\n    environment:\n      - POSTGRES_USER=postgres\n      - POSTGRES_PASSWORD=postgres\n      - PGPASSWORD=postgres\n      - POSTGRES_DB=postgres\n      - DATABASE_HOST=localhost\n    volumes:\n      - $PWD/infrastructure/postgres/init:/docker-entrypoint-initdb.d\n      #      this volume allows for data persistence; if you make data persistent, make /infrastructure/postgres/data as empty directory and you remove comment out from this volume.\n      #      - $PWD/infrastructure/postgres/data:/var/lib/postgresql/data\n-- domain/model/sample_model.go --\npackage model\n\ntype Sample struct {\n\tID int64\n\tText string\n}\n-- domain/repository/README.md --\n### Application layer\n\n### OverView\n\n- this layer is like UseCase . This layer receives information from the repository layer.\n\n### How to use\n\n#### Create Repository Handler\n\n- At first: Create a function that belongs to a structure\n```\nex...\n\nfunc (e *exampleRepository) Fetch() (*model.Example, error) {\n    // abridgement\n}\n```\n\n- And then: Fill the interface\n```\nex...\n\ntype ExampleRepository interface {\n    Fetch() (*model.Example, error) // additional codes\n}\n```\n-- domain/repository/sample_repository.go --\npackage repository\n\nimport (\n\t\"database/sql\"\n\t\"@@.ImportPath@@/domain/model\"\n)\n\ntype (\n\tsampleRepository struct{\n\t\tconn *sql.DB\n\t}\n\tSampleRepository interface {\n\t\tFetch() ([]*model.Sample, error)\n\t}\n)\n\nfunc NewSampleRepository(Conn *sql.DB) SampleRepository {\n\treturn &sampleRepository{Conn}\n}\n\nfunc (s *sampleRepository) Fetch() ([]*model.Sample, error) {\n\tvar samples []*model.Sample\n\trows, err := s.conn.Query(\"SELECT id, text FROM samples;\")\n\tif rows == nil { return nil, err }\n\tfor rows.Next() {\n\t\tsample := &model.Sample{}\n\t\terr = rows.Scan(&sample.ID, &sample.Text)\n\t\tif err == nil {\n\t\t\tsamples = append(samples, sample)\n\t\t}\n\t}\n\treturn samples, err\n}\n-- go.mod --\nmodule @@.ImportPath@@\n\ngo @@.GoVer@@\n-- infrastructure/postgres/conf/database.go --\npackage conf\n\nimport (\n\t\"database/sql\"\n\t\"fmt\"\n\t_ \"github.com/lib/pq\"\n)\n\nvar (\n\tDRIVER = \"postgres\"\n\tHOSTNAME = \"@@.ImportPath@@db\"\n\tUSER = \"postgres\"\n\tDBNAME = \"postgres\"\n\tPASSWORD = \"postgres\"\n)\n\nfunc NewDatabaseConnection() (*sql.DB, error) {\n\tsource := fmt.Sprintf(\"host=%s user=%s dbname=%s password=%s sslmode=disable\", HOSTNAME, USER, DBNAME, PASSWORD)\n\tconn, err := sql.Open(DRIVER, source)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn conn, nil\n}\n\n-- infrastructure/postgres/init/1_init.sql --\nDROP TABLE IF EXISTS samples;\n\nCREATE TABLE IF NOT EXISTS samples\n(\n   id SERIAL NOT NULL,\n   text TEXT NOT NULL,\n   PRIMARY KEY (id)\n);\n\nINSERT INTO samples(text) VALUES ('sample');\n-- interactor/README.md --\n### OverView\n\n- This file is a DI(Dependency Injection) container\n\n- It has each layer of Repository, Application, and Presenter as a structure.\nEach structure is initialized in the NewXXX method of that structure.\n\n### How to use\n\n- When adding a new model structure\n\n```\nex... \n// When example is added to model\ntype Example struct {\n    Name string\n}\n```\n\nadd example_repository.go in /domain/repository/ \n\n```\n// example_repository.go\npackage repository\n\nimport (\n\t\"database/sql\"\n\t\"$YourProjectName/domain/model\"\n)\n\ntype (\n\texampleRepository struct{\n\t\tconn *sql.DB\n\t}\n\tExampleRepository interface {\n\t\tFetch() ([]*model.Example, error)\n\t}\n)\n\nfunc NewExampleRepository(Conn *sql.DB) ExampleRepository {\n\treturn &exampleRepository{Conn}\n}\n\nfunc (s *exampleRepository) Fetch() ([]*model.Example, error) {\n    // TODO\n}\n```\n\nadd example_application.go in /application/\n\n```\npackage application\n\nimport (\n\t\"$YourProjectName/domain/model\"\n\t\"$YourProjectName/domain/repository\"\n)\n\ntype (\n\texampleApplication struct {\n\t\trepository.ExampleRepository\n\t}\n\tExampleApplication interface {\n\t\tGetExamples() ([]*model.Example, error)\n\t}\n)\n\nfunc NewExampleApplication(rs repository.ExampleRepository) ExampleApplication {\n\treturn &exampleApplication{rs}\n}\n\nfunc (s *exampleApplication) GetExamples() ([]*model.Example, error) {\n\treturn s.ExampleRepository.Fetch()\n}\n```\n\nadd example_handler.go in /presenter/handler/\n\n```\npackage handler\n\nimport (\n\t\"encoding/json\"\n\t\"log\"\n\t\"net/http\"\n\t\"$YourProjectName/application\"\n\t\"$YourProjectName/domain/model\"\n)\n\ntype(\n\texampleHandler struct {\n\t\tapplication.ExampleApplication\n\t}\n\tExampleHandler interface {\n\t\tExampleIndex(w http.ResponseWriter, r *http.Request)\n\t}\n\tresponse struct {\n\t\tStatus int\n\t\tExamples []*model.Example\n\t}\n)\n\nfunc NewExampleHandler(as application.ExampleApplication) ExampleHandler {\n\treturn &exampleHandler{as}\n}\n\nfunc (s *exampleHandler) ExampleIndex(w http.ResponseWriter, r *http.Request) {\n    // TODO\n}\n```\n\nwrite additional interfaces in interactor.go\n\n```\nex...\n\ntype Repository struct {\n\trepository.SampleRepository\n    repository.ExampleRepository // additional codes\n}\n\n...\n\nfunc (i *interactor) NewRepository() *Repository {\n\tr := &Repository{}\n\tr.SampleRepository = repository.NewSampleRepository(i.conn)\n\tr.ExampleRepository = repository.NewExampleRepository(i.conn)\n\treturn r\n}\n\n...\n\n```\n-- interactor/interactor.go --\npackage interactor\n\nimport (\n\t\"database/sql\"\n\t\"@@.ImportPath@@/application\"\n\t\"@@.ImportPath@@/domain/repository\"\n\t\"@@.ImportPath@@/presenter/handler\"\n)\n\ntype (\n\tinteractor struct {\n\t\tconn *sql.DB\n\t}\n\tInteractor interface {\n\t\tNewRepository() *Repository\n\t\tNewApplication(r *Repository) *Application\n\t\tNewHandler(a *Application) *Handler\n\t}\n\tRepository struct {\n\t\trepository.SampleRepository\n\t}\n\tApplication struct {\n\t\tapplication.SampleApplication\n\t}\n\tHandler struct {\n\t\thandler.SampleHandler\n\t}\n)\n\nfunc NewInteractor(conn *sql.DB) Interactor {\n\treturn &interactor{conn}\n}\n\nfunc (i *interactor) NewRepository() *Repository {\n\tr := &Repository{}\n\tr.SampleRepository = repository.NewSampleRepository(i.conn)\n\treturn r\n}\n\nfunc (i *interactor) NewApplication(r *Repository) *Application {\n\ta := &Application{}\n\ta.SampleApplication = application.NewSampleApplication(r.SampleRepository)\n\treturn a\n}\n\nfunc (i *interactor) NewHandler(a *Application) *Handler {\n\th := &Handler{}\n\th.SampleHandler = handler.NewSampleHandler(a.SampleApplication)\n\treturn h\n}\n\n\n\n\n\n-- main.go --\npackage main\n\nimport (\n\t\"fmt\"\n\t\"net/http\"\n\t\"@@.ImportPath@@/infrastructure/postgres/conf\"\n\t\"@@.ImportPath@@/interactor\"\n\t\"@@.ImportPath@@/presenter/middleware\"\n\t\"@@.ImportPath@@/presenter/router\"\n)\n\nfunc main() {\n\tconn, err := conf.NewDatabaseConnection()\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer conn.Close()\n\tfmt.Println(`\n    * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * \n  *        ####    #####    #####     ####    ##  ##   ######   ######   #####  *\n  *      ##  ##   ##  ##   ##  ##     ##     ### ##     ##     ##       ##  ##  *\n  *     ##       ##  ##   ##  ##     ##     ######     ##     ##       ##  ##   *\n  *     ####    #####    #####      ##     ######     ##     ####     #####     *\n  *       ##   ##       ####       ##     ## ###     ##     ##       ####       *\n  *  ##  ##   ##       ## ##      ##     ##  ##     ##     ##       ## ##       *\n  *  ####    ##       ##  ##    ####    ##  ##     ##     ######   ##  ##       *\n    * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *\n    `)\n\n\tfmt.Println(`HTML:\tGET http://localhost:8080`)\n\tfmt.Println(`API:\tGET http://localhost:8080/api/v1`)\n\n\ti := interactor.NewInteractor(conn)\n\tr := i.NewRepository()\n\ta := i.NewApplication(r)\n\th := i.NewHandler(a)\n\tm := middleware.NewMiddleware()\n\ts := router.NewRouter()\n\ts.Router(h, m)\n\n\t_ = http.ListenAndServe(\":8080\", s.Route)\n}\n-- presenter/handler/README.md --\n### Presenter layer\n\n### OverView\n\n- This layer receives information from the application layer.\n\n### How to use\n\n#### Create Presenter Handler\n\n- At first: Create a function that belongs to a structure\n```\nex...\n\nfunc (e *exampleHandler) ExampleIndex() (*model.Example, error) {\n    // abridgement\n}\n```\n\n- And then: Fill the interface\n```\nex...\n\ntype ExampleHandler interface {\n    ExampleIndex() (*model.Example, error) // additional codes\n}\n```\n-- presenter/handler/handler_util.go --\npackage handler\n\nimport (\n\t\"fmt\"\n\t\"html/template\"\n\t\"io/ioutil\"\n\t\"os\"\n\t\"path/filepath\"\n\t\"strings\"\n)\n\ntype response struct {\n\tStatus int\n\tResult interface{}\n}\n\nfunc parseTemplate(dir string, fileName string) (*template.Template, error) {\n\ttmpl := template.New(\"\")\n\n\tvar layout string\n\n\tif err := filepath.Walk(\"presenter/template/layout\", func(path string, info os.FileInfo, err error) error {\n\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\tif !info.IsDir() && (strings.HasSuffix(path, \".html\") || strings.HasSuffix(path, \".js\")) {\n\t\t\tfile, err := ioutil.ReadFile(path)\n\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\n\t\t\tlayout += string(file)\n\t\t}\n\n\t\treturn nil\n\t}); err != nil {\n\t\tfmt.Println(err)\n\t}\n\n\tif err := filepath.Walk(\"presenter/template/\" + dir, func(path string, info os.FileInfo, err error) error {\n\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\n\t\tif !info.IsDir() && (strings.HasSuffix(path, \".html\") || strings.HasSuffix(path, \".js\")) {\n\t\t\tfile, err := ioutil.ReadFile(path)\n\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\n\t\t\tfilename := strings.Replace(path, \"presenter/template/\" + dir, \"\", -1)\n\n\n\t\t\tif strings.Contains(filename, fileName) {\n\t\t\t\ttmpl = tmpl.New(filename)\n\n\t\t\t\ttmpl, err = tmpl.Parse(string(file) + layout)\n\n\t\t\t\tif err != nil {\n\t\t\t\t\treturn err\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\n\t\treturn nil\n\t}); err != nil {\n\t\treturn nil, err\n\t}\n\n\treturn tmpl, nil\n}\n\n-- presenter/handler/sample_handler.go --\npackage handler\n\nimport (\n\t\"@@.ImportPath@@/application\"\n\t\"encoding/json\"\n\t\"log\"\n\t\"net/http\"\n)\n\ntype(\n\tsampleHandler struct {\n\t\tapplication.SampleApplication\n\t}\n\tSampleHandler interface {\n\t\tSampleIndex(w http.ResponseWriter, r *http.Request)\n\t\tSampleHTML(w http.ResponseWriter, r *http.Request)\n\t}\n)\n\nfunc NewSampleHandler(as application.SampleApplication) SampleHandler {\n\treturn &sampleHandler{as}\n}\n\nfunc (s *sampleHandler) SampleIndex(w http.ResponseWriter, r *http.Request) {\n\tsamples, err := s.SampleApplication.GetSamples()\n\n\tif err != nil {\n\t\thttp.Error(w, err.Error(), http.StatusNotFound)\n\t}\n\n\tresp := &response{\n\t\tStatus: http.StatusOK,\n\t\tResult: samples,\n\t}\n\n\tres, err := json.Marshal(resp)\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\n\t_ , _ = w.Write(res)\n}\n\nfunc (s *sampleHandler) SampleHTML(w http.ResponseWriter, r *http.Request) {\n\ttmpl, err := parseTemplate(\"sample\", \"index\")\n\n\tif err != nil {\n\t\tlog.Fatal(\"err :\", err)\n\t}\n\n\tif err := tmpl.Execute(w, nil); err != nil {\n\t\tlog.Printf(\"failed to execute template: %v\", err)\n\t}\n}\n-- presenter/middleware/main.go --\npackage middleware\n\ntype middleware struct {}\n\ntype Middleware interface {}\n\nfunc NewMiddleware() Middleware {\n\treturn &middleware{}\n}\n-- presenter/router/router.go --\npackage router\n\nimport (\n\t\"github.com/go-chi/chi\"\n\t\"github.com/go-chi/chi/middleware\"\n\t\"@@.ImportPath@@/interactor\"\n\tmid \"@@.ImportPath@@/presenter/middleware\"\n)\n\ntype Server struct {\n\tRoute *chi.Mux\n}\n\nfunc NewRouter() *Server {\n\treturn &Server{\n\t\tRoute: chi.NewRouter(),\n\t}\n}\n\nfunc (s *Server) Router(h *interactor.Handler, m mid.Middleware) {\n\ts.Route.Use(middleware.Logger)\n\ts.Route.Use(middleware.Recoverer)\n\ts.Route.Route(\"/\", func(r chi.Router) {\n\t\tr.Get(\"/\", h.SampleHandler.SampleHTML)\n\t})\n\ts.Route.Route(\"/api/v1\", func(r chi.Router) {\n\t\tr.Get(\"/\", h.SampleHandler.SampleIndex)\n\t\t// TODO\n\t})\n}\n-- presenter/router/router_test.go --\npackage router_test\n\nimport (\n\t\"github.com/go-chi/chi\"\n\t\"github.com/stretchr/testify/assert\"\n\t\"reflect\"\n\t\"@@.ImportPath@@/presenter/router\"\n\t\"testing\"\n)\n\nfunc TestNewRouter(t *testing.T) {\n\troute := chi.NewRouter()\n\twantNewRouterForm := &router.Server{\n\t\tRoute: route,\n\t}\n\n\tr := router.NewRouter()\n\n\tv := reflect.ValueOf(r)\n\tw := reflect.ValueOf(wantNewRouterForm)\n\n\tassert.Equal(t, v.Type(), w.Type())\n}\n-- presenter/template/layout/_footer.html --\n{{ define \"footer\" }}\n    <footer>sample</footer>\n    </body>\n    </html>\n{{ end }}\n-- presenter/template/layout/_header.html --\n{{ define \"header\" }}\n    <!DOCTYPE html>\n    <html lang=\"en\">\n    <head>\n        <meta charset=\"UTF-8\">\n        <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n        <title>Sample</title>\n        <style>\n            * {\n                margin: 0;\n                padding: 0;\n                box-sizing: border-box;\n            }\n        </style>\n    </head>\n    <body>\n    <header><a href=\"/\">sample</a></header>\n{{ end }}\n-- presenter/template/sample/index.html --\n{{ template \"header\" }}\n<h1><i>SPRINTER</i></h1>\n{{ template \"footer\" }}\n"))
