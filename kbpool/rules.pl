far(X,Y) :-	seen(X,K,MIDPORTX), 
		seen(Y,K,MIDPORTY),
		seen(X,Y,PORTX),
		seen(Y,X,PORTY),
		switch(K), switch(X), switch(Y), X \= Y, Y \= K, X \= K, 
		PORTX == MIDPORTX,
		PORTY == MIDPORTY.
 
direct(XNAME,YNAME) :-	switchname(X,XNAME),
		switchname(Y,YNAME),
		\+ far(X,Y).

directp(XNAME,YNAME,PORTX,PORTY) :- switchname(X,XNAME),
		switchname(Y,YNAME),
		seen(X,Y,PORTX),
		seen(Y,X,PORTY),
		\+ far(X,Y).

edgeswitch(XNAME) :-	switchname(X,XNAME),
			seen(X,Y,PORT1),
			seen(X,Z,PORT2),
			switch(X), switch(Y), switch(Z), X \= Y, X \= Z, Y \= Z,
			PORT1 == PORT2.

edgeswitchp(XNAME,PORT) :-	switchname(X,XNAME),
			seen(X,Y,PORT),
			seen(X,Z,PORT),
			switch(X), switch(Y), switch(Z), X \= Y, X \= Z, Y \= Z.
