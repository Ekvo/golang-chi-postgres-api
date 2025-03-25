// REST API property and features.

// CRUD
// Data Transfer Object (DTO)
// SOLID

// DataBase    - PostgresSQL
// multiplexer - chi

// Rules for getting object from Request and write to Response
// - github.com/gin-gonic/gin/binding
// - github.com/go-playground/validator/v10

// package main ~> ../cmd/app
// logic of application

// package model ~> ../internal/model
// describe property of Task - object stored in the database
/*
 * struct - Task
 * 4 interface - object maintenance
 */

// packege server ~> ../internal/server
// rules for use http.Server in application
/*
 * struct - Connect        - contain http.Server
 * func   - Init function  - get property from .env for initialize http.Serve
 * func   - ListenAndServe - property of connect and shut http.Server
 */

// packege servises ~> ../internal/servises
// Data Transfer Object
/*
 * - validator.go
 * struct - TaskValidator - rules for body from Request
 * func   - Bind          - TaskValidator member - get body for Task
 ---------------------------------------------------------------------------
 * - serializer.go
 * struct - TaskSerializer - rules for body to ResponseWriter
 * func   - Response       - TaskSerializer member - create body of Task
*/

// packege source ~> ../internal/source
// PostgresSQL - github.com/lib/pq
/*
 - source.go
 * struct - Dbinstance    - contain ptr of sql.DB
 * func   - Init function - get property from .env for initialize database
 * func   - URLParam      - created a string type URL to connect to the database
 ---------------------------------------------------------------------------
 - query.go
 * describe logic of interfaces Task look. (look: package model ~> ../internal/model)
 * interface - RowScaner - logic for 'Scan' data from a database
*/

// packege transport ~> ../internal/transport
// chi - github.com/go-chi/chi/v5
/*
 - transport.go
 * struct - Transport  - contain ptr of chi.Mux
 * Routes - Transport member
 * func   - taskRoutes - logic application handlers
 * func   - Timeout    - middleware func
 ---------------------------------------------------------------------------
 - route.go
 * describe application handlers
*/

// packege variables ~> ../internal/variables
// contain only const varaibles

// packege common ~> ../pkg/common
// universal utilities
/*
 * struct - MessageError   - wrap error for Response
 * strucе - FiledValidator - custom error of 'validator.FieldError'
 * describes error getting an invalid object from Request

 * FieldName - FiledValidator member
 * func - NewMessageErrorFromValidator - create MessageError from 'validator.FieldError' object
 * func - Bind
 */
package docs
