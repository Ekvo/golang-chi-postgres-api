// This application: gets, saves, reads, updates and deletes a Task object, using the following patterns and tools.
//==========================================================================================================
// REST API property and features.
//==========================================================================================================
// Design patterns
/*
* CRUD
* Data Transfer Object (DTO)
* SOLID
 */
//==========================================================================================================
// DataBase    - PostgresSQL
// multiplexer - chi
//==========================================================================================================

// package main ~> ../cmd/app
// logic of application

// package config ~> ../internal/config
// parse data for run application from file or ENV
/*
 - config.go
 * struct - Config      - contain critical data for run of application
 * func   - NewConfig
 * func   - getNameENV  - returns array of string  with hanes all name of ENV variables
 * func   - validConfig - member of Config - create 'common.Message' see pkg/common/common.go
check all fields for validity. If field after viper.Unmarhal is broken exept 'DBNameForTest'
add name field (key) and set Error(value).
if 'common.Message' is not empty, an error describing all corrupted fields is returned
*/

// package model ~> ../internal/model
// describe property of Task - object stored in the database
/*
 - model.go
 * struct - Task
 * 4 interface - object maintenance in strore
*/

// packege server ~> ../internal/server
// rules for use http.Server in application
/*
 - server.go
 * struct - Connect        - contain http.Server
 * func   - Init function  - get property from  config.Config for initialize http.Serve
 * func   - ListenAndServe - property of connect and shut http.Server
*/

// packege servises ~> ../internal/servises
// Data Transfer Object
/*
 - validator.go
 * struct - TaskValidator - rules for body from Request
 * func   - DecodeJSON    - TaskValidator member - get body for Task
 * func   - TaskModel     - return object Task
------------------------------------------------------------------------------------------------------------
 - serializer.go
 * struct - TaskSerializer     - rules for creating a body for ResponseWriter from one Task
 * func   - Response           - TaskSerializer member - create body of Task
 * struct - TaskListSerializer - body for ResponseWriter from array of Tasks
 * func   - Response           - member of TaskListSerializer
*/

// packege source ~> ../internal/source
// PostgresSQL - github.com/lib/pq
/*
 - source.go
 * struct - Dbinstance    - contain ptr of sql.DB
 * func   - Init function - get property from .env for initialize database
 * func   - URLParam      - created a string type URL to connect to the database
------------------------------------------------------------------------------------------------------------
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
 * func   - Timeout    - midddleware func
------------------------------------------------------------------------------------------------------------
 - middlweare.go
 * func - Timeout - middlweare function,
create context.WithTimeout,
request = request.WithContext(ctx),
call next(w,r)
------------------------------------------------------------------------------------------------------------
 - route.go
describe application handlers
 * struct - responseData - body and status for 'ResponseWriter'
 * type   - taskFunc     - layout of func for 'Decode', work with Store(DB) and Encode object
 * func TaskHandler - main function on route
accepts an interface for interaction with the database,
function 'taskFunc' describing the logic of processing the object and obtaining the result.
create chan 'responseData', call in goroutine function 'taskFunc' for create 'responseData',
in select inside TaskHandler get data from chan 'responseData',
call 'common.EncodeJSON' for  create Response
 * func(s) - create, read, update and delete of Task
*/

// packege variables ~> ../internal/variables
// contain only const varaibles

// packege common ~> ../pkg/common
// universal utilities
/*
 - common.go
 * type   - Message      - map for ResponseWriter (map[string]any)
 * func   - String       - member of Message, create line {key1:value1},{key2:value2}...
 * struct - MessageError - wrap error for Response
 * func   - DecodeJSON   - rules for getting object body from Request
 * func   - EncodeJSON   - set COntent-type, code status, and write body for ResponseWriter
*/
package docs
