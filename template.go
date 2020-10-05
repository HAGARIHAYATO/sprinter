// Code generated by _tools/txtar/main.go; DO NOT EDIT.

package main

import "text/template"

var tmpl = template.Must(template.New("template").Delims(`@@`, `@@`).Parse("-- Dockerfile --\nFROM golang:1.14.2-alpine3.11\n\nENV GO111MODULE=on\n\nWORKDIR /app\nCOPY go.mod .\n\nRUN go mod tidy\nCOPY ../.. .\n-- README.md --\n# QUIC START\n\n- this app was made by github.com/hagarihayato/sprint\n-- application/sample_application.go --\npackage application\n\nimport (\n\t\"@@.ImportPath@@/domain/model\"\n\t\"@@.ImportPath@@/domain/repository\"\n)\n\ntype (\n\tsampleApplication struct {\n\t\trepository.SampleRepository\n\t}\n\tSampleApplication interface {\n\t\tGetSamples() ([]*model.Sample, error)\n\t}\n)\n\nfunc NewSampleApplication(rs repository.SampleRepository) SampleApplication {\n\treturn &sampleApplication{rs}\n}\n\nfunc (s *sampleApplication) GetSamples() ([]*model.Sample, error) {\n\treturn s.SampleRepository.Fetch()\n}\n-- build.sh --\n#!/bin/bash\n\ngo build -o app && ./app\n-- docker-compose.yml --\nversion: \"3.5\"\nservices:\n  app:\n    container_name: app\n    build: \"\"\n    tty: true\n    restart: always\n    volumes:\n      - .:/app\n    ports:\n      - 8080:8080\n    command: sh ./build.sh\n  postgres:\n    image: postgres:10-alpine\n    container_name: postgres\n    ports:\n      - \"5432:5432\"\n    environment:\n      - POSTGRES_USER=postgres\n      - POSTGRES_PASSWORD=postgres\n      - PGPASSWORD=postgres\n      - POSTGRES_DB=postgres\n      - DATABASE_HOST=localhost\n    volumes:\n      - $PWD/infrastructure/postgres/init:/docker-entrypoint-initdb.d\n-- domain/model/sample_model.go --\npackage model\n\ntype Sample struct {\n\tID int64\n\tText string\n}\n-- domain/repository/sample_repository.go --\npackage repository\n\nimport (\n\t\"database/sql\"\n\t\"@@.ImportPath@@/domain/model\"\n)\n\ntype (\n\tsampleRepository struct{\n\t\tconn *sql.DB\n\t}\n\tSampleRepository interface {\n\t\tFetch() ([]*model.Sample, error)\n\t}\n)\n\nfunc NewSampleRepository(Conn *sql.DB) SampleRepository {\n\treturn &sampleRepository{Conn}\n}\n\nfunc (s *sampleRepository) Fetch() ([]*model.Sample, error) {\n\tvar samples []*model.Sample\n\trows, err := s.conn.Query(\"SELECT id, text FROM samples;\")\n\tif rows == nil { return nil, err }\n\tfor rows.Next() {\n\t\tsample := &model.Sample{}\n\t\terr = rows.Scan(&sample.ID, &sample.Text)\n\t\tif err == nil {\n\t\t\tsamples = append(samples, sample)\n\t\t}\n\t}\n\treturn samples, err\n}\n-- go.mod --\nmodule sprinter\n\ngo 1.15\n\nrequire (\n\tgithub.com/go-chi/chi v4.1.2+incompatible\n\tgithub.com/lib/pq v1.8.0\n\tgithub.com/stretchr/objx v0.3.0 // indirect\n\tgithub.com/stretchr/testify v1.6.1\n\tgolang.org/x/net v0.0.0-20201002202402-0a1ea396d57c // indirect\n\tgolang.org/x/tools v0.0.0-20201002184944-ecd9fd270d5d // indirect\n\tgopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect\n)\n-- infrastructure/postgres/conf/database.go --\npackage conf\n\nimport (\n\t\"database/sql\"\n\t\"fmt\"\n\t_ \"github.com/lib/pq\"\n)\n\nvar (\n\tDRIVER = \"postgres\"\n\tHOSTNAME = \"postgres\"\n\tUSER = \"postgres\"\n\tDBNAME = \"postgres\"\n\tPASSWORD = \"postgres\"\n)\n\nfunc NewDatabaseConnection() (*sql.DB, error) {\n\tsource := fmt.Sprintf(\"host=%s user=%s dbname=%s password=%s sslmode=disable\", HOSTNAME, USER, DBNAME, PASSWORD)\n\tconn, err := sql.Open(DRIVER, source)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn conn, nil\n}\n\n-- infrastructure/postgres/init/1_init.sql --\nDROP TABLE IF EXISTS samples;\n\nCREATE TABLE IF NOT EXISTS samples\n(\n   id SERIAL NOT NULL,\n   text TEXT NOT NULL,\n   PRIMARY KEY (id)\n);\n\nINSERT INTO samples(text) VALUES ('sample');\n-- interactor/interactor.go --\npackage interactor\n\nimport (\n\t\"database/sql\"\n\t\"@@.ImportPath@@/application\"\n\t\"@@.ImportPath@@/domain/repository\"\n\t\"@@.ImportPath@@/presenter/handler\"\n)\n\ntype (\n\tinteractor struct {\n\t\tconn *sql.DB\n\t}\n\tInteractor interface {\n\t\tNewRepository() *Repository\n\t\tNewApplication(r *Repository) *Application\n\t\tNewHandler(a *Application) *Handler\n\t}\n\tRepository struct {\n\t\trepository.SampleRepository\n\t}\n\tApplication struct {\n\t\tapplication.SampleApplication\n\t}\n\tHandler struct {\n\t\thandler.SampleHandler\n\t}\n)\n\nfunc NewInteractor(conn *sql.DB) Interactor {\n\treturn &interactor{conn}\n}\n\nfunc (i *interactor) NewRepository() *Repository {\n\tr := &Repository{}\n\tr.SampleRepository = repository.NewSampleRepository(i.conn)\n\treturn r\n}\n\nfunc (i *interactor) NewApplication(r *Repository) *Application {\n\ta := &Application{}\n\ta.SampleApplication = application.NewSampleApplication(r.SampleRepository)\n\treturn a\n}\n\nfunc (i *interactor) NewHandler(a *Application) *Handler {\n\th := &Handler{}\n\th.SampleHandler = handler.NewSampleHandler(a.SampleApplication)\n\treturn h\n}\n\n\n\n\n\n-- main.go --\npackage main\n\nimport (\n\t\"fmt\"\n\t\"net/http\"\n\t\"@@.ImportPath@@/infrastructure/postgres/conf\"\n\t\"@@.ImportPath@@/interactor\"\n\t\"@@.ImportPath@@/presenter/middleware\"\n\t\"@@.ImportPath@@/presenter/router\"\n)\n\nfunc main() {\n\tconn, err := conf.NewDatabaseConnection()\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer conn.Close()\n\tfmt.Println(`\n    * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * \n  *        ####    #####    #####     ####    ##  ##   ######   ######   #####  *\n  *      ##  ##   ##  ##   ##  ##     ##     ### ##     ##     ##       ##  ##  *\n  *     ##       ##  ##   ##  ##     ##     ######     ##     ##       ##  ##   *\n  *     ####    #####    #####      ##     ######     ##     ####     #####     *\n  *       ##   ##       ####       ##     ## ###     ##     ##       ####       *\n  *  ##  ##   ##       ## ##      ##     ##  ##     ##     ##       ## ##       *\n  *  ####    ##       ##  ##    ####    ##  ##     ##     ######   ##  ##       *\n    * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *\n    `)\n\n\tfmt.Println(`\tGET http://localhost:8080/api/v1`)\n\ti := interactor.NewInteractor(conn)\n\tr := i.NewRepository()\n\ta := i.NewApplication(r)\n\th := i.NewHandler(a)\n\tm := middleware.NewMiddleware()\n\ts := router.NewRouter()\n\ts.Router(h, m)\n\n\t_ = http.ListenAndServe(\":8080\", s.Route)\n}\n-- presenter/handler/sample_handler.go --\npackage handler\n\nimport (\n\t\"encoding/json\"\n\t\"log\"\n\t\"net/http\"\n\t\"@@.ImportPath@@/application\"\n\t\"@@.ImportPath@@/domain/model\"\n)\n\ntype(\n\tsampleHandler struct {\n\t\tapplication.SampleApplication\n\t}\n\tSampleHandler interface {\n\t\tSampleIndex(w http.ResponseWriter, r *http.Request)\n\t}\n\tresponse struct {\n\t\tStatus int\n\t\tSamples []*model.Sample\n\t}\n)\n\nfunc NewSampleHandler(as application.SampleApplication) SampleHandler {\n\treturn &sampleHandler{as}\n}\n\nfunc (s *sampleHandler) SampleIndex(w http.ResponseWriter, r *http.Request) {\n\tsamples, err := s.SampleApplication.GetSamples()\n\n\tif err != nil {\n\t\thttp.Error(w, err.Error(), http.StatusNotFound)\n\t}\n\n\tresp := &response{\n\t\tStatus: http.StatusOK,\n\t\tSamples: samples,\n\t}\n\n\tres, err := json.Marshal(resp)\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\n\t_ , _ = w.Write(res)\n}\n-- presenter/middleware/main.go --\npackage middleware\n\ntype middleware struct {}\n\ntype Middleware interface {}\n\nfunc NewMiddleware() Middleware {\n\treturn &middleware{}\n}\n-- presenter/router/router.go --\npackage router\n\nimport (\n\t\"github.com/go-chi/chi\"\n\t\"github.com/go-chi/chi/middleware\"\n\t\"@@.ImportPath@@/interactor\"\n\tmid \"@@.Important@@/presenter/middleware\"\n)\n\ntype Server struct {\n\tRoute *chi.Mux\n}\n\nfunc NewRouter() *Server {\n\treturn &Server{\n\t\tRoute: chi.NewRouter(),\n\t}\n}\n\nfunc (s *Server) Router(h *interactor.Handler, m mid.Middleware) {\n\ts.Route.Use(middleware.Logger)\n\ts.Route.Use(middleware.Recoverer)\n\ts.Route.Route(\"/api/v1\", func(r chi.Router) {\n\t\tr.Get(\"/\", h.SampleHandler.SampleIndex)\n\t\t// TODO\n\t})\n}\n-- presenter/router/router_test.go --\npackage router_test\n\nimport (\n\t\"github.com/go-chi/chi\"\n\t\"github.com/stretchr/testify/assert\"\n\t\"reflect\"\n\t\"@@.ImportPath@@/presenter/router\"\n\t\"testing\"\n)\n\nfunc TestNewRouter(t *testing.T) {\n\troute := chi.NewRouter()\n\twantNewRouterForm := &router.Server{\n\t\tRoute: route,\n\t}\n\n\tr := router.NewRouter()\n\n\tv := reflect.ValueOf(r)\n\tw := reflect.ValueOf(wantNewRouterForm)\n\n\tassert.Equal(t, v.Type(), w.Type())\n}\n"))
