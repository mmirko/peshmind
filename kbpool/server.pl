:- use_module(library(pengines)).
:- use_module(library(http/http_server)).

% launch an HTTP server with pengine endpoint /p
:- http_server(http_dispatch, [port(3030)]).
