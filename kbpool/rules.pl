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

directp(X,Y,PORTX,PORTY) :- seen(X,K,MIDPORTX),
		seen(Y,K,MIDPORTY),
		seen(X,Y,PORTX),
		seen(Y,X,PORTY),
		switch(K), switch(X), switch(Y), X \= Y, Y \= K, X \= K,
		(PORTX \= MIDPORTX ; PORTY \= MIDPORTY).

