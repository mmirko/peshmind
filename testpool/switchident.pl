
name(x, 'SX').
name(w, 'SW').
name(z, 'SZ').
seen(x,d,1).
seen(x,w,8).
seen(x,c,8).
seen(x,a,8).
seen(x,b,8).
seen(x,z,8).
seen(w,c,6).
seen(w,b,7).
seen(w,x,10).
seen(w,d,10).
seen(w,a,4).
seen(w,z,4).
seen(w,b,4).
seen(z,a,1).
seen(z,b,2).
seen(z,c,5).
seen(z,d,5).
seen(z,w,5).
seen(z,x,5).

:- dynamic switch/1.

build_switch_index :-
    setof(Switch, Host^Port^seen(Switch, Host, Port), Switches),
    forall(member(S, Switches),
	   assertz(switch(S))).

:- build_switch_index.

far(X,Y) :-	seen(X,K,MIDPORTX), 
		seen(Y,K,MIDPORTY),
		seen(X,Y,PORTX),
		seen(Y,X,PORTY),
		switch(K), switch(X), switch(Y), X \= Y, Y \= K, X \= K, 
		PORTX == MIDPORTX,
		PORTY == MIDPORTY.

direct(X,Y) :-	seen(X,K,MIDPORTX),
		seen(Y,K,MIDPORTY),
		seen(X,Y,PORTX),
		seen(Y,X,PORTY),
		switch(K), switch(X), switch(Y), X \= Y, Y \= K, X \= K,
		(PORTX \= MIDPORTX ; PORTY \= MIDPORTY).

edgeswitch(XNAME) :-	switchname(X,XNAME),
			seen(X,Y,PORT1),
			seen(X,Z,PORT2),
			switch(X), switch(Y), switch(Z), X \= Y, X \= Z, Y \= Z,
			PORT1 == PORT2.