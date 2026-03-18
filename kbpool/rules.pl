% Two switches are far if they are connected to a third switch on the same port as each other. (it means a switch among the two)
far(X,Y) :-	seen(X,K,MIDPORTX), 
		seen(Y,K,MIDPORTY),
		seen(X,Y,PORTX),
		seen(Y,X,PORTY),
		switch(K), switch(X), switch(Y), X \= Y, Y \= K, X \= K, 
		PORTX == MIDPORTX,
		PORTY == MIDPORTY.

% Far switches by name
farn(XNAME,YNAME) :- switchname(X,XNAME),
		switchname(Y,YNAME),
		far(X,Y).

% Two switches are directly connected if they are not far.
direct(X,Y) :- \+ far(X,Y).

% Directly connected switches by name
directn(XNAME,YNAME) :- switchname(X,XNAME),
		switchname(Y,YNAME),
		direct(X,Y).

% Two switches directly connected on specific ports
directp(X,Y,PORTX,PORTY) :- seen(X,Y,PORTX),
		seen(Y,X,PORTY),
		\+ far(X,Y).

% Two directly connected switches on specific ports by name
directpn(XNAME,YNAME,PORTX,PORTY) :- switchname(X,XNAME),
		switchname(Y,YNAME),
		directp(X,Y,PORTX,PORTY).		

% A switch is an internal switch if it is connected to two different switches on different ports.
internalswitch(X) :- seen(X,Y,PORT1),
		seen(X,Z,PORT2),
		switch(X), switch(Y), switch(Z), X \= Y, X \= Z, Y \= Z,
		PORT1 \= PORT2.

% Internal switches by name
internalswitchn(XNAME) :- switchname(X,XNAME),
		internalswitch(X).

% A switch is an edge switch if it is not an internal switch.
edgeswitch(X) :- switch(X), \+ internalswitch(X).

% Edge switches by name
edgeswitchn(XNAME) :- switchname(X,XNAME),
			edgeswitch(X).

% An enge switch and it's port
edgeswitchp(X,PORT) :-  seen(X,Y,PORT),
		\+ internalswitch(X),
		switch(X), switch(Y), X \= Y.

% An edge switch and it's port by name
edgeswitchpn(XNAME,PORT) :- switchname(X,XNAME),
			edgeswitchp(X,PORT).
